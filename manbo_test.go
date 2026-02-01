package main

import (
	"log"
	"testing"

	"github.com/toszyrt/manbo/openclaw"
)

func TestOpenclawHello(t *testing.T) {
	c := openclaw.NewWsClient(nil)
	c.ConnectAndHandshake()
	if _, err := c.SendQueryAndWait("你好"); err != nil {
		log.Fatal(err)
	}
}

func TestManbo(t *testing.T) {
	m := NewManbo()
	m.Init()
	m.Loop()
}
