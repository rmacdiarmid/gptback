package handlers

import (
	"log"
	"net/http"

	"github.com/rmacdiarmid/GPTSite/database"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("IndexHandler called")

	// Retrieve data from the database
	articles, err := database.GetArticles()
	if err != nil {
		log.Printf("Error fetching articles: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Rendering index.html with data: %#v", articles)

	// Render the HTML template with the data
	renderTemplateWithData(w, "index.html", articles)
}
