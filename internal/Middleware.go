package internal

import (
	"database/sql"
	"net/http"
	"github.com/Danaj1412/vk/database"

	_ "github.com/go-sql-driver/mysql"

)

func WithDBContextMiddleware(next http.HandlerFunc, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = database.WithDB(ctx, db)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}