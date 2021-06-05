package src

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
)

var Maria *sql.DB
var QueryAPILog = `INSERT INTO osu.api_log (time, request_id, remote_ip, host, method, uri, user_agent, status, error, latency, latency_human, bytes_in, bytes_out) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?);`

func QueryOnly(sql string, parm ...interface{}) error {

	raws, err := Maria.Query(sql, parm...)
	if err != nil {
		if raws != nil {
			_ = raws.Close()
		}
		return err
	}

	return raws.Close()
}

func ConnectMaria() {

	db, err := sql.Open("mysql", Settings.Config.Sql.Id+":"+Settings.Config.Sql.Passwd+"@tcp("+Settings.Config.Sql.Url+")/")
	if Maria = db; db != nil {
		Maria.SetMaxOpenConns(100)
		fmt.Println("mariaDB connected")
	} else {
		panic(err)
	}
}

func Upsert(query string, data []interface{}) {
	data = append(data, data[1:]...)
	err := QueryOnly(
		query,
		data...,
	)
	if err != nil {
		fmt.Println(err)

	}
}

func ToDateTime(t interface{}) string {
	if t == nil {
		return "0000-00-00T00:00:00"
	}
	myDate, _ := time.Parse("2006-01-02T15:04:05-07:00", t.(string))
	return myDate.Format("2006-01-02T15:04:05")
}

func InsertAPILog(s ...interface{}) (err error) {
	rows, err := Maria.Query(QueryAPILog, s...)
	if err != nil {
		return
	}
	defer rows.Close()
	return
}
