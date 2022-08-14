package usi

import (
	"bufio"
	"context"
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/kk-no/go-usi/usicmd"
)

type ReadWriteProcessor interface {
	Start(ctx context.Context)
	Stop() error
	Write(ctx context.Context)
	Read(ctx context.Context)
	SendCommand(command usicmd.Command)
}

type process struct {
	cmd       *exec.Cmd
	wg        *sync.WaitGroup
	cancel    context.CancelFunc
	procIn    io.WriteCloser
	procOut   io.ReadCloser
	sendQueue chan usicmd.Command
}

func NewReadWriteProcessor(ctx context.Context, name string) (ReadWriteProcessor, error) {
	p := &process{}
	p.wg = &sync.WaitGroup{}
	p.cmd = exec.CommandContext(ctx, name)
	p.sendQueue = make(chan usicmd.Command)

	var err error
	if p.procIn, err = p.cmd.StdinPipe(); err != nil {
		return nil, err
	}

	if p.procOut, err = p.cmd.StdoutPipe(); err != nil {
		return nil, err
	}

	if err := p.cmd.Start(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *process) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	p.wg.Add(1)
	go p.Read(ctx)

	p.wg.Add(1)
	go p.Write(ctx)
}

func (p *process) Stop() error {
	if p.cancel != nil {
		p.cancel()
		p.wg.Wait()
	}
	if p.sendQueue != nil {
		close(p.sendQueue)
	}
	if p.procIn != nil {
		if err := p.procIn.Close(); err != nil {
			return err
		}
	}
	if p.procOut != nil {
		if err := p.procOut.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *process) Write(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.wg.Done()
			return
		default:
			command := <-p.sendQueue
			log.Println(">", command)
			if _, err := p.procIn.Write([]byte(command + "\n")); err != nil {
				return
			}
		}
	}
}

func (p *process) Read(ctx context.Context) {
	scanner := bufio.NewScanner(p.procOut)
	for {
		select {
		case <-ctx.Done():
			p.wg.Done()
			return
		default:
			if scanner.Scan() {
				log.Println("<", scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return
			}
		}
	}
}

func (p *process) SendCommand(command usicmd.Command) {
	p.sendQueue <- command
}
