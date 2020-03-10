package utils

import (
	"encoding/json"
	"net/http"
)

func CreateAnswer(w http.ResponseWriter, statusCode int, value interface{})  {
	encoded , err := json.Marshal(value)
	if err != nil {
		//log.Println(err)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(encoded)
	if err != nil {
		//log.Println(err)
	}
}

func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		var value string
		for _, v := range vv {
			value = value + v
			dst.Add(k, v)
		}
	}
}
