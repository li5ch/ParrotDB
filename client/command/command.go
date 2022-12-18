package command

// Command represents a command from a client to an NSQ daemon
type Command struct {
	Name   []byte
	Params [][]byte
	Body   []byte
}

func Set(key, value string) (*Command, error) {
	params := [][]byte{[]byte(key)}
	if len(value) > 0 {
		params = append(params, []byte(value))
	}
	return &Command{[]byte("SET"), params, nil}, nil
}

