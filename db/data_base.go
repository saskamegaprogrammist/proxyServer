package db

import (
	"fmt"
	"github.com/jackc/pgx"
	"log"
)

var dataBasePool *pgx.ConnPool

func CreateAddress(user, password, host, name string) string {
	return  fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=%s",
	user, password, host, name)
}

func CreateDataBaseConnection(user, password, host, name string, maxConn int) {
	dataBaseConfig := CreateAddress(user, password, host, name)
	connectionConfig, err := pgx.ParseConnectionString(dataBaseConfig)
	if err != nil {
		//log.Println(err);
		return
	}
	dataBasePool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: connectionConfig,
		MaxConnections: maxConn,
	})
	if err != nil {
		//log.Println(err);
		return
	}
}


func InitDataBase() {
	_, err := dataBasePool.Exec(`
	DROP TABLE IF EXISTS requests;
	
	CREATE TABLE IF NOT EXISTS requests (
    id SERIAL NOT NULL PRIMARY KEY,
    method text,
    urlhost text,
    urlscheme text,
    headers text,
    body text,
    contentlength int,
    host text,
    remoteaddr text,
    requesturi text
);`)
	 if err != nil {
		log.Println(err)
	}

}

func GetDataBase() *pgx.ConnPool {
	return dataBasePool
}