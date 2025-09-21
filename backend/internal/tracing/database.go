package tracing

import (
	"context"
	"database/sql/driver"
	"time"

	"database/sql"

	"go.opentelemetry.io/otel/attribute"
)

// DBTracer is a wrapper around a driver.Driver that traces database operations
type DBTracer struct {
	driver driver.Driver
}

// NewDBTracer creates a new DBTracer
func NewDBTracer(d driver.Driver) *DBTracer {
	return &DBTracer{driver: d}
}

// Open implements the driver.Driver interface
func (t *DBTracer) Open(name string) (driver.Conn, error) {
	conn, err := t.driver.Open(name)
	if err != nil {
		return nil, err
	}
	return &tracedConn{conn: conn}, nil
}

// tracedConn is a wrapper around a driver.Conn that traces database operations
type tracedConn struct {
	conn driver.Conn
}

// Prepare implements the driver.Conn interface
func (c *tracedConn) Prepare(query string) (driver.Stmt, error) {
	ctx := context.Background()
	_, span := StartSpan(ctx, "db.prepare")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.query", query),
	)

	stmt, err := c.conn.Prepare(query)
	if err != nil {
		SetSpanError(ctx, err)
		return nil, err
	}

	return &tracedStmt{stmt: stmt, query: query}, nil
}

// Close implements the driver.Conn interface
func (c *tracedConn) Close() error {
	ctx := context.Background()
	_, span := StartSpan(ctx, "db.close")
	defer span.End()

	err := c.conn.Close()
	if err != nil {
		SetSpanError(ctx, err)
		return err
	}

	return nil
}

// Begin implements the driver.Conn interface
func (c *tracedConn) Begin() (driver.Tx, error) {
	ctx := context.Background()
	_, span := StartSpan(ctx, "db.begin")
	defer span.End()

	tx, err := c.conn.Begin()
	if err != nil {
		SetSpanError(ctx, err)
		return nil, err
	}

	return &tracedTx{tx: tx}, nil
}

// tracedStmt is a wrapper around a driver.Stmt that traces database operations
type tracedStmt struct {
	stmt  driver.Stmt
	query string
}

// Close implements the driver.Stmt interface
func (s *tracedStmt) Close() error {
	ctx := context.Background()
	_, span := StartSpan(ctx, "db.stmt.close")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.query", s.query),
	)

	err := s.stmt.Close()
	if err != nil {
		SetSpanError(ctx, err)
		return err
	}

	return nil
}

// NumInput implements the driver.Stmt interface
func (s *tracedStmt) NumInput() int {
	return s.stmt.NumInput()
}

// Exec implements the driver.Stmt interface
func (s *tracedStmt) Exec(args []driver.Value) (driver.Result, error) {
	ctx := context.Background()
	_, span := StartSpan(ctx, "db.stmt.exec")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.query", s.query),
		attribute.Int("db.args_count", len(args)),
	)

	start := time.Now()
	result, err := s.stmt.Exec(args)
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Float64("db.duration_ms", float64(duration.Milliseconds())),
	)

	if err != nil {
		SetSpanError(ctx, err)
		return nil, err
	}

	return result, nil
}

// Query implements the driver.Stmt interface
func (s *tracedStmt) Query(args []driver.Value) (driver.Rows, error) {
	ctx := context.Background()
	_, span := StartSpan(ctx, "db.stmt.query")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.query", s.query),
		attribute.Int("db.args_count", len(args)),
	)

	start := time.Now()
	rows, err := s.stmt.Query(args)
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Float64("db.duration_ms", float64(duration.Milliseconds())),
	)

	if err != nil {
		SetSpanError(ctx, err)
		return nil, err
	}

	return rows, nil
}

// tracedTx is a wrapper around a driver.Tx that traces database operations
type tracedTx struct {
	tx driver.Tx
}

// Commit implements the driver.Tx interface
func (t *tracedTx) Commit() error {
	ctx := context.Background()
	_, span := StartSpan(ctx, "db.tx.commit")
	defer span.End()

	err := t.tx.Commit()
	if err != nil {
		SetSpanError(ctx, err)
		return err
	}

	return nil
}

// Rollback implements the driver.Tx interface
func (t *tracedTx) Rollback() error {
	ctx := context.Background()
	_, span := StartSpan(ctx, "db.tx.rollback")
	defer span.End()

	err := t.tx.Rollback()
	if err != nil {
		SetSpanError(ctx, err)
		return err
	}

	return nil
}

// TraceDBQuery traces a database query
func TraceDBQuery(ctx context.Context, operation, query string, args []interface{}, fn func() error) error {
	_, span := StartSpan(ctx, "db.query")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.operation", operation),
		attribute.String("db.query", query),
		attribute.Int("db.args_count", len(args)),
	)

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Float64("db.duration_ms", float64(duration.Milliseconds())),
	)

	if err != nil {
		SetSpanError(ctx, err)
		return err
	}

	return nil
}

// TraceDBExec traces a database exec operation
func TraceDBExec(ctx context.Context, operation, query string, args []interface{}, fn func() (driver.Result, error)) (driver.Result, error) {
	_, span := StartSpan(ctx, "db.exec")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.operation", operation),
		attribute.String("db.query", query),
		attribute.Int("db.args_count", len(args)),
	)

	start := time.Now()
	result, err := fn()
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Float64("db.duration_ms", float64(duration.Milliseconds())),
	)

	if err != nil {
		SetSpanError(ctx, err)
		return nil, err
	}

	return result, nil
}

// TraceDBQueryRow traces a database query row operation
func TraceDBQueryRow(ctx context.Context, operation, query string, args []interface{}, fn func() *sql.Row) *sql.Row {
	_, span := StartSpan(ctx, "db.query_row")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.operation", operation),
		attribute.String("db.query", query),
		attribute.Int("db.args_count", len(args)),
	)

	start := time.Now()
	row := fn()
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Float64("db.duration_ms", float64(duration.Milliseconds())),
	)

	return row
}

// TraceDBBegin traces a database begin transaction operation
func TraceDBBegin(ctx context.Context, fn func() (*sql.Tx, error)) (*sql.Tx, error) {
	_, span := StartSpan(ctx, "db.begin")
	defer span.End()

	start := time.Now()
	tx, err := fn()
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Float64("db.duration_ms", float64(duration.Milliseconds())),
	)

	if err != nil {
		SetSpanError(ctx, err)
		return nil, err
	}

	return tx, nil
}

// TraceDBCommit traces a database commit transaction operation
func TraceDBCommit(ctx context.Context, fn func() error) error {
	_, span := StartSpan(ctx, "db.commit")
	defer span.End()

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Float64("db.duration_ms", float64(duration.Milliseconds())),
	)

	if err != nil {
		SetSpanError(ctx, err)
		return err
	}

	return nil
}

// TraceDBRollback traces a database rollback transaction operation
func TraceDBRollback(ctx context.Context, fn func() error) error {
	_, span := StartSpan(ctx, "db.rollback")
	defer span.End()

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Float64("db.duration_ms", float64(duration.Milliseconds())),
	)

	if err != nil {
		SetSpanError(ctx, err)
		return err
	}

	return nil
}
