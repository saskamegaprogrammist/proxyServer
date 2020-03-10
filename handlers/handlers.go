package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/saskamegaprogrammist/proxyServer/requests"
	"github.com/saskamegaprogrammist/proxyServer/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var RepeaterClient *http.Client

func GetLastRequests (writer http.ResponseWriter, req *http.Request) {
	var requestModel requests.Request
	requestsRetrieved, err := requestModel.GetRequests()
	if err != nil {
		utils.CreateAnswer(writer, 404, requests.CreateError(err.Error()))
		return
	}
	utils.CreateAnswer(writer, 200, requestsRetrieved)
}

func MakeRequest (writer http.ResponseWriter, req *http.Request) {
	var requestModel requests.Request
	reqId := mux.Vars(req)["id"]
	reqIdInt, err := strconv.Atoi(reqId)
	if err != nil {
		utils.CreateAnswer(writer, 404, requests.CreateError(err.Error()))
		return
	}
	err = requestModel.GetRequest(reqIdInt)
	if err != nil {
		utils.CreateAnswer(writer, 404, requests.CreateError(err.Error()))
		return
	}
	request, err := http.NewRequest(requestModel.Method, fmt.Sprintf("%s://%s", requestModel.URLscheme, requestModel.URLhost), strings.NewReader(requestModel.Body))
	if err != nil {
		utils.CreateAnswer(writer, 404, requests.CreateError(err.Error()))
		return
	}
	response, err := RepeaterClient.Do(request)
	if err != nil {
		utils.CreateAnswer(writer, 404, requests.CreateError(err.Error()))
		return
	}
	utils.CopyHeader(writer.Header(), response.Header)
	writer.WriteHeader(response.StatusCode)
	io.Copy(writer, response.Body)
}
