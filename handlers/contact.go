package handlers

import (
	"net/http"
)

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "contact.gohtml")
}
