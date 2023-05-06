package internal

import (
	"net/http"

	"github.com/rmacdiarmid/gptback/logger"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Log message for 404 error
	logger.DualLog.Printf("404 error: %s", r.URL.Path)

	w.WriteHeader(http.StatusNotFound)

	data := map[string]interface{}{
		"ContentTemplateName": "404",
	}

	// Call the function without using the result as a value
	RenderTemplateWithData(w, "base.gohtml", "404.gohtml", data)
}
