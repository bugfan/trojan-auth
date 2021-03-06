package models

import (
	"fmt"
	"os"

	"github.com/bugfan/trojan-auth/env"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var (
	x      *xorm.Engine
	tables []interface{}
)

func init() {
	Register(&Credential{})

	_, err := SetEngine(&Config{
		User:     env.Get("db_user"),
		Password: env.Get("db_pwd"),
		Host:     env.Get("db_host"),
		Name:     env.Get("db_name"),
		ShowSQL:  env.GetBool("db_show_sql"),
	}, env.Get("db_scheme"))
	if err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
}
func Register(obj ...interface{}) {
	tables = append(tables, obj...)

}

func GetEngine() *xorm.Engine {
	return x
}

type Config struct {
	User     string
	Password string
	Host     string
	Name     string
	Log      string
	ShowSQL  bool
}

func newEngine(config *Config) (*xorm.Engine, error) {
	var cnnstr string
	if config.Host[0] == '/' { // looks like a unix socket
		cnnstr = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8&parseTime=true",
			config.User, config.Password, config.Host, config.Name)
	} else {
		cnnstr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true",
			config.User, config.Password, config.Host, config.Name)
	}
	x, err := xorm.NewEngine("mysql", cnnstr)
	if err != nil {
		return nil, err
	}
	x.SetMapper(core.GonicMapper{})

	if config.Log != "" {
		f, err := os.Create(config.Log)
		if err != nil {
			return nil, fmt.Errorf("Fail to create xorm.log: %v", err)
		}
		x.SetLogger(xorm.NewSimpleLogger(f))
	}
	x.ShowSQL(config.ShowSQL)
	return x, nil
}
func newSqlite3Engine(config *Config) (*xorm.Engine, error) {
	return xorm.NewEngine("sqlite3", "./auth.db")
}
func SetEngine(config *Config, args ...string) (*xorm.Engine, error) {
	var err error
	if len(args) > 0 && "sqlite3" == args[0] {
		if x, err = newSqlite3Engine(config); err != nil {
			return nil, err
		}
	} else {
		if x, err = newEngine(config); err != nil {
			return nil, err
		}
	}
	if err = x.StoreEngine("InnoDB").Sync2(tables...); err != nil {
		fmt.Printf("sync database struct error: %v\n", err)
		return nil, err
	}
	initData()
	return x, nil
}
func initData() {
	// InitData()
}
func All(obj interface{}) error {
	return x.Find(obj)
}

func Get(id int64, obj interface{}) (has bool, err error) {
	return x.Id(id).Get(obj)
}

func Insert(beans ...interface{}) (int64, error) {
	return x.Insert(beans...)
}

func Delete(id int64, obj interface{}) error {
	_, err := x.ID(id).Delete(obj)
	return err
}

func Update(id int64, obj interface{}) error {
	_, err := x.ID(id).Update(obj)
	return err
}

func UpdateAll(id int64, obj interface{}) error {
	_, err := x.ID(id).AllCols().Update(obj)
	return err
}
