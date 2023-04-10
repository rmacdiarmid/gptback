package handlers

import (
	"net/http"

	"github.com/rmacdiarmid/GPTSite/logger"
)

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("ContactHandler called")

	// Prepare the data structure for the base template
	data := struct {
		Content string
	}{
		Content: "contact.gohtml",
	}

	// Render the Base template with the content template name
	RenderTemplateWithData(w, "base.gohtml", data)
}
