package config

/*
Main Database manager
*/

import (
	"database/sql"
	"fmt"
	"kanggo/pkg/entity/schema"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Native *sql.DB
	Gorm   *gorm.DB
)

func ConnectDb() {
	dbHost := EnvFile.DbHost
	dbPort := EnvFile.DbPort
	dbUser := EnvFile.DbUser
	dbPass := EnvFile.DbPassword
	dbName := EnvFile.DbName

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}

	Gorm = dbConn
	Native, err = dbConn.DB()
	if err != nil {
		log.Fatal(err)
	}

	Native.SetMaxOpenConns(25)
	Native.SetMaxIdleConns(25)
	Native.SetConnMaxLifetime(5 * time.Minute)

	autoCreate := EnvFile.DbAutoMigrate
	if autoCreate {
		fmt.Println("Dropping and recreating all tables...")

		// Auto migrate functionality
		Gorm.AutoMigrate(

			&schema.Order{},
			&schema.Product{},
			&schema.User{},
		)

		fmt.Println("All tables recreated successfully...")
	}

}
