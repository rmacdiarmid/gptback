package handlers

import (
	"net/http"

	"github.com/rmacdiarmid/GPTSite/logger"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("IndexHandler called")
	articles, err := database.GetArticles()
	if err != nil {
		http.Error(w, "Error fetching articles", http.StatusInternalServerError)
		return
	}
	logger.DualLog.Printf("Fetched articles: %v", articles)

	data := map[string]interface{}{
		"ContentTemplateName": "index.gohtml",
		"Articles":            articles,
	}

	RenderTemplateWithData(w, r, "base.gohtml", data)
}
