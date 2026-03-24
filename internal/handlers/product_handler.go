package handlers

import "github.com/jmoiron/sqlx"

type ProductHandler struct {
	// We will add the DB connection here later when we implement the methods
	DB *sqlx.DB
}
