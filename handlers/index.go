package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/rmacdiarmid/GPTSite/database"
	"github.com/rmacdiarmid/GPTSite/logger"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("IndexHandler called")

	// Retrieve data from the database
	articles, err := database.GetArticles()
	if err != nil {
		log.Printf("Error fetching articles: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Rendering index.gohtml with data: %#v", articles)

	// Load the Go template
	tpl, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the Go template with the data
	err = tpl.Execute(w, articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
