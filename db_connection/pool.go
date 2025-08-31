package db_connection

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBPool struct {
	DB *sql.DB
}

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func NewDBPool(dsn string, config PoolConfig) (*DBPool, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool settings
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connection pool configured:")
	log.Printf("  Max Open Connections: %d", config.MaxOpenConns)
	log.Printf("  Max Idle Connections: %d", config.MaxIdleConns)
	log.Printf("  Connection Max Lifetime: %v", config.ConnMaxLifetime)
	log.Printf("  Connection Max Idle Time: %v", config.ConnMaxIdleTime)

	return &DBPool{DB: db}, nil
}

func (p *DBPool) Close() error {
	return p.DB.Close()
}

func (p *DBPool) GetStats() sql.DBStats {
	return p.DB.Stats()
}

func (p *DBPool) PrintStats() {
	stats := p.GetStats()
	log.Printf("Connection Pool Stats:")
	log.Printf("  Open Connections: %d", stats.OpenConnections)
	log.Printf("  In Use: %d", stats.InUse)
	log.Printf("  Idle: %d", stats.Idle)
	log.Printf("  Wait Count: %d", stats.WaitCount)
	log.Printf("  Wait Duration: %v", stats.WaitDuration)
	log.Printf("  Max Idle Closed: %d", stats.MaxIdleClosed)
	log.Printf("  Max Idle Time Closed: %d", stats.MaxIdleTimeClosed)
	log.Printf("  Max Lifetime Closed: %d", stats.MaxLifetimeClosed)
}

func (p *DBPool) SimulateClientAbort(ctx context.Context, queryDuration time.Duration, abortAfter time.Duration) error {
	abortCtx, cancel := context.WithTimeout(ctx, abortAfter)
	defer cancel()

	query := "SELECT SLEEP(?)"

	log.Printf("Starting query that will run for %v, but client will abort after %v", queryDuration, abortAfter)

	_, err := p.DB.QueryContext(abortCtx, query, queryDuration.Seconds())
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Printf("Client aborted connection after %v (simulating client abort)", abortAfter)
			return fmt.Errorf("client abort simulated: %w", err)
		}
		return fmt.Errorf("query failed: %w", err)
	}

	log.Printf("Query completed successfully (no abort occurred)")
	return nil
}

func (p *DBPool) SimulateMultipleClientAborts(
	ctx context.Context,
	numClients int,
	queryDuration time.Duration,
	abortAfter time.Duration,
	wg *sync.WaitGroup,
) {
	log.Printf("Simulating %d client aborts", numClients)
	wg.Add(numClients)
	for i := 0; i < numClients; i++ {
		go func(clientID int) {
			defer wg.Done()
			log.Printf("Client %d: Starting simulation", clientID)
			err := p.SimulateClientAbort(ctx, queryDuration, abortAfter)
			if err != nil {
				log.Printf("Client %d: %v", clientID, err)
			} else {
				log.Printf("Client %d: Completed successfully", clientID)
			}
		}(i + 1)
	}
}
