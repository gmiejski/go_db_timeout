package main

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const timeout = 2 * time.Second

func TestDriverTimeout(t *testing.T) {
	const requestsCount = 5

	connStr := "user=postgres password=postgres dbname=db_timeout sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(requestsCount)

	db.SetMaxOpenConns(1)

	channelDoneAll := make(chan bool)

	for i := 0; i < requestsCount; i++ {
		go run(db, &wg)
	}
	go wiatFor(&wg, channelDoneAll)

	select {
	case <-time.After(3 * time.Second):
		t.Fatalf("should timeout all operations")
	case <-channelDoneAll:
		println("OK postgres")
	}
}
func wiatFor(group *sync.WaitGroup, errors chan bool) {
	group.Wait()
	errors <- true
}

func run(db *sql.DB, wg *sync.WaitGroup) {
	ctx1, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := db.ExecContext(ctx1, "SELECT pg_sleep(10)")
	if err != nil {

	}
	defer wg.Done()
}
