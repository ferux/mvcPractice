//go:generate stringer -type DatabaseType

//Package mvcPractice.
//When building application outside of docker container you should specify flag go build -ldflags "-X main.outsideDocker=true"
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ferux/mvcPractice/controller"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

//DatabaseType implementation
type DatabaseType int

//DB Types for DatabaseType
const (
	POSTGRESQL DatabaseType = iota + 1
	MONGO
	MYSQL
	MSSQL
	SQLDB
)

type config struct {
	ip, port, dbAddr, dbPort, user, pwd, dbName, driver, extIP string
}

func getParamsDocker() config {
	log.Println("Using docker params")
	var ip, port, dbAddr, dbPort, user, pwd, dbName, driver, extIP string
	flag.StringVar(&ip, "ip", "127.0.0.1", "Hosting IP")
	flag.StringVar(&port, "port", "8080", "Hosting Port")
	flag.StringVar(&dbAddr, "dbAddr", "127.0.0.1", "Database IP")
	flag.StringVar(&dbPort, "dbPort", "5432", "Database Port")
	flag.StringVar(&user, "user", "user", "Database User")
	flag.StringVar(&pwd, "password", "pass", "Database Password")
	flag.StringVar(&dbName, "db", "default", "Database Name")
	flag.StringVar(&driver, "driver", "postgres", "Database driver")
	flag.StringVar(&extIP, "extIp", "127.0.0.1", "External IP for CORS")
	return config{
		ip:     ip,
		port:   port,
		dbAddr: dbAddr,
		dbPort: dbPort,
		user:   user,
		pwd:    pwd,
		dbName: dbName,
		driver: driver,
		extIP:  extIP,
	}
}

func getParamsAlone() config {
	log.Println("Using console params")
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	dbAddr := os.Getenv("DB_IP")
	dbPort := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PWD")
	dbName := os.Getenv("DB")
	extIP := os.Getenv("EXT_IP")
	driver := "postgres"
	return config{
		ip:     ip,
		port:   port,
		dbAddr: dbAddr,
		dbPort: dbPort,
		user:   user,
		pwd:    pwd,
		dbName: dbName,
		driver: driver,
		extIP:  extIP,
	}
}

var outsideDocker string

//In case of runnig application outside of docker change function from getParamsDocker to getParamsAlone
func main() {
	log.Printf("ousideDocker value: %v", outsideDocker)
	var mainConfig config
	if outsideDocker == "false" {
		mainConfig = getParamsDocker()
	} else {
		mainConfig = getParamsAlone()
	}
	conf := controller.Config{
		ListenIP:   mainConfig.ip,
		ListenPort: mainConfig.port,
		ExternalIP: mainConfig.extIP,
	}

	ds := fmt.Sprintf(`%s://%s:%s@%s:%s/%s?sslmode=disable`,
		mainConfig.driver,
		mainConfig.user,
		mainConfig.pwd,
		mainConfig.dbAddr,
		mainConfig.dbPort,
		mainConfig.dbName,
	)

	db, err := func() (*sqlx.DB, error) {
		var err error
		var db *sqlx.DB
		for i := 1; i < 11; i++ {
			db, err = sqlx.Connect(mainConfig.driver, ds)
			if err == nil {
				return db, nil
			}
			log.Println("Got an error while connecting to DB: ", err)
			log.Printf("Attempt #%d. Reconnecting in 3 seconds.", i)
			time.Sleep(time.Second * 3)
		}
		return nil, err
	}()

	if err != nil {
		log.Fatalf("Can't connect to DB. Reason: %v", err)
	}
	controller.Run(conf, db)
}
