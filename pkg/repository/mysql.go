package repository

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GormConn *gorm.DB

type MysqlDriver struct {
	Unix     string
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// NewMysqlDriver creates a new mysql driver
func NewMysqlDriver(unix, host, username, password, database string, port int) *MysqlDriver {
	return &MysqlDriver{
		Unix:     unix,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}
}

// MysqlParseUrl returns a mysql connection string
func (d *MysqlDriver) MysqlParseUrl() string {
	var connHost string
	if d.Unix != "" {
		connHost = fmt.Sprintf("unix(%s)", d.Unix)
	} else {
		connHost = fmt.Sprintf("tcp(%s:%d)", d.Host, d.Port)
	}
	return fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local", d.Username, d.Password, connHost, d.Database)
}

// GetGromConn returns a gorm connection
func (d *MysqlDriver) GetGromConn(ctx context.Context) (*gorm.DB, error) {
	dns := d.MysqlParseUrl()
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		slog.Error("open mysql connection failed", slog.Any("error", err))
		return nil, err
	}
	return db, nil
}
