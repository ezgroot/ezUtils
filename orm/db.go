package orm

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB(conf Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch conf.SQLName {
	case "postgres":
		{

			dsn := url.URL{
				User:     url.UserPassword(conf.Account, conf.Password),
				Scheme:   conf.SQLName,
				Host:     fmt.Sprintf("%s:%d", conf.Address, conf.Port),
				Path:     conf.Database,
				RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
			}
			db, err = gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
			if err != nil {
				fmt.Printf("create postgres db handle error = %s\n", err)
				return nil, err
			}
		}
	case "mysql":
		{
			dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				conf.Account, conf.Password, conf.Address, conf.Port, conf.Database)
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				fmt.Printf("create mysql db handle error = %s\n", err)
				return nil, err
			}
		}
	default:
		{
			fmt.Printf("%s db type not support\n", conf.SQLName)
			return nil, fmt.Errorf("%s db type not support", conf.SQLName)
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("get std db error = %s\n", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdleConn)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLifetime))

	return db, nil
}

func ShowDbStatus(gdb *gorm.DB) {
	if gdb == nil {
		return
	}

	db, err := gdb.DB()
	if err != nil {
		fmt.Printf("get std db error = %s\n", err)
		return
	}

	state := db.Stats()
	fmt.Printf("---------------------------- Db-info ----------------------------\n"+
		"* Maximum number of open connections to the database : %d\n"+
		"             **************** pool state ****************\n"+
		"* The number of established connections both in use and idle : %d\n"+
		"* The number of connections currently in use : %d\n"+
		"* The number of idle connections : %d\n"+
		"           **************** statistics info ****************\n"+
		"* The total number of connections waited for : %d\n"+
		"* The total time blocked waiting for a new connection : %d\n"+
		"* The total number of connections closed due to SetMaxIdleConns : %d\n"+
		"* The total number of connections closed due to SetConnMaxLifetime : %d\n"+
		"-----------------------------------------------------------------\n",
		state.MaxOpenConnections, state.OpenConnections, state.InUse, state.Idle,
		state.WaitCount, state.WaitDuration, state.MaxIdleClosed, state.MaxLifetimeClosed)

	return
}
