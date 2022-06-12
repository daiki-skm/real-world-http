package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

func reqReWrite() {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = ":9001"
	}
	rp := &httputil.ReverseProxy{
		Director: director,
	}
	server := http.Server{
		Addr:    "127.0.0.1:9000",
		Handler: rp,
	}
	log.Println("Start listening at :9000")
	log.Fatalln(server.ListenAndServe())
}

func resReWrite() {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = ":9001"
	}
	modifier := func(res *http.Response) error {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Reading body error: %w", err)
		}
		newBody := bytes.NewBuffer(body)
		newBody.WriteString(" via Proxy")
		res.Body = ioutil.NopCloser(newBody)
		res.Header.Set("Content-Length", strconv.Itoa(newBody.Len()))
		return nil
	}
	rp := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifier,
	}
	server := http.Server{
		Addr:    "127.0.0.1:9000",
		Handler: rp,
	}
	log.Println("Start listening at :9000")
	log.Fatalln(server.ListenAndServe())
}

type RetryTransport struct {
}

func (RetryTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	for i := 0; i < 3; i++ {
		resp, err = http.DefaultTransport.RoundTrip(req)
		if err != nil {
			log.Println("fail")
			time.Sleep(time.Second)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("failed to request to %s", req.URL.String())
}

func retry() {
	target, _ := url.Parse("http://127.0.0.1:9001")
	rp := httputil.NewSingleHostReverseProxy(target)
	rp.Transport = &RetryTransport{}
	server := http.Server{
		Addr:    "127.0.0.1:9000",
		Handler: rp,
	}
	log.Println("Start listening at :9000")
	log.Fatalln(server.ListenAndServe())
}

func main() {
	//reqReWrite()
	//resReWrite()
	//retry()
}
