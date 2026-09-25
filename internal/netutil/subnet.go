package netutil

import (
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
)

type SubnetInfo struct {
	Address     string
	Network     string
	Netmask     string
	Wildcard    string
	Broadcast   string
	FirstHost   string
	LastHost    string
	UsableHosts string
	TotalHosts  string
	Class       string
	Private     bool
	IsIPv6      bool
}

// ParseSubnet computes subnet details for addr, which may be in CIDR
// notation (e.g. "192.168.1.10/24"). If mask is non-empty, addr is treated
// as a bare IPv4 address and mask as its netmask instead, given either in
// dotted-decimal (255.255.255.0) or hex (0xffffff00) form.
func ParseSubnet(addr string, mask string) (*SubnetInfo, error) {
	var ip net.IP
	var ipNet *net.IPNet

	if mask != "" {
		ip = net.ParseIP(addr)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address: %s", addr)
		}
		if ip.To4() == nil {
			return nil, fmt.Errorf("netmask argument can only be used with an IPv4 address")
		}

		m, err := parseIPv4Mask(mask)
		if err != nil {
			return nil, err
		}

		ipNet = &net.IPNet{IP: ip.Mask(m), Mask: m}
	} else {
		var err error
		ip, ipNet, err = net.ParseCIDR(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid subnet %q: expected CIDR notation like 192.168.1.0/24 (or pass a netmask as a second argument)", addr)
		}
	}

	if ip4 := ip.To4(); ip4 != nil {
		// an IPv4-mapped IPv6 CIDR (::ffff:192.168.1.1/120) comes with a
		// 128-bit mask, of which only the last 32 bits cover the IPv4 part
		if ones, bits := ipNet.Mask.Size(); bits == 8*net.IPv6len {
			if ones < 96 {
				return ipv6Subnet(ip, ipNet), nil
			}
			ipNet = &net.IPNet{IP: ipNet.IP.To4(), Mask: net.CIDRMask(ones-96, 8*net.IPv4len)}
		}
		return ipv4Subnet(ip4, ipNet), nil
	}

	return ipv6Subnet(ip, ipNet), nil
}

func ipv4Subnet(ip net.IP, ipNet *net.IPNet) *SubnetInfo {
	ones, bits := ipNet.Mask.Size()
	network := ipToUint32(ipNet.IP.To4())
	mask := ipToUint32(net.IP(ipNet.Mask).To4())
	wildcard := ^mask
	broadcast := network | wildcard

	total := uint64(1) << uint(bits-ones)

	var first, last uint32
	var usable uint64

	switch ones {
	case 32:
		first, last, usable = network, network, 1
	case 31:
		first, last, usable = network, broadcast, 2
	default:
		first, last, usable = network+1, broadcast-1, total-2
	}

	return &SubnetInfo{
		Address:     ip.String(),
		Network:     fmt.Sprintf("%s/%d", uint32ToIP(network), ones),
		Netmask:     uint32ToIP(mask).String(),
		Wildcard:    uint32ToIP(wildcard).String(),
		Broadcast:   uint32ToIP(broadcast).String(),
		FirstHost:   uint32ToIP(first).String(),
		LastHost:    uint32ToIP(last).String(),
		UsableHosts: fmt.Sprintf("%d", usable),
		TotalHosts:  fmt.Sprintf("%d", total),
		Class:       ipv4Class(network),
		Private:     ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast(),
	}
}

func ipv6Subnet(ip net.IP, ipNet *net.IPNet) *SubnetInfo {
	ones, bits := ipNet.Mask.Size()

	network := new(big.Int).SetBytes(ipNet.IP.To16())
	total := new(big.Int).Lsh(big.NewInt(1), uint(bits-ones))
	last := new(big.Int).Sub(new(big.Int).Add(network, total), big.NewInt(1))

	return &SubnetInfo{
		Address:     ip.String(),
		Network:     fmt.Sprintf("%s/%d", ipNet.IP.String(), ones),
		Netmask:     net.IP(ipNet.Mask).String(),
		FirstHost:   ipNet.IP.String(),
		LastHost:    bigToIP(last).String(),
		UsableHosts: total.String(),
		TotalHosts:  total.String(),
		Private:     ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast(),
		IsIPv6:      true,
	}
}

func parseIPv4Mask(mask string) (net.IPMask, error) {
	var m net.IPMask

	if rest, ok := strings.CutPrefix(strings.ToLower(mask), "0x"); ok {
		n, err := strconv.ParseUint(rest, 16, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid netmask: %s", mask)
		}
		m = net.IPMask(uint32ToIP(uint32(n)).To4())
	} else {
		maskIP := net.ParseIP(mask)
		if maskIP == nil || maskIP.To4() == nil {
			return nil, fmt.Errorf("invalid netmask: %s", mask)
		}
		m = net.IPMask(maskIP.To4())
	}

	// Size reports 0 bits for masks like 255.0.255.0, which have no prefix length
	if _, bits := m.Size(); bits == 0 {
		return nil, fmt.Errorf("invalid netmask %s: mask bits must be contiguous (e.g. 255.255.255.0)", mask)
	}

	return m, nil
}

func ipToUint32(ip net.IP) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func uint32ToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

func bigToIP(n *big.Int) net.IP {
	b := n.Bytes()
	ip := make(net.IP, net.IPv6len)
	copy(ip[net.IPv6len-len(b):], b)
	return ip
}

func ipv4Class(network uint32) string {
	switch first := byte(network >> 24); {
	case first < 128:
		return "A"
	case first < 192:
		return "B"
	case first < 224:
		return "C"
	case first < 240:
		return "D (multicast)"
	default:
		return "E (reserved)"
	}
}
