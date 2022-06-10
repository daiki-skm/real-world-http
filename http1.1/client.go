package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/rpc/jsonrpc"
	"strconv"
	"time"
)

func GET() {
	resp, err := http.Get("https://localhost:18443")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}

func GETWith() {
	cert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		panic(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}
	tlsConfig.BuildNameToCertificate()
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	resp, err := client.Get("https://localhost:18443")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}

func GETWithClient() {
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	resp, err := client.Get("https://localhost:18443")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}

func GETWithUpgrade() {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	conn, err := dialer.Dial("tcp", "localhost:18888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	request, _ := http.NewRequest("GET", "http://localhost:18888/upgrade", nil)
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "MyProtocol")
	err = request.Write(conn)
	if err != nil {
		panic(err)
	}

	response, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", response.Status)
	log.Println("Headers:", response.Header)

	counter := 10
	for {
		data, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		fmt.Println("<-", string(bytes.TrimSpace(data)))
		fmt.Fprintf(conn, "%d\n", counter)
		fmt.Println("->", counter)
		counter--
	}
}

func GETWithChunk() {
	resp, err := http.Get("https://localhost:18888/chunked")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		log.Println(string(bytes.TrimSpace(line)))
	}
}

func GETWithChunkComplete() {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	conn, err := dialer.Dial("tcp", "localhost:18888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	request, err := http.NewRequest("GET", "http://localhost:18888/chunked", nil)
	err = request.Write(conn)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}
	if resp.TransferEncoding[0] != "chunked" {
		panic("wrong transfer encoding")
	}
	for {
		sizeStr, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		size, err := strconv.ParseInt(string(sizeStr[:len(sizeStr)-2]), 16, 64)
		if size == 0 {
			break
		}
		if err != nil {
			panic(err)
		}
		line := make([]byte, int(size))
		reader.Read(line)
		reader.Discard(2)
		log.Println(" ", string(line))
	}
}

func RPCClient() {
	client, err := jsonrpc.Dial("tcp", "localhost:18888")
	if err != nil {
		panic(err)
	}
	var result int
	args := &Args{4, 5}
	err = client.Call("Arith.Multiply", args, &result)
	if err != nil {
		panic(err)
	}
	log.Printf("%d*%d=%d\n", args.A, args.B, result)
}

func main() {
	//GET()
	//GETWith()
	//GETWithClient()
	//GETWithUpgrade()
	//GETWithChunk()
	//GETWithChunkComplete()
	//RPCClient()
}
