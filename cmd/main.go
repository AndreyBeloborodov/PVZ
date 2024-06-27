package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/console"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/services"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/storage"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/storage/transactor"
	"log"
	"os"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "0000"
	dbname   = "go-homework"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	txManager := transactor.NewTransactionManager(pool)

	cmd := console.NewConsole(services.NewOrderService(storage.NewOrderStorage(txManager), txManager))

	for {
		var command string
		command, _ = in.ReadString('\n')
		go cmd.Solve(ctx, command)
	}
}
