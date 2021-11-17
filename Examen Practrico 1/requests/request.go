package requests

const (
	HELLO int = iota
	TXT
	FILE
	GOODBYE
)

type Request struct {
	Id     int
	Sender string
	Info   string
	Data   []byte
}
