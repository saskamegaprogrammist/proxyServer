package handlers

import (
	"github.com/saskamegaprogrammist/proxyServer/utils"
	"net/http"
)

func GetLastRequests(writer http.ResponseWriter, req *http.Request) {
	&http.Request{
		Method:           "",
		URL:              nil,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             nil,
		GetBody:          nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Host:             "",
		Form:             nil,
		PostForm:         nil,
		MultipartForm:    nil,
		Trailer:          nil,
		RemoteAddr:       "",
		RequestURI:       "",
		TLS:              nil,
		Cancel:           nil,
		Response:         nil,
	}


	utils.CreateAnswer(writer, 200, )
}
