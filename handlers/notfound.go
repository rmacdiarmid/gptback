package handlers

import (
	"net/http"

	"github.com/rmacdiarmid/GPTSite/logger"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Log message for 404 error
	logger.DualLog.Printf("404 error: %s", r.URL.Path)

	w.WriteHeader(http.StatusNotFound)
	if err := RenderTemplate(w, "404.gohtml"); err != nil {
		logger.DualLog.Printf("Error rendering 404 template: %v", err)
	}
}
