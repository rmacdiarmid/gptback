package handlers

import (
	"net/http"

	"github.com/rmacdiarmid/GPTSite/database"
	"github.com/rmacdiarmid/GPTSite/logger"
)

func ArticleGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	templates.ExecuteTemplate(w, "article_generator.gohtml", nil)
}

func GenerateArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		prompt := r.FormValue("prompt")
		title, imageURL, articleText, err := GenerateArticle(prompt)
		if err != nil {
			// Log the error
			logger.DualLog.Printf("Error generating article: %v", err)
			http.Error(w, "Error generating article", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Generated":   true,
			"Title":       title,
			"ImageURL":    imageURL,
			"ArticleText": articleText,
		}
		err = templates.ExecuteTemplate(w, "article_generator.gohtml", data)
		if err != nil {
			// Log the error
			logger.DualLog.Printf("Error executing template: %v", err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	} else {
		err := templates.ExecuteTemplate(w, "article_generator.gohtml", nil)
		if err != nil {
			// Log the error
			logger.DualLog.Printf("Error executing template: %v", err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	}
}

func AcceptArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		imageURL := r.FormValue("image_url")
		articleText := r.FormValue("article_text")

		_, err := database.InsertArticle(title, imageURL, articleText)
		if err != nil {
			// Handle error
			http.Error(w, "Error uploading article", http.StatusInternalServerError)
			return
		}

		// Redirect to a success or confirmation page, or any other desired page
		http.Redirect(w, r, "/success", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
