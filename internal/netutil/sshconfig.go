package netutil

import (
	"bufio"
	"bytes"
	"strings"
	"unicode"
)

type SSHHost struct {
	Alias        string
	HostName     string
	User         string
	Port         string
	IdentityFile string
}

// ParseSSHConfig lists the hosts defined in an ssh_config file, one entry
// per name that can be passed to ssh. Wildcard patterns such as "Host *"
// and Match blocks only hold defaults, so they are not listed. As in ssh,
// the first value given for a keyword wins.
func ParseSSHConfig(data []byte) ([]SSHHost, error) {
	var hosts []SSHHost
	var block []int // hosts defined by the Host line being read

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		key, value, ok := splitSSHConfigLine(scanner.Text())
		if !ok {
			continue
		}

		switch key {
		case "host":
			block = nil
			for _, name := range concreteHostNames(value) {
				hosts = append(hosts, SSHHost{Alias: name})
				block = append(block, len(hosts)-1)
			}
		case "match":
			block = nil
		default:
			for _, i := range block {
				setSSHOption(&hosts[i], key, value)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return hosts, nil
}

// splitSSHConfigLine splits "Keyword value", "Keyword=value" and
// "Keyword = value" (with any whitespace) into a lower-cased keyword and
// its value.
func splitSSHConfigLine(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}

	i := strings.IndexFunc(line, func(r rune) bool {
		return unicode.IsSpace(r) || r == '='
	})
	if i < 0 {
		return "", "", false
	}

	value := strings.TrimSpace(line[i:])
	value = strings.TrimSpace(strings.TrimPrefix(value, "="))
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-1]
	}
	if value == "" {
		return "", "", false
	}

	return strings.ToLower(line[:i]), value, true
}

// concreteHostNames drops wildcard and negated patterns from a Host line,
// what remains are names that can be passed to ssh as they are.
func concreteHostNames(patterns string) []string {
	var names []string
	for _, p := range strings.Fields(patterns) {
		if strings.ContainsAny(p, "*?") || strings.HasPrefix(p, "!") {
			continue
		}
		names = append(names, p)
	}
	return names
}

func setSSHOption(h *SSHHost, key, value string) {
	var field *string
	switch key {
	case "hostname":
		field = &h.HostName
	case "user":
		field = &h.User
	case "port":
		field = &h.Port
	case "identityfile":
		field = &h.IdentityFile
	default:
		return
	}

	if *field == "" {
		*field = value
	}
}
