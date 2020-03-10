package main

import (
	"crypto/tls"
	"github.com/saskamegaprogrammist/proxyServer/certificate"
	"github.com/saskamegaprogrammist/proxyServer/db"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var rootCertificate certificate.Cert


func copyHeader(dst, src http.Header) []string {
	header := make([]string, 0)
	for k, vv := range src {
		header = append(header, k)
		var value string
		for _, v := range vv {
			value = value + v
			dst.Add(k, v)
		}
		header = append(header, value)
	}
	return header
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	var a []byte
	source.Read(a)
	log.Println(a)
	io.Copy(destination, source)
}

func handleCONNECT(writer http.ResponseWriter, req *http.Request) {
	log.Println(req.Host)
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

	go transfer(destinationConnection, tlsConnection)
	go transfer(tlsConnection, destinationConnection)
}

func handleHTTPRequests(writer http.ResponseWriter, req *http.Request) {
	//req.URL, _ = url.ParseRequestURI(req.RequestURI)
	log.Println(req)
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Printf("proxy error: %s, request: %+v", err.Error(), req)
		http.Error(writer, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Printf("status code: %d", resp.StatusCode)
	defer resp.Body.Close()
	headerString := copyHeader(writer.Header(), resp.Header)
	writer.WriteHeader(resp.StatusCode)
	io.Copy(writer, resp.Body)

	dataBase := db.GetDataBase()
	transaction, _ := dataBase.Begin()
	body,_ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	log.Println(bodyString)
	_, err = transaction.Exec("INSERT INTO requests (headers) VALUES ($1) ", headerString)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
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