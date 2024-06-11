package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"semay.com/configs"
)

var (
	DBConn *gorm.DB
)

func GormLoggerFile() *os.File {

	gormLogFile, gerr := os.OpenFile("gormblue.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if gerr != nil {
		log.Fatalf("error opening file: %v", gerr)
	}
	return gormLogFile
}

func ReturnSession() *gorm.DB {

	//  setting up database connection based on DB type

	app_env := configs.AppConfig.Get("DB_TYPE")
	//  This is file to output gorm logger on to
	gormlogger := GormLoggerFile()
	gormFileLogger := log.Logger{}
	gormFileLogger.SetOutput(gormlogger)
	gormFileLogger.Writer()

	gormLogger := log.New(gormFileLogger.Writer(), "\r\n", log.LstdFlags|log.Ldate|log.Ltime|log.Lshortfile)
	newLogger := logger.New(
		gormLogger, // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			Colorful:                  true,        // Enable color
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			// ParameterizedQueries:      true,        // Don't include params in the SQL log

		},
	)

	var DBSession *gorm.DB

	switch app_env {
	case "postgres":
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  configs.AppConfig.Get("POSTGRES_URI"),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage,

		}), &gorm.Config{
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
		})
		if err != nil {
			panic(err)
		}

		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)

		DBSession = db
	case "":
		//  this is sqlite connection
		db, _ := gorm.Open(sqlite.Open("goframedb"), &gorm.Config{
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
		})
		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)

		DBSession = db
	default:
		//  this is sqlite connection
		db, _ := gorm.Open(sqlite.Open(configs.AppConfig.Get("SQLITE_URI")), &gorm.Config{
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
		})

		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)

		DBSession = db

	}
	return DBSession

}
