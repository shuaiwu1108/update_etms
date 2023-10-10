package model

type Application struct {
	Datasource Datasource
}

type Datasource struct {
	DriverName string // oracle or mysql
	Host       string // localhost or  192.168.x.x
	Port       int    // 1521
	Database   string // orcl or xxx
	Username   string
	Password   string
}
