package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"todolist-api/config"

	"github.com/jmoiron/sqlx"
	// library pq
	_ "github.com/go-sql-driver/mysql"
)

// DB logical wrapper for database object
type DB struct {
	mtx    sync.RWMutex
	driver string
	dbs    []*sql.DB
	count  uint64
}

func scatter(n int, fn func(idx int) error) error {
	errCh := make(chan error, n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			errCh <- fn(idx)
		}(i)
	}

	var err error
	for i := 0; i < cap(errCh); i++ {
		if e := <-errCh; e != nil {
			err = e
		}
	}

	return err
}

// Open concurrently opens each underlying physical db
func Open(dbSetting *config.DBConfig) (*DB, error) {
	if dbSetting == nil {
		return nil, errors.New("database setting is required")
	}

	if dbSetting.Name == "" {
		return nil, errors.New("database driver name should not empty")
	}

	dsns := []string{dbSetting.Host}

	db := &DB{
		driver: dbSetting.Name,
		dbs:    make([]*sql.DB, len(dsns)),
	}

	err := scatter(len(db.dbs), func(idx int) error {
		dbConn, err := sql.Open(dbSetting.Name, dsns[idx])
		if err != nil {
			return err
		}

		dbConn.SetMaxOpenConns(dbSetting.MaxOpenConn)
		dbConn.SetMaxIdleConns(dbSetting.MaxIdleConn)
		dbConn.SetConnMaxLifetime(time.Duration(dbSetting.ConnMaxLifetime) * time.Second)
		db.dbs[idx] = dbConn

		return nil
	})

	if err != nil {
		return nil, err
	}

	fmt.Printf("success connect to database %s", dbSetting.Host)

	return db, nil
}

// Close closes all physical databases concurrently, releasing any open resources.
func (db *DB) Close() error {
	return scatter(len(db.dbs), func(idx int) error {
		return db.dbs[idx].Close()
	})
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool for each underlying physical db.
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns then the
// new MaxIdleConns will be reduced to match the MaxOpenConns limit
// If n <= 0, no idle connections are retained.
func (db *DB) SetMaxIdleConns(n int) {
	for idx := range db.dbs {
		db.dbs[idx].SetMaxIdleConns(n)
	}
}

// SetMaxOpenConns sets the maximum number of open connections
// to each physical database.
// If MaxIdleConns is greater than 0 and the new MaxOpenConns
// is less than MaxIdleConns, then MaxIdleConns will be reduced to match
// the new MaxOpenConns limit. If n <= 0, then there is no limit on the number
// of open connections. The default is 0 (unlimited).
func (db *DB) SetMaxOpenConns(n int) {
	for idx := range db.dbs {
		db.dbs[idx].SetMaxOpenConns(n)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
// Expired connections may be closed lazily before reuse.
// If d <= 0, connections are reused forever.
func (db *DB) SetConnMaxLifetime(d time.Duration) {
	for idx := range db.dbs {
		db.dbs[idx].SetConnMaxLifetime(d)
	}
}

// Master returns the master physical database
func (db *DB) Master() *sqlx.DB {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	return sqlx.NewDb(db.dbs[0], db.driver)
}

// spread all slave connection to all available slaves
func (db *DB) slave(n int) int {
	if n <= 1 {
		return 0
	}

	return int(1 + (atomic.AddUint64(&db.count, 1) % uint64(n-1)))
}

// Slave returns one of the physical databases which is a slave
func (db *DB) Slave() *sqlx.DB {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	return sqlx.NewDb(db.dbs[db.slave(len(db.dbs))], db.driver)
}

// Begin starts a transaction on the master. The isolation level is dependent on the driver.
func (db *DB) Begin(ctx context.Context) (*sqlx.Tx, error) {
	return db.Master().BeginTxx(ctx, nil)
}

// BeginTx starts a transaction with the provided context on the master.
//
// The provided TxOptions is optional and may be nil if defaults should be used.
// If a non-default isolation level is used that the driver doesn't support,
// an error will be returned.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.Master().BeginTx(ctx, opts)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
// Exec uses the master as the underlying physical db.
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Master().Exec(query, args...)
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
// Exec uses the master as the underlying physical db.
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.Master().ExecContext(ctx, query, args...)
}

// Ping verifies if a connection to each physical database is still alive,
// establishing a connection if necessary.
func (db *DB) Ping() error {
	return scatter(len(db.dbs), func(idx int) error {
		return db.dbs[idx].Ping()
	})
}

// PingContext verifies if a connection to each physical database is still
// alive, establishing a connection if necessary.
func (db *DB) PingContext(ctx context.Context) error {
	return scatter(len(db.dbs), func(idx int) error {
		return db.dbs[idx].PingContext(ctx)
	})
}

// Prepare creates a prepared statement for later queries or executions
// on each physical database, concurrently.
func (db *DB) Prepare(query string) (Stmt, error) {
	stmts := make([]*sql.Stmt, len(db.dbs))

	err := scatter(len(db.dbs), func(idx int) (err error) {
		stmts[idx], err = db.dbs[idx].Prepare(query)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &stmt{db: db, stmts: stmts}, nil
}

// PrepareContext creates a prepared statement for later queries or executions
// on each physical database, concurrently.
//
// The provided context is used for the preparation of the statement, not for
// the execution of the statement.
func (db *DB) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	stmts := make([]*sql.Stmt, len(db.dbs))

	err := scatter(len(db.dbs), func(idx int) (err error) {
		stmts[idx], err = db.dbs[idx].PrepareContext(ctx, query)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &stmt{db: db, stmts: stmts}, nil
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
// Query uses a slave as the physical db.
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Slave().Query(query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
// QueryContext uses a slave as the physical db.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.Slave().QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always return a non-nil value.
// Errors are deferred until Row's Scan method is called.
// QueryRow uses a slave as the physical db.
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.Slave().QueryRow(query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always return a non-nil value.
// Errors are deferred until Row's Scan method is called.
// QueryRowContext uses a slave as the physical db.
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.Slave().QueryRowContext(ctx, query, args...)
}
