package database

import (
	"database/sql"
	"fmt"
	"strings"
	"strconv"
	"time"
)

/**
初始化连接池(`sql.DB`对象)，及设置各种参数
基于驱动:go-sql-driver/mysql
*/

const(
	defaultDriver = "mysql"  //默认数据库驱动
	defaultPort = 3306  		//默认数据库端口
	defaultFormat = "%s:%s@tcp(%s:%d)/%s?%s"  //DSN格式:user:password@tcp(host:port)/dbName?settings
	settingFormat = "%s%s=%s&"	//设置参数格式:源串、请求参数名=请求值
)

// 数据源结构：驱动、主机地址、端口、数据库名、用户名、用户密码、数据源设置参数
type DataSource struct{
	driver string
	host string
	port int
	dbName string
	user string
	password string
	settings []Setting
}

//设置参数接口
type Setting func(string)string

//创建基础数据源，所需参数：主机地址、数据库名称、用户名称、用户密码
//默认驱动:mysql,默认端口:3306
func New(host,dbName,user,password string)*DataSource{
	return &DataSource{
		driver:defaultDriver,
		host:host,
		port:defaultPort,
		dbName:dbName,
		user:user,
		password:password,
	}
}

//设置自定义驱动
func (ds *DataSource) Driver(driver string) *DataSource{
	ds.driver = driver
	return ds
}

//设置自定义端口
func (ds *DataSource) Port(port int) *DataSource{
	ds.port = port
	return ds
}

//设置数据源请求参数列表
func (ds *DataSource) Set(sets... Setting) *DataSource{
	ds.settings = append(ds.settings, sets...)
	return ds
}

//创建数据库连接池
func (ds *DataSource) Open(ping bool) (*sql.DB,error){
	db, err := open(ds)
	if nil != err {
		return nil,err
	}
	if ping {
		err = db.Ping()
	}
	return db,err
}

//调用databse/sql创建连接池
func open(ds *DataSource) (*sql.DB,error){
	return sql.Open(ds.driver,realDSN(ds))
}

//构建完整数据库连接请求串（Data Source Name）
func realDSN(ds *DataSource) string{
	settings := concatSettings(ds.settings)
	dsn := fmt.Sprintf(defaultFormat,ds.user,ds.password,ds.host,ds.port,ds.dbName,settings)
	fmt.Println("dsn:",dsn)
	return strings.TrimRight(dsn,"?") //无设置参数时，去掉参数拼接符?
}

//拼接所有设置参数串,使用&串接
func concatSettings(settings []Setting)  string{
	sets := ""
	for _,fun := range settings {
		sets = fun(sets)
	}
	return strings.TrimRight(sets,"&")
}

//设置字符串类型参数：源串、请求参数名、请求参数值
func stringSetting(source,name,value string) string{
	if "" == value {
		return ""
	}
	return fmt.Sprintf(settingFormat,source,name,value)
}

//设置bool类型参数
func boolSetting(source,name string , value bool) string{
	return stringSetting(source,name,strconv.FormatBool(value))
}

//设置时间类型参数 超时时间不小于1ms,不超过24h
func timeSetting(source,name string, value time.Duration) string{
	if value < time.Millisecond || value > 24 * time.Hour {
		return ""
	}
	return fmt.Sprintf(settingFormat,source,name,value)
}

//设置字符编码
func SetCharset(value string) Setting{
	return func(source string) string{
		return stringSetting(source,"charset",value)
	}
}

//设置是否允许明文传输密码 allowCleartextPasswords
func SetAllowCleartextPasswords(value bool) Setting{
	return func(source string) string{
		return boolSetting(source,"allowCleartextPasswords",value)
	}
}

//设置是否格式化输出时间：0000-00-00 00:00:00 parseTime
func SetParseTime(value bool) Setting {
	return func(source string) string{
		return boolSetting(source,"parseTime",value)
	}
}

//设置请求超时时间 timeout
func SetTimeOut(value time.Duration) Setting{
	return func(source string) string{
		return timeSetting(source,"timeout",value)
	}
}

//设置读超时时间 readTimeout
func SetReadTimeout(value time.Duration) Setting{
	return func(source string) string{
		return timeSetting(source,"readTimeout",value)
	}
}

//设置写超时时间 writeTimeout
func SetWriteTimeout(value time.Duration) Setting{
	return func(source string) string{
		return timeSetting(source,"writeTimeout",value)
	}
}

