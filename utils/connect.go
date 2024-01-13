package utils

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/jackc/pgconn"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var MockPoolsDB = poolsDB{}

type PGXPool interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
	Close()
}

type poolsDB struct {
	pgConn    map[string]*pgconn.Notification
	connMutex sync.Mutex
	enabled   bool
}

func ConnectDBPool(connString string) *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		LogAndPanicIfError(err, "unable to create connection pool")
	}

	return dbpool
}

func ConnectDB(connString string) *sql.DB {
	dbc, err := sql.Open("pgx", connString)
	LogAndPanicIfError(err, "failed when connecting to database")
	// Maximum Idle Connections
	dbc.SetMaxIdleConns(10)
	// Maximum Open Connections
	dbc.SetMaxOpenConns(20)
	// Idle Connection Timeout
	dbc.SetConnMaxIdleTime(15 * time.Second)
	// Connection Lifetime
	dbc.SetConnMaxLifetime(60 * time.Second)

	err = dbc.Ping()
	LogAndPanicIfError(err, "failed when ping to database")

	return dbc
}

func OpenPool(connString string) (*pgxpool.Pool, error) {
	if MockPoolsDB.enabled {
		return &pgxpool.Pool{}, nil
	}

	return pgxpool.Connect(context.Background(), connString)
}

// func (m *poolsDB) Start() {
// 	m.enabled = true
// }

// func (m *poolsDB) Stop() {
// 	m.enabled = false
// }

// func (m *poolsDB) GetMockPool(channel string) *pgconn.Notification {
// 	return m.pgConn[channel]
// }

// func (m *poolsDB) DeleteMockPool() {
// 	m.connMutex.Lock()
// 	defer m.connMutex.Unlock()

// 	m.pgConn = make(map[string]*pgconn.Notification)
// }

// func (m *poolsDB) Add(pgConn map[string]*pgconn.Notification) {
// 	m.connMutex.Lock()
// 	defer m.connMutex.Unlock()

// 	m.pgConn = pgConn
// }
