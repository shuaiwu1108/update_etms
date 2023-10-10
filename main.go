package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"update_etms/db"
	"update_etms/handler"
	"update_etms/model"
	"update_etms/util"
)

var c model.Application
var conn *sql.DB

var (
	Logger      = logrus.New()
	wd, _       = os.Getwd()
	logFileName = "upEtms.log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	//读取配置文件，连接orcl
	conn = db.ConnDB(c)
	defer func() {
		conn.Close()
		if err := recover(); err != nil {
			handler.Logger.Infof("程序异常， %v", err)
		}
	}()

	src, err := os.OpenFile(path.Join(wd, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer src.Close()
	if err != nil {
		fmt.Println("err", err)
	}
	r := gin.Default()
	r.Use(handler.LogerMiddleware(src, logFileName))
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "upEtms succ",
		})
	})
	r.POST("/searchorcl", func(c *gin.Context) {
		var body interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		bodyMap := body.(map[string]interface{})
		sql := bodyMap["sql"].(string)
		handler.Logger.Infof("查询sql:%v", sql)

		// 使用连接，查询orcl
		listRes := util.CustomQuery(conn, sql)
		bodyMap["resData"] = &listRes

		// 返回结果集
		c.JSON(http.StatusOK, bodyMap)
	})
	r.Run(":9005")
}

func init() {
	InitConfig()
}

func InitConfig() {
	path := "./application.conf"
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Println(fmt.Errorf("%v配置文件读取错误, %v", path, err))
		path = "/home/cty/BusSale/conf/application.conf"
		handler.Logger.Infof("开始读取%v\n", path)
		f, err = os.Open(path)
		if err != nil {
			log.Println(fmt.Errorf("%v配置文件读取错误, %v", path, err))
			path = "/data/cty/BusSale/conf/application.conf"
			handler.Logger.Infof("开始读取%v\n", path)
			f, err = os.Open(path)
			if err != nil {
				log.Println(fmt.Errorf("%v配置文件读取错误, %v", path, err))
				path = "/data/BusSale/conf/application.conf"
				handler.Logger.Infof("开始读取%v\n", path)
				f, err = os.Open(path)
				if err != nil {
					panic(fmt.Errorf("配置文件读取错误, %v", err))
				}
			}
		}
	}
	scann := bufio.NewScanner(f)
	stationConfig := make(map[string]string)

	re := regexp.MustCompile(`^db_station\d?\.(?P<key>.*)=(?P<value>.*)`)
	for scann.Scan() {
		line := scann.Text()
		if strings.Contains(line, "db_station") {
			handler.Logger.Infof("数据：%v", line)
			match := re.FindStringSubmatch(line)
			stationConfig[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
		}
	}
	url := strings.Split(stationConfig["url"], ":")[3][1:len(strings.Split(stationConfig["url"], ":")[3])]
	port, _ := strconv.Atoi(strings.Split(stationConfig["url"], ":")[4])
	database := strings.Split(stationConfig["url"], ":")[5]
	d := model.Datasource{
		DriverName: "oracle",
		Host:       url,
		Port:       port,
		Database:   database,
		Username:   stationConfig["user"],
		Password:   stationConfig["pass"],
	}
	c = model.Application{
		Datasource: d,
	}
}
