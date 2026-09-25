package netutil

import "testing"

func TestParseSubnetIPv4(t *testing.T) {
	tests := []struct {
		name          string
		addr          string
		mask          string
		wantNetwork   string
		wantBroadcast string
		wantFirst     string
		wantLast      string
		wantUsable    string
		wantTotal     string
	}{
		{"cidr /24", "192.168.1.10/24", "", "192.168.1.0/24", "192.168.1.255", "192.168.1.1", "192.168.1.254", "254", "256"},
		{"dotted mask", "192.168.1.10", "255.255.255.0", "192.168.1.0/24", "192.168.1.255", "192.168.1.1", "192.168.1.254", "254", "256"},
		{"hex mask", "192.168.1.10", "0xffffff00", "192.168.1.0/24", "192.168.1.255", "192.168.1.1", "192.168.1.254", "254", "256"},
		{"hex mask uppercase", "192.168.1.10", "0XFFFFFF00", "192.168.1.0/24", "192.168.1.255", "192.168.1.1", "192.168.1.254", "254", "256"},
		{"cidr /30", "10.0.0.5/30", "", "10.0.0.4/30", "10.0.0.7", "10.0.0.5", "10.0.0.6", "2", "4"},
		{"point-to-point /31", "10.0.0.4/31", "", "10.0.0.4/31", "10.0.0.5", "10.0.0.4", "10.0.0.5", "2", "2"},
		{"host route /32", "10.0.0.4/32", "", "10.0.0.4/32", "10.0.0.4", "10.0.0.4", "10.0.0.4", "1", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseSubnet(tt.addr, tt.mask)
			if err != nil {
				t.Fatalf("ParseSubnet(%q, %q) unexpected error: %v", tt.addr, tt.mask, err)
			}
			if info.Network != tt.wantNetwork {
				t.Errorf("Network = %q, want %q", info.Network, tt.wantNetwork)
			}
			if info.Broadcast != tt.wantBroadcast {
				t.Errorf("Broadcast = %q, want %q", info.Broadcast, tt.wantBroadcast)
			}
			if info.FirstHost != tt.wantFirst {
				t.Errorf("FirstHost = %q, want %q", info.FirstHost, tt.wantFirst)
			}
			if info.LastHost != tt.wantLast {
				t.Errorf("LastHost = %q, want %q", info.LastHost, tt.wantLast)
			}
			if info.UsableHosts != tt.wantUsable {
				t.Errorf("UsableHosts = %q, want %q", info.UsableHosts, tt.wantUsable)
			}
			if info.TotalHosts != tt.wantTotal {
				t.Errorf("TotalHosts = %q, want %q", info.TotalHosts, tt.wantTotal)
			}
		})
	}
}

func TestParseSubnetClassAndPrivate(t *testing.T) {
	tests := []struct {
		addr      string
		wantClass string
		wantPriv  bool
	}{
		{"10.0.0.1/8", "A", true},
		{"172.16.0.1/12", "B", true},
		{"192.168.1.1/24", "C", true},
		{"8.8.8.8/32", "A", false},
		{"224.0.0.1/24", "D (multicast)", false},
	}

	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			info, err := ParseSubnet(tt.addr, "")
			if err != nil {
				t.Fatalf("ParseSubnet(%q) unexpected error: %v", tt.addr, err)
			}
			if info.Class != tt.wantClass {
				t.Errorf("Class = %q, want %q", info.Class, tt.wantClass)
			}
			if info.Private != tt.wantPriv {
				t.Errorf("Private = %v, want %v", info.Private, tt.wantPriv)
			}
		})
	}
}

func TestParseSubnetIPv6(t *testing.T) {
	info, err := ParseSubnet("2001:db8::1/64", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !info.IsIPv6 {
		t.Errorf("IsIPv6 = false, want true")
	}
	if info.Network != "2001:db8::/64" {
		t.Errorf("Network = %q, want %q", info.Network, "2001:db8::/64")
	}
	if info.FirstHost != "2001:db8::" {
		t.Errorf("FirstHost = %q, want %q", info.FirstHost, "2001:db8::")
	}
	if info.LastHost != "2001:db8::ffff:ffff:ffff:ffff" {
		t.Errorf("LastHost = %q, want %q", info.LastHost, "2001:db8::ffff:ffff:ffff:ffff")
	}
}

func TestParseSubnetIPv4Mapped(t *testing.T) {
	info, err := ParseSubnet("::ffff:192.168.1.10/120", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.IsIPv6 {
		t.Errorf("IsIPv6 = true, want false")
	}
	if info.Network != "192.168.1.0/24" {
		t.Errorf("Network = %q, want %q", info.Network, "192.168.1.0/24")
	}
	if info.Broadcast != "192.168.1.255" {
		t.Errorf("Broadcast = %q, want %q", info.Broadcast, "192.168.1.255")
	}
	if info.UsableHosts != "254" {
		t.Errorf("UsableHosts = %q, want %q", info.UsableHosts, "254")
	}
}

func TestParseSubnetErrors(t *testing.T) {
	tests := []struct {
		name string
		addr string
		mask string
	}{
		{"invalid cidr", "not-an-ip/24", ""},
		{"invalid ip with mask", "not-an-ip", "255.255.255.0"},
		{"invalid mask", "192.168.1.1", "not-a-mask"},
		{"invalid hex mask", "192.168.1.1", "0xzzzzzzzz"},
		{"mask with ipv6", "2001:db8::1", "255.255.255.0"},
		{"non-contiguous mask", "192.168.1.10", "255.0.255.0"},
		{"non-contiguous hex mask", "192.168.1.10", "0xff00ff00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseSubnet(tt.addr, tt.mask); err == nil {
				t.Errorf("ParseSubnet(%q, %q) expected error, got nil", tt.addr, tt.mask)
			}
		})
	}
}
