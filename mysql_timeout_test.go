package main

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestMysqlDriverTimeout(t *testing.T) {
	const requestsCount = 5
	db, err := sql.Open("mysql", "root:example@/db_timeout")
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(requestsCount)

	db.SetMaxOpenConns(1)

	channelDoneAll := make(chan bool)

	for i := 0; i < requestsCount; i++ {
		go runMysql(db, &wg)
	}
	go wiatFor(&wg, channelDoneAll)

	select {
	case <-time.After(3 * time.Second):
		t.Fatalf("should timeout all operations")
	case <-channelDoneAll:
		println("OK Mysql")
	}
}

func runMysql(db *sql.DB, wg *sync.WaitGroup) {
	ctx1, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := db.ExecContext(ctx1, "SELECT sleep(10)")
	if err != nil {

	}
	defer wg.Done()
}
