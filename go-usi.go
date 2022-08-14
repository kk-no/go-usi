package usi

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/kk-no/go-usi/usicmd"
)

type USI struct {
	enginePath string
	engineName string

	status  Status
	process ReadWriteProcessor
	result  ResultManager
}

func New(enginePath string, options ...Option) (*USI, error) {
	absPath, err := filepath.Abs(enginePath)
	if err != nil {
		return nil, err
	}
	i := strings.LastIndex(absPath, string(filepath.Separator))

	usi := &USI{
		enginePath: absPath[:i],
		engineName: "." + absPath[i:],

		status: Disconnected,
	}

	for _, o := range options {
		o(usi)
	}

	return usi, nil
}

func (u *USI) SetStatus(to Status) {
	u.status = to
}

func (u *USI) Connect(ctx context.Context) error {
	u.SetStatus(WaitConnecting)

	if err := os.Chdir(u.enginePath); err != nil {
		return err
	}

	p, err := NewReadWriteProcessor(ctx, u.engineName)
	if err != nil {
		return err
	}

	u.process = p
	u.result = NewResultManager()

	u.process.Start(ctx)
	u.SetStatus(Connected)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in := scanner.Text()
		if usicmd.IsCommand(in) {
			if usicmd.IsQuit(in) {
				return errors.New("quit inputted")
			}
			u.SendCommand(in)
		}
	}
	return scanner.Err()
}

func (u *USI) Disconnect() error {
	if err := u.process.Stop(); err != nil {
		return err
	}
	u.process = nil

	u.SetStatus(Disconnected)
	return nil
}

func (u *USI) IsConnected() bool {
	return u.process != nil
}

func (u *USI) SendCommand(command string) {
	u.process.SendCommand(usicmd.Command(command))
}

func (u *USI) HandleMessage(message string) {
	u.result.ReceiveMessage(message)

	var token string
	if index := strings.Index(message, " "); index == -1 {
		token = message
	} else {
		token = message[0:index]
	}

	switch Response(token) {
	case ReadyOK:
		u.SetStatus(WaitCommand)
	case BestMove:
		u.result.HandleBestMove(message)
		u.SetStatus(WaitCommand)
	case Info:
		u.result.HandleInfo(message)
	}
}
