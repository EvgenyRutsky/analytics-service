package storage

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"
	"time"
)

const (
	clickhouseHost = "kq74nrnq49.eu-central-1.aws.clickhouse.cloud:9440"
	database       = "default"
)

const (
	createTableSQL = `CREATE TABLE IF NOT EXISTS analytics (UUID UUID, Time DateTime, Value Int32) ORDER BY (Time)`
)

type ClickHouse struct {
	conn clickhouse.Conn
}

// NewClickHouse function creates ClickHouse object with opened connection and created table
func NewClickHouse(username string, password string) (*ClickHouse, error) {
	c := &ClickHouse{}

	if err := c.connect(username, password); err != nil {
		return c, err
	}

	if err := c.createTable(); err != nil {
		return c, err
	}

	return c, nil
}

func (c *ClickHouse) connect(username string, password string) error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{clickhouseHost},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		TLS:          &tls.Config{},
		Protocol:     clickhouse.HTTP,
		MaxOpenConns: 100, //~85 connections is enough for handling 50 RPS
		DialTimeout:  time.Second * 30,
	})
	if err != nil {
		return err
	}

	c.conn = conn

	v, err := c.conn.ServerVersion()
	fmt.Println(v)

	return err
}

func (c *ClickHouse) createTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := c.conn.Exec(ctx, createTableSQL); err != nil {
		return err
	}

	return nil
}

// SendAsync method performs async insert of the values to the ClickHouse DB
func (c *ClickHouse) SendAsync(t time.Time, value int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := c.conn.AsyncInsert(ctx, fmt.Sprintf(`INSERT INTO analytics VALUES ('%s', %d, %v)`, uuid.New(), t.UTC().Unix(), int32(value)), false); err != nil {
		return err
	}

	return nil
}

// GetValuesForLastMin method performs reading of values in the last min
func (c *ClickHouse) GetValuesForLastMin(t time.Time) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	queryTime := t.Add(-time.Minute)
	rows, err := c.conn.Query(ctx, fmt.Sprintf(`Select Value FROM analytics Where Time >= %v `, queryTime.Unix()))
	if err != nil {
		return nil, err
	}

	defer rows.Close() // Could be omitted as we perform defer cancel(), which releases all resources related to this context

	result := make([]int, 0)

	for rows.Next() {
		var tempInt32 int32
		if err := rows.Scan(&tempInt32); err != nil {
			return result, err
		}
		result = append(result, int(tempInt32))
	}

	return result, rows.Err()
}
