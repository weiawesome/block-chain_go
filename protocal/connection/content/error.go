package content

import (
	"encoding/json"
	"net"
)

type ErrorMessage struct {
	Error string `json:"error"`
}

func SentErrorMessage(conn net.Conn, message string) {
	response, err := json.Marshal(ErrorMessage{Error: message})
	if err != nil {
		//fmt.Println("Error found ", err.Error(), " in SEC")
	}
	SendContent(conn, string(response))
}
