package dealer

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/ellexo2456/tp_security_hw/src/utils"
	"net"
	"net/http"
	"os"
)

func TlsConnect(host, port string) (net.Conn, error) {
	d := tls.Dialer{}

	ctx, c := context.WithTimeout(context.Background(), utils.DefaultTimeout)
	defer c()

	return d.DialContext(ctx, "tcp", host+":"+port)
}

func TcpConnect(host, port string) (net.Conn, error) {
	d := net.Dialer{}

	ctx, c := context.WithTimeout(context.Background(), utils.DefaultTimeout)
	defer c()

	return d.DialContext(ctx, "tcp", host+":"+port)
}

func SendRequest(conn net.Conn, req *http.Request) (*http.Response, error) {
	err := req.Write(conn)
	if err != nil {
		return nil, err
	}

	fmt.Print("\n\n\n\n")
	fmt.Println("###############################################################")
	fmt.Println("###############################################################")
	fmt.Println("############################REQUEST############################")
	fmt.Println("###############################################################")
	fmt.Println("###############################################################")
	req.Write(os.Stdout)
	fmt.Print("\n\n\n\n")
	return http.ReadResponse(bufio.NewReader(conn), req)
}

func WriteResponse(conn net.Conn, resp *http.Response) error {
	err := resp.Write(conn)
	if err != nil {
		return err
	}
	fmt.Print("\n\n\n\n")
	fmt.Println("###############################################################")
	fmt.Println("###############################################################")
	fmt.Println("############################RESPONSE############################")
	fmt.Println("###############################################################")
	fmt.Println("###############################################################")
	resp.Write(os.Stdout)
	fmt.Print("\n\n\n\n")

	return nil
}
