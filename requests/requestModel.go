package requests

import (
	"fmt"
	"github.com/saskamegaprogrammist/proxyServer/db"
	"log"
	"strings"
)

type Request struct {
	Id int `json:"id"`
	Method string`json:"method"`
	URLhost string `json:"host"`
	URLscheme string `json:"scheme"`
	Header map[string]string `json:"header"`
	Body string`json:"body"`
	ContentLength int `json:"contentlength"`
	Host string `json:"host"`
	RemoteAddr string `json:"remoteaddr"`
	RequestURI string `json:"requestURI"`
}


func (reqModel*Request) SaveRequest() error {
	dataBase := db.GetDataBase()
	transaction, _ := dataBase.Begin()
	header := make([]string, 0)
	for k, v := range reqModel.Header {
		header = append(header, fmt.Sprintf("%s : %s", k, v))
	}
	_, err := transaction.Exec("INSERT INTO requests (method, urlhost, urlscheme, headers, body, contentlength, host, remoteaddr, requesturi) VALUES  ($1, $2, $3, $4, $5, $6, $7, $8, $9) ",
		reqModel.Method, reqModel.URLhost, reqModel.URLscheme, header, reqModel.Body, reqModel.ContentLength, reqModel.Host,
		reqModel.RemoteAddr, reqModel.RequestURI)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
			return err
		}
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return  nil
}

func (reqModel*Request) GetRequests() ([]Request, error) {
	requestsFound := make([]Request, 0)
	dataBase := db.GetDataBase()
	transaction, _ := dataBase.Begin()
	rows, err := transaction.Query("SELECT id, method, urlhost, urlscheme, headers::text[], body, contentlength, host, remoteaddr, requesturi FROM requests ORDER BY id DESC LIMIT 10;")
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return requestsFound, err
	}
	for rows.Next() {
		var reqRetrieved Request
		var headerArray []string
		err = rows.Scan(&reqRetrieved.Id, &reqRetrieved.Method, &reqRetrieved.URLhost, &reqRetrieved.URLscheme, &headerArray,
			&reqRetrieved.Body, &reqRetrieved.ContentLength, &reqRetrieved.Host, &reqRetrieved.RemoteAddr, &reqRetrieved.RequestURI)
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if err != nil {
				log.Fatalln(errRollback)
			}
			return requestsFound, err
		}
		reqRetrieved.Header = getHeader(headerArray)
		requestsFound = append(requestsFound, reqRetrieved)
	}

	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return requestsFound, nil
}


func (reqModel*Request) GetRequest(id int) error {
	dataBase := db.GetDataBase()
	transaction, _ := dataBase.Begin()
	var headerArray []string
	row := transaction.QueryRow("SELECT id, method, urlhost, urlscheme, headers::text[], body, contentlength, host, remoteaddr, requesturi FROM requests WHERE id = $1", id)
	err := row.Scan(&reqModel.Id, &reqModel.Method, &reqModel.URLhost, &reqModel.URLscheme, &headerArray,
		&reqModel.Body, &reqModel.ContentLength, &reqModel.Host, &reqModel.RemoteAddr, &reqModel.RequestURI)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf("can't find request with id %d", id)
	}
	reqModel.Header = getHeader(headerArray)
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func getHeader(headerDB []string) map[string]string {
	headerMap := make(map[string]string, 0)
	for _, header := range headerDB {
		keyVal := strings.Split(header, " : ")
		headerMap[keyVal[0]] = keyVal[1]
	}
	return headerMap
}