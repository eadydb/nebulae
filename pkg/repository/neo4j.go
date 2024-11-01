package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Neo4jConn neo4j.DriverWithContext

// Neo4jDriver is a wrapper around the Neo4j driver
type Neo4jDriver struct {
	Host     string
	Port     int
	Username string
	Password string
	Driver   neo4j.DriverWithContext
	Ctx      context.Context
}

func NewNeo4jDriver(ctx context.Context, host, username, password string, port int) *Neo4jDriver {
	return &Neo4jDriver{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Ctx:      ctx,
	}
}

// Connect to Neo4j
func (n *Neo4jDriver) Connect() error {
	uri := fmt.Sprintf("neo4j://%s:%d", n.Host, n.Port)
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(n.Username, n.Password, ""))
	if err != nil {
		slog.Error("create Neo4j driver failed", slog.String("uri", uri), slog.String("username", n.Username), slog.String("err", err.Error()))
		return err
	}

	// Test the connection
	err = driver.VerifyConnectivity(n.Ctx)
	if err != nil {
		slog.Error("connect to Neo4j failed", slog.String("uri", uri), slog.String("username", n.Username), slog.String("err", err.Error()))
		return err
	}

	// Store the driver for later use
	n.Driver = driver
	Neo4jConn = driver
	return nil
}

func (n *Neo4jDriver) Close() error {
	return n.Driver.Close(n.Ctx)
}
