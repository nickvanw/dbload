package main

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/*
CREATE TABLE `data` (
  `id` int NOT NULL AUTO_INCREMENT,
  `host` varchar(100) DEFAULT NULL,
  `data` varchar(100) DEFAULT NULL,
  `now` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
*/
var query string = "INSERT INTO data (host, data) VALUES (?, ?)"

func main() {
	if err := realMain(5*time.Second, 25); err != nil {
		fmt.Printf("error running realMain: %s\n", err)
		os.Exit(1)
	}
}

func realMain(delay time.Duration, parallelism int) error {
	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return errors.New("no DATABASE_URL")
	}

	h, err := os.Hostname()
	if err != nil {
		return err
	}

	for i := 0; i < parallelism; i++ {
		fmt.Printf("launching insert loop %d\n", i)
		go func() {
			if err := insertLoop(delay, dsn, h); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}()
	}

	done := make(chan struct{})

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		close(done)
	}()

	<-done
	fmt.Println("i'm done")

	return nil
}

func insertLoop(delay time.Duration, dsn, hostname string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	q, err := db.Prepare(query)
	if err != nil {
		return err
	}
	for {
		data := randData(20)
		fmt.Printf("inputting id=%s, data=%s\n", hostname, data)
		if _, err := q.Exec(hostname, data); err != nil {
			fmt.Printf("error inserting data: %s", err)
			return err
		}
		time.Sleep(delay)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randData(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
