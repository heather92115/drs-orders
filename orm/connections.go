package orm

import (
	"drs-orders/orders"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var globalDb *gorm.DB

func CreatePool() (err error) {
	// https://github.com/go-gorm/postgres
	globalDb, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=heather.stevens dbname=orders port=5432 sslmode=disable",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	sqlDB, err := globalDb.DB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Printf("Created %d db connections\n", sqlDB.Stats().OpenConnections)
	return nil
}

func GetConnection() (db *gorm.DB) {
	db = globalDb
	return
}

func MigrateTables() (err error) {

	err = globalDb.AutoMigrate(orders.Order{})
	if err != nil {
		return err
	}

	err = globalDb.AutoMigrate(orders.OrderItem{})
	if err != nil {
		return err
	}

	err = globalDb.AutoMigrate(orders.OrderBatch{})
	if err != nil {
		return err
	}

	err = globalDb.AutoMigrate(orders.Status{})
	if err != nil {
		return err
	}

	err = globalDb.AutoMigrate(orders.AuditTrail{})
	if err != nil {
		return err
	}

	return
}
