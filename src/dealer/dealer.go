package dealer

import (
	"bufio"
	"context"
	"crypto/tls"
	"github.com/ellexo2456/tp_security_hw/src/utils"
	"net"
	"net/http"
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

	return http.ReadResponse(bufio.NewReader(conn), req)
}

func WriteResponse(conn net.Conn, resp *http.Response) error {
	err := resp.Write(conn)
	if err != nil {
		return err
	}

	return nil
}
