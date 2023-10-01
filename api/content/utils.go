package content

import (
	"net"
)

func SendContent(conn net.Conn, val string) {
	responseBytes := []byte(string(val) + SuffixString)

	_, err := conn.Write(responseBytes)
	if err != nil {
		//fmt.Println("Error found ", err.Error(), " in SC")
		return
	}
}
