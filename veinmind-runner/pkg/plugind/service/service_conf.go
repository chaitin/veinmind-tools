package service

import (
	"github.com/google/uuid"
	"os"
	"os/exec"
	"time"
)

type Runner struct {
	Uuid    string
	Command string
	Stderr  *os.File
	Stdout  *os.File
	Cmd     *exec.Cmd
	TimeOut time.Duration
	Signal  chan string
}

type Conf struct {
	Command   string `toml:"Command"`
	StdoutLog string `toml:"StdoutLog"`
	StderrLog string `toml:"StderrLog"`
	TimeOut   int    `toml:"TimeOut"`
}

func NewRunner(s Conf) (*Runner, error) {
	StderrFile, err := os.OpenFile(s.StderrLog, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	StdoutFile, err := os.OpenFile(s.StdoutLog, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Runner{
		Uuid:    uuid.New().String(),
		Command: s.Command,
		Stdout:  StdoutFile,
		Stderr:  StderrFile,
		TimeOut: time.Duration(s.TimeOut) * time.Second,
		Signal:  make(chan string),
	}, nil
}
