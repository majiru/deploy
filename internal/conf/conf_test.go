package conf

import (
	"strings"
	"testing"
)

const sampleConf = `
{
	"Cmd": ["ssh", "{{HOST}}", "sh"],
	"Hosts": ["localhost", "192.168.0.5"],
	"Script": "./deploy.sh"
}`

func TestReadConf(t *testing.T) {
	cf := strings.NewReader(sampleConf)
	conf, err := ReadConf(cf)
	if err != nil {
		t.Error(err)
	}
	if conf.Cmd[0] != "ssh" || conf.Cmd[1] != "{{HOST}}" || conf.Cmd[2] != "sh" {
		t.Errorf("Cmd mismatch")
	}
	if conf.Hosts[0] != "localhost" || conf.Hosts[1] != "192.168.0.5" {
		t.Errorf("Hosts mismatch")
	}
	if conf.Script != "./deploy.sh" {
		t.Errorf("Script mismatch")
	}
}

func TestCmdList(t *testing.T) {
	cf := strings.NewReader(sampleConf)
	conf, err := ReadConf(cf)
	if err != nil {
		t.Error(err)
	}
	l, err := conf.CmdList()
	if err != nil {
		t.Error(err)
	}
	const l0 = "ssh localhost sh"
	if l[0] != l0 {
		t.Errorf("List mismatch: Got %s, expected %s", l0, l[0])
	}
	const l1 = "ssh 192.168.0.5 sh"
	if l[1] != l1 {
		t.Errorf("List mismatch: Got %s, expected %s", l1, l[1])
	}
}

func TestCmdListNoHost(t *testing.T) {
	conf := &Conf{Cmd: []string{"ssh", "sh"}}
	_, err := conf.CmdList()
	if err != ErrNoHost {
		t.Errorf("Did not get ErrNoHost without {{HOST}}")
	}
}
