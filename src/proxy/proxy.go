package proxy

import (
	"fmt"
	"net"
)

func StartServer(port int, handler func(conn net.Conn) error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: port,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func() {
			defer conn.Close()

			if err := handler(conn); err != nil {
				fmt.Println(err)
			}
		}()

	}
}
