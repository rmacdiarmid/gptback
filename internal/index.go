package handlers

import (
	"net/http"

	"github.com/rmacdiarmid/GPTSite/logger"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
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

	// Prepare the data structure for the base template
	data := struct {
		Content  string
		Articles []database.Article
	}{
		Content:  "index.gohtml",
		Articles: articles,
	}

	// Render the Base template with the fetched data
	RenderTemplateWithData(w, "base.gohtml", data)
}
