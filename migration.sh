goose -dir ./migrations postgres "postgresql://postgres:0000@localhost:5432/go-homework" status

goose -dir ./migrations postgres "postgresql://postgres:0000@localhost:5432/go-homework" up
