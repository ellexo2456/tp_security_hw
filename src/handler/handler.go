package handler

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/ellexo2456/tp_security_hw/src/dealer"
	"github.com/ellexo2456/tp_security_hw/src/utils"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ellexo2456/tp_security_hw/src/certs"
)

type Handler struct {
	certs map[string][]byte
	mu    sync.RWMutex
}

func New() (*Handler, error) {
	cs, err := certs.Load()
	if err != nil {
		return nil, err
	}

	return &Handler{
		certs: cs,
		mu:    sync.RWMutex{},
	}, nil
}

func (h *Handler) Handle(source net.Conn) error {
	req, err := http.ReadRequest(bufio.NewReader(source))
	if err != nil {
		return err
	}

	dest, source, err := h.makeConnection(req, source)
	if err != nil {
		return err
	}
	defer func(dest net.Conn) {
		err := dest.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(dest)

	req.Header.Del("Proxy-Connection")
	req.RequestURI = req.URL.Path
	return h.makeExchange(source, dest, req)
}

func (h *Handler) makeExchange(source, dest net.Conn, req *http.Request) error {
	resp, err := dealer.SendRequest(dest, req)
	if err != nil {
		return err
	}

	err = dealer.WriteResponse(source, resp)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) makeConnection(req *http.Request, source net.Conn) (net.Conn, net.Conn, error) {
	host := req.URL.Hostname()
	port := utils.GetPort(req)

	if req.Method != "CONNECT" {
		return h.makeTcp(req, source, host, port)
	}
	if req.Method == "CONNECT" {
		return h.makeTls(req, source, host, port)
	}

	return nil, nil, nil
}

func (h *Handler) makeTcp(req *http.Request, src net.Conn, host, port string) (net.Conn, net.Conn, error) {
	req.URL.Scheme = "http"
	dest, err := dealer.TcpConnect(host, port)
	if err != nil {
		return nil, nil, err
	}

	return dest, src, nil
}

func (h *Handler) makeTls(req *http.Request, source net.Conn, host, port string) (net.Conn, net.Conn, error) {
	_, err := source.Write([]byte("HTTP/1.0 200 Connection established\n\n"))
	if err != nil {
		return nil, nil, err
	}

	source, err = h.upgradeToTls(source, host)
	if err != nil {
		return nil, nil, err
	}

	r, err := http.ReadRequest(bufio.NewReader(source))
	if err != nil {
		return nil, nil, err
	}
	*req = *r

	//req.URL.Scheme = "https"
	dest, err := dealer.TlsConnect(host, port)
	if err != nil {
		return nil, nil, err
	}

	return dest, source, nil
}

func (h *Handler) upgradeToTls(conn net.Conn, host string) (net.Conn, error) {
	cfg, err := h.getTlsConfig(host)
	if err != nil {
		return nil, err
	}

	tlsConn := tls.Server(conn, cfg)
	err = tlsConn.SetReadDeadline(time.Now().Add(utils.DefaultTimeout))
	if err != nil {
		return nil, err
	}

	return tlsConn, nil
}

func (h *Handler) getTlsConfig(host string) (*tls.Config, error) {
	if err := h.addCert(host); err != nil {
		return nil, err
	}

	key, err := os.ReadFile("./src/certs/cert.key")
	if err != nil {
		return nil, err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	c, err := tls.X509KeyPair(h.certs[host], key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{c},
	}, nil
}

func (h *Handler) addCert(host string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.certs[host]; exists {
		return nil
	}

	c, err := certs.Generate(host)
	if err != nil {
		return err
	}

	h.certs[host] = c
	return nil
}
