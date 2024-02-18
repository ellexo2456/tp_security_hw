package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Server struct{}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleHTTP(w, r)
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
		log.Fatal(err.Error())
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
		log.Fatal(err.Error())
	}

	for h, v := range from.Header {
		for _, v := range v {
			to.Header().Add(h, v)
		}
	}
}
