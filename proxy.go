package main

import (
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
)

var rootCertificate certificate.Cert

func saveToDB(req *http.Request) {
	var reqModel requests.Request
	reqModel.Method = req.Method
	reqModel.URLhost = req.URL.Host
	reqModel.URLscheme = req.URL.Scheme
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


func transfer(destination io.WriteCloser, source io.ReadCloser, copy bool) {
	defer destination.Close()
	defer source.Close()
	if copy {
		var buffer []byte
		io.ReadFull(source,buffer)
		fmt.Println("buffer", buffer)
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

	go transfer(destinationConnection, tlsConnection, true)
	go transfer(tlsConnection, destinationConnection, false)
}

func handleHTTPRequests(writer http.ResponseWriter, req *http.Request) {
	saveToDB(req)
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
	db.CreateDataBaseConnection("alexis", "sinope27", "localhost", "proxy", 20)
	db.InitDataBase()

	server := &http.Server{
		Handler:      http.HandlerFunc(handleRequests),
		Addr:         "127.0.0.1:5000",
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	rootCertificate = certificate.GetRootCertificate()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}