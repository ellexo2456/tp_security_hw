package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"slices"
)

var certPool []string

type Server struct{}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "CONNECT" {
		handleHTTPS(w, r)
	} else {
		handleHTTP(w, r)
	}
}

func handleHTTPS(w http.ResponseWriter, r *http.Request) {
	connDest, err := connectHandshake(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	connSrc, _, err := hijacker.Hijack()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = connSrc.Write([]byte("HTTP/1.0 200 Connection established\\r\\n\\r\\n"))
	if err != nil {
		log.Fatal(err)
		return
	}

	go exchangeData(connDest, connSrc)
	go exchangeData(connSrc, connDest)
}

func exchangeData(to io.WriteCloser, from io.ReadCloser) {
	defer func() {
		to.Close()
		from.Close()
	}()

	_, err := io.Copy(to, from)
	if err != nil {
		fmt.Println(err)
	}
}

func connectHandshake(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	if !slices.Contains(certPool, r.Host) {
		if err := generateCert(r.Host); err != nil {
			return nil, err
		}
	}

	cert, err := os.ReadFile(r.Host)
	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(cert)
	if !ok {
		return nil, errors.New("failed to parse root certificate")
	}

	conn, err := tls.Dial("tcp", r.Host, &tls.Config{
		ClientCAs: roots,
	})
	if err != nil {
		return nil, err
	}

	w.WriteHeader(http.StatusOK)
	return conn, nil
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")
	r.RequestURI = ""

	c := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	targetResponse, err := c.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(targetResponse.Body)

	copyResponse(targetResponse, w)
}

func copyResponse(from *http.Response, to http.ResponseWriter) {
	to.WriteHeader(from.StatusCode)

	_, err := io.Copy(to, from.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for h, v := range from.Header {
		for _, v := range v {
			to.Header().Add(h, v)
		}
	}
}
