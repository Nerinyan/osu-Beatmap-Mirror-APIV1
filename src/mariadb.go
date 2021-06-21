package src

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nerina1241/osu-beatmap-mirror-api/ConsoleLogger"
	"github.com/nerina1241/osu-beatmap-mirror-api/Settings"
)

var Maria *sql.DB

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
		ConsoleLogger.Consolelog("DBSM", "Succesfully connected DBSM Server.")
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
		ConsoleLogger.WarningConsolelog("Warning", err.Error())

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
		ConsoleLogger.WarningConsolelog("Logger", err.Error())
		return
	}
	defer rows.Close()
	return
}
