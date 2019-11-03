package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/majiru/deploy/internal/conf"
)

func Run(stdout io.Writer, stderr io.Writer, script string, commands ...string) {
	var wg sync.WaitGroup
	wg.Add(len(commands))
	for _, c := range commands {
		cl := strings.Split(c, " ")
		cmd := exec.Command(cl[0], cl[1:]...)
		cmd.Stdin = strings.NewReader(script)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Start(); err != nil {
			log.Printf("Error starting command %s: %v\n", c, err)
			continue
		}
		go func(cmd *exec.Cmd){
			if err := cmd.Wait(); err != nil {
				log.Printf("Error running command %v: %v\n", cmd, err)
			}
			wg.Done()
		}(cmd)
	}
	wg.Wait()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Please specify at least one config file\n");
		os.Exit(1)
	}
	var jobs []*conf.Conf
	for _, f := range os.Args[1:] {
		cf, err := os.Open(f)
		if err != nil {
			log.Fatal("Error opening config file:", err)
		}
		conf, err := conf.ReadConf(cf)
		if err != nil {
			log.Fatal("Error parsing config file:", err)
		}
		jobs = append(jobs, conf)
		cf.Close()
	}
	for _, j := range jobs {
		s, err := os.Open(j.Script)
		if err != nil {
			log.Fatalf("Could not open script file %s: %v", j.Script, err)
		}
		b, err := ioutil.ReadAll(s)
		if err != nil {
			log.Fatalf("Could not read file %s: %v", j.Script, err)
		}
		cl, err := j.CmdList()
		if err != nil {
			log.Fatal("Error creating cmdList:", err)
		}
		Run(os.Stdout, os.Stderr, string(b), cl...)
	}
}
