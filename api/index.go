package handler

import (
	"net/http"

	"money-tracker-api/app"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	app.Handler(w, r)
}
