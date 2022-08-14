package usicmd

import "strings"

type Command string

const (
	Go       Command = "go"
	Stop     Command = "stop"
	Position Command = "position"
	Quit     Command = "quit"
	Moves    Command = "moves"
	Side     Command = "side"
	IsReady  Command = "isready"
	NewGame  Command = "usinewgame"
	GameOver Command = "gameover"
)

func IsCommand(command string) bool {
	switch Command(strings.Split(command, " ")[0]) {
	case Go, Stop, Position, Quit, Moves, Side, IsReady, NewGame, GameOver:
		return true
	default:
		return false
	}
}

func IsQuit(command string) bool {
	switch Command(command) {
	case Quit:
		return true
	default:
		return false
	}
}
