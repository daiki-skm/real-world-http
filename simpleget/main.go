package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/idna"
)

func GET() {
	resp, err := http.Get("http://localhost:18888")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Println("Status:", resp.Status)
	log.Println("StatusCode:", resp.StatusCode)
	log.Println("Header:", resp.Header)
	log.Println("Content-Length:", resp.Header.Get("Content-Length"))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
}

func GET_Q() {
	values := url.Values{
		"query": {"hello world"},
	}
	resp, err := http.Get("http://localhost:18888?" + values.Encode())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
}

func HEAD() {
	resp, err := http.Head("http://localhost:18888")
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
	log.Println("StatusCode:", resp.StatusCode)
	log.Println("Header:", resp.Header)
}

func POST() {
	values := url.Values{
		"test": {"value"},
	}
	resp, err := http.PostForm("http://localhost:18888", values)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
}

func POST_OP() {
	//file, err := os.Open("main.go")
	file := strings.NewReader("hello world")
	//if err != nil {
	//	panic(err)
	//}
	resp, err := http.Post("http://localhost:18888", "text/plain", file)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
	log.Println(string(body))
}

func POST_MULT() {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("test", "value")

	fileWriter, err := writer.CreateFormFile("thumbnail", "thx.jpg")
	if err != nil {
		panic(err)
	}
	readFile, err := os.Open("thx.jpg")
	if err != nil {
		panic(err)
	}
	defer readFile.Close()
	io.Copy(fileWriter, readFile)

	writer.Close()
	resp, err := http.Post("http://localhost:18888", writer.FormDataContentType(), &buffer)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
}

func POST_MULT_MIME() {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("test", "value")

	part := make(textproto.MIMEHeader)
	part.Set("Content-Type", "image/jpeg")
	part.Set("Content-Disposition", "form-data; name=\"thumbnail\"; filename=\"thx.jpg\"")
	fileWriter, err := writer.CreatePart(part)
	if err != nil {
		panic(err)
	}
	readFile, err := os.Open("thx.jpg")
	if err != nil {
		panic(err)
	}
	io.Copy(fileWriter, readFile)

	writer.Close()
	resp, err := http.Post("http://localhost:18888", writer.FormDataContentType(), &buffer)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
}

func COOKIE() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	//client := http.Client{Jar: jar}
	http.DefaultClient = &http.Client{Jar: jar}
	for i := 0; i < 2; i++ {
		//resp, err := client.Get("http://localhost:18888/cookie")
		resp, err := http.Get("http://localhost:18888/cookie")
		if err != nil {
			panic(err)
		}
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			panic(err)
		}
		log.Println(string(dump))
	}
}

func PROXY() {
	proxyUrl, err := url.Parse("http://localhost:18888")
	if err != nil {
		panic(err)
	}
	client := http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	resp, err := client.Get("http://github.com")
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}

func LOCAL() {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	client := http.Client{Transport: transport}
	resp, err := client.Get("file://./main.go")
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}

func DELETE() {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://localhost:18888", nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}

func international() {
	src := "大空"
	ascii, err := idna.ToASCII(src)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s -> %s\n", src, ascii)
}

func main() {
	//GET()
	//GET_Q()
	//HEAD()
	//POST()
	//POST_OP()
	//POST_MULT()
	//POST_MULT_MIME()
	//COOKIE()
	//PROXY()
	//LOCAL()
	//DELETE()
	international()
}
