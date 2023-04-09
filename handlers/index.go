package handlers

import (
	"net/http"

	"github.com/rmacdiarmid/GPTSite/database"
	"github.com/rmacdiarmid/GPTSite/logger"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("IndexHandler called")

	// Retrieve data from the database
	articles, err := database.GetArticles()
	if err != nil {
		logger.DualLog.Printf("Error fetching articles: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.DualLog.Printf("Rendering index.gohtml with data: %#v", articles)

	// Render the Base template with the fetched data
	RenderTemplateWithData(w, "base.gohtml", articles)
}
