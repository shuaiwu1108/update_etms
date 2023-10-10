package db

import (
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	goora "github.com/sijms/go-ora/v2"
	"update_etms/handler"
	"update_etms/model"
)

func ConnDB(c model.Application) *sql.DB {
	handler.Logger.Infof("orcl连接初始化...")
	url := goora.BuildUrl(c.Datasource.Host, c.Datasource.Port, c.Datasource.Database, c.Datasource.Username, c.Datasource.Password, nil)
	db, err := sql.Open(c.Datasource.DriverName, url)
	if err != nil {
		panic(fmt.Errorf("Unmarshal error config file: %s \n", err))
	}
	handler.Logger.Infof("orcl连接初始化完毕")
	return db
}
