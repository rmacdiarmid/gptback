package handlers

import (
	"net/http"

	"github.com/rmacdiarmid/GPTSite/logger"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Log message for 404 error
	logger.DualLog.Printf("404 error: %s", r.URL.Path)

	w.WriteHeader(http.StatusNotFound)

	data := struct {
		Content string
	}{
		Content: "404.gohtml",
	}

	// Call the function without using the result as a value
	RenderTemplateWithData(w, "base.gohtml", data)
}
