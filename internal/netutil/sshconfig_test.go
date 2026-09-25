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

Host *
    ServerAliveInterval 60
`

	expected := []SSHHost{
		{Alias: "myserver", HostName: "10.0.0.5", User: "deploy"},
		{Alias: "tabbed", HostName: "10.0.0.6", Port: "2222"},
		{Alias: "eq", HostName: "10.0.0.7", User: "admin", IdentityFile: "~/.ssh/my key"},
		{Alias: "web1", User: "www"},
		{Alias: "web2", User: "www"},
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
