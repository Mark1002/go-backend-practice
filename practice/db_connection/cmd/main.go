package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mark1002/practice/db_connection"
)

func simulateClientAborts(pool *db_connection.DBPool, ch chan<- struct{}) {
	var wg sync.WaitGroup

	ctx := context.Background()
	fmt.Println("Multiple concurrent client aborts:")
	numClients := 5
	queryDuration := 15 * time.Minute
	abortAfter := 10 * time.Minute

	pool.SimulateMultipleClientAborts(ctx, numClients, queryDuration, abortAfter, &wg)
	wg.Wait()
	pool.PrintStats()
	ch <- struct{}{}
}

func main() {
	dsn := "appuser:apppassword@tcp(localhost:3306)/practice_db?parseTime=true"
	var ch chan struct{} = make(chan struct{})
	config := db_connection.PoolConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    1,
		ConnMaxLifetime: 1 * time.Second,
		ConnMaxIdleTime: 1 * time.Second,
	}

	pool, err := db_connection.NewDBPool(dsn, config)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer pool.Close()

	fmt.Println("=== Database Connection Pool with Client Abort Simulation ===")
	pool.PrintStats()
	go simulateClientAborts(pool, ch)
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ch:
			fmt.Println("Client abort simulation finished")
			return
		case <-ticker.C:
			pool.PrintStats()
		}
	}
}
