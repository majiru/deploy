package conf

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"	
	"strings"
)

type Conf struct {
	Cmd    []string
	Hosts  []string
	Script string
}

func ReadConf(r io.Reader) (*Conf, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	conf := &Conf{}
	if err = json.Unmarshal(b, conf); err != nil {
		return nil, err
	}
	return conf, nil
}

var ErrNoHost = errors.New("CmdList: no {{host}} in conf.Cmd")

const hostTemplate = "{{HOST}}"

func (conf *Conf) CmdList() ([]string, error) {
	var out []string
	var pivot int = -1
	var hostStr string = ""
	for i, c := range conf.Cmd {
		if strings.Contains(c, hostTemplate) {
			hostStr = c
			pivot = i
			break
		}
	}
	if pivot == -1 {
		return nil, ErrNoHost
	}
	for _, h := range conf.Hosts {
		b := []string{}
		b = append(b, conf.Cmd[:pivot]...)
		b = append(b, strings.Replace(hostStr, hostTemplate, h, -1))
		if pivot+1 < len(conf.Cmd) {
			b = append(b, conf.Cmd[pivot+1:]...)
		}
		out = append(out, strings.Join(b, " "))
	}
	return out, nil
}
