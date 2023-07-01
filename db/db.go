package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	c  = config.Config.DBConfig
	db *gorm.DB
)

func getLogger() logger.Interface {
	recover()
	if c.GormLogMode == config.ENV_MODE_RELEASE {
		return nil
	}
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
		},
	)
}

func Connect() (err error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		c.Host,
		c.User,
		c.Password,
		c.DBName,
		c.Port,
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: getLogger(),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	return
}

func Get() *gorm.DB {
	return db
}
