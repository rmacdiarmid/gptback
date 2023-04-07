package handlers

import (
	"net/http"
)

func ActivityHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "activity.html")
}
