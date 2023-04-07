package handlers

import (
	"log"
	"net/http"

	"github.com/rmacdiarmid/GPTSite/database"
)

type Article struct {
	ID      int
	Title   string
	Image   string
	Preview string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("IndexHandler called")

	// Retrieve data from the database
	articles, err := database.GetArticles(DB)
	if err != nil {
		log.Printf("Error fetching articles: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Rendering index.html with data: %#v", articles)

	// Render the HTML template with the data
	renderTemplateWithData(w, "index.html", articles)
}
