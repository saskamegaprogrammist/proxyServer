package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/saskamegaprogrammist/proxyServer/certificate"
	"github.com/saskamegaprogrammist/proxyServer/db"
	"github.com/saskamegaprogrammist/proxyServer/requests"
	"github.com/saskamegaprogrammist/proxyServer/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var rootCertificate certificate.Cert

func saveToDB(req *http.Request, scheme string) {
	var reqModel requests.Request
	reqModel.Method = req.Method
	reqModel.URLhost = req.URL.Host
	reqModel.URLscheme = scheme
	header := make(map[string]string, 0)
	for k, v := range req.Header {
		header[k] = v[0]
	}
	reqModel.Header = header
	body,_ := ioutil.ReadAll(req.Body)
	bodyString := string(body)
	reqModel.Body = bodyString
	reqModel.ContentLength = int(req.ContentLength)
	reqModel.Host = req.Host
	reqModel.RemoteAddr = req.RemoteAddr
	reqModel.RequestURI = req.RequestURI
	err := reqModel.SaveRequest()
	if err != nil {
		log.Println(err)
	}
}

func copyBufferToDB(buffer *bytes.Buffer, requestHost string) {
	p := make([]byte, 1024*1024*8)

	for {
		n, err := buffer.Read(p)
		if err != nil{
			if err == io.EOF {
				fmt.Println(string(p[:n])) //should handle any remaining bytes.
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(p[:n]))
	}
	reader:=bufio.NewReader(bytes.NewReader(p))
	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Println(err)
		return
	} else {
		req.URL.Host = requestHost
		saveToDB(req, "https")
	}
}


func transfer(destination io.WriteCloser, source io.ReadCloser, copy bool, requestHost string) {
	defer destination.Close()
	defer source.Close()
	if copy {
		buffer := &bytes.Buffer{}
		duplicateSources := io.MultiWriter(destination, buffer) //we copy data from source into buffer and destination
		io.Copy(duplicateSources, source)
		go copyBufferToDB(buffer,  requestHost)

	} else {
		io.Copy(destination, source)
	}
}

func handleCONNECT(writer http.ResponseWriter, req *http.Request) {
	log.Println(req)
	cert, err := certificate.CreateLeafCertificate(req.Host)
	if err != nil {
		log.Println(err)
	}
	tlsConfig := & tls.Config{
		Certificates:                []tls.Certificate{*cert},
		GetCertificate:              func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return certificate.CreateLeafCertificate(info.ServerName)
		},
	}
	destinationConnection, err := tls.Dial("tcp", req.Host, tlsConfig)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusServiceUnavailable)
		return
	}
	hijacker, ok := writer.(http.Hijacker)
	if !ok {
		http.Error(writer, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConnection, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusServiceUnavailable)
	}
	_, err = clientConnection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	tlsConnection := tls.Server(clientConnection, tlsConfig)
	err = tlsConnection.Handshake()

	go transfer(destinationConnection, tlsConnection, true, req.URL.Host)
	go transfer(tlsConnection, destinationConnection, false, req.URL.Host)
}

func handleHTTPRequests(writer http.ResponseWriter, req *http.Request) {
	saveToDB(req, "http")
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Printf("proxy error: %s, request: %+v", err.Error(), req)
		http.Error(writer, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Printf("status code: %d", resp.StatusCode)
	defer resp.Body.Close()
	utils.CopyHeader(writer.Header(), resp.Header)
	writer.WriteHeader(resp.StatusCode)
	io.Copy(writer, resp.Body)

}

func handleRequests(writer http.ResponseWriter, req *http.Request) {
	log.Println(req.Method)
	if req.Method == "CONNECT" {
		handleCONNECT(writer, req)
	} else {
		handleHTTPRequests(writer, req)
	}
}


func main() {

	db.CreateDataBaseConnection("docker", "docker", "localhost", "docker", 20)
	//db.CreateDataBaseConnection("alexis", "sinope27", "localhost", "proxy", 20)
	db.InitDataBase()

	server := &http.Server{
		Handler:      http.HandlerFunc(handleRequests),
		Addr:         ":5000",
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	rootCertificate = certificate.GetRootCertificate()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}