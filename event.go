package main

import (
	"log"
	"os/exec"
	"strings"
)

var events = []Event{
	{
		Name: "打开电脑",
		Match: func(query string) bool {
			return strings.Contains(query, "打开电脑")
		},
		Handle: func(query string) error {
			return runScript("/root/wake.sh")
		},
	},
	{
		Name: "关闭电脑",
		Match: func(query string) bool {
			return strings.Contains(query, "关闭电脑")
		},
		Handle: func(query string) error {
			return runScript("/root/sleep.sh")
		},
	},
}

type Event struct {
	Name   string
	Match  func(query string) bool
	Handle func(query string) error
}

func DispatchEvent(query string) (string, error) {
	for _, ev := range events {
		if ev.Match(query) {
			log.Printf("[EVENT] matched: %s, query=%q\n", ev.Name, query)
			if err := ev.Handle(query); err != nil {
				log.Println("[EVENT] handle failed:", err)
				return ev.Name, err
			}
			return ev.Name, nil
		}
	}
	return "", nil
}

func runScript(scriptPath string) error {
	cmd := exec.Command("/bin/sh", scriptPath)
	output, err := cmd.CombinedOutput()
	errMsg := "nil"
	if err != nil {
		errMsg = err.Error()
	}
	log.Printf("run, scriptPath: %s, output: %s, err: %s\n", scriptPath, output, errMsg)
	if err != nil {
		return err
	}
	return nil
}
