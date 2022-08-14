package usi

type Status string

const (
	WaitConnecting Status = "WaitConnecting"
	Connected      Status = "Connected"
	WaitReadyOK    Status = "WaitReadyOk"
	WaitCommand    Status = "WaitCommand"
	WaitBestMove   Status = "WaitBestMove"
	WaitOneLine    Status = "WaitOneLine"
	Disconnected   Status = "Disconnected"
)

type Response string

const (
	ReadyOK  Response = "readyok"
	BestMove Response = "bestmove"
	Info     Response = "info"
)
