package netutil

import (
	"reflect"
	"testing"
)

func TestParseSSHConfig(t *testing.T) {
	config := `# personal servers
Host myserver
    HostName 10.0.0.5
    User deploy
    User ignored

Host	tabbed
	HostName	10.0.0.6
	Port	2222

Host eq
    HostName=10.0.0.7
    User = admin
    IdentityFile "~/.ssh/my key"

Match host foo
    User matchuser

Host web1 web2 *.internal !bad
    User www

Host pi pi.local
    HostName pi.local
    User milos

Host *
    ServerAliveInterval 60
`

	expected := []SSHHost{
		{Aliases: []string{"myserver"}, HostName: "10.0.0.5", User: "deploy"},
		{Aliases: []string{"tabbed"}, HostName: "10.0.0.6", Port: "2222"},
		{Aliases: []string{"eq"}, HostName: "10.0.0.7", User: "admin", IdentityFile: "~/.ssh/my key"},
		// no HostName, so each alias connects to a different host
		{Aliases: []string{"web1"}, User: "www"},
		{Aliases: []string{"web2"}, User: "www"},
		// a HostName shared by both aliases, listed once
		{Aliases: []string{"pi", "pi.local"}, HostName: "pi.local", User: "milos"},
	}

	hosts, err := ParseSSHConfig([]byte(config))
	if err != nil {
		t.Fatalf("ParseSSHConfig() error = %v", err)
	}
	if !reflect.DeepEqual(hosts, expected) {
		t.Errorf("ParseSSHConfig() =\n%+v\nwant\n%+v", hosts, expected)
	}
}

func TestParseSSHConfigEmpty(t *testing.T) {
	hosts, err := ParseSSHConfig([]byte("# only defaults\nHost *\n    User me\n"))
	if err != nil {
		t.Fatalf("ParseSSHConfig() error = %v", err)
	}
	if len(hosts) != 0 {
		t.Errorf("ParseSSHConfig() = %+v, want no hosts", hosts)
	}
}
