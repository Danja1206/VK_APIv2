package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type contextKey int


const dbKey contextKey = iota


func InitDB(dbUser,dbPassword,dbIP,dbPort,dbName string ) (*sql.DB, error) {
	
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbIP, dbPort, dbName)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func WithDB(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func FromContext(ctx context.Context) (*sql.DB, bool) {
	db, ok := ctx.Value(dbKey).(*sql.DB)
	return db, ok
}