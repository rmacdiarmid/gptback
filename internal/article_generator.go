package internal

import (
	"net/http"

	"github.com/rmacdiarmid/GPTSite/logger"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
)

func ArticleGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting the ArticleGeneratorHandler function...")
	defer logger.DualLog.Println("Exiting the ArticleGeneratorHandler function.")

	if r.Method != "GET" {
		logger.DualLog.Printf("Invalid request method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := map[string]interface{}{
		"ContentTemplateName": "article_generator",
	}

	RenderTemplateWithData(w, "base.gohtml", data)
}

func GenerateArticleHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting the GenerateArticleHandler function...")
	defer logger.DualLog.Println("Exiting the GenerateArticleHandler function.")

	if r.Method != "POST" {
		logger.DualLog.Printf("Invalid request method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		logger.DualLog.Printf("Invalid content type: %s", r.Header.Get("Content-Type"))
		http.Error(w, "Invalid content type", http.StatusUnsupportedMediaType)
		return
	}
	prompt := r.FormValue("prompt")
	title, imageURL, articleText, err := GenerateArticle(prompt)
	if err != nil {
		// Log the error
		logger.DualLog.Printf("Error generating article: %v", err)
		http.Error(w, "Error generating article", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Content":     "article_generator.gohtml",
		"Generated":   true,
		"Title":       title,
		"ImageURL":    imageURL,
		"ArticleText": articleText,
	}

	RenderTemplateWithData(w, "base.gohtml", data)
}

func AcceptArticleHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting the AcceptArticleHandler function...")
	defer logger.DualLog.Println("Exiting the AcceptArticleHandler function.")

	if r.Method != "POST" {
		logger.DualLog.Printf("Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	imageURL := r.FormValue("image_url")
	articleText := r.FormValue("article_text")

	_, err := database.InsertArticle(title, imageURL, articleText)
	if err != nil {
		// Handle error
		logger.DualLog.Printf("Error uploading article: %v", err)
		http.Error(w, "Error uploading article", http.StatusInternalServerError)
		return
	}

	// Redirect to a success or confirmation page, or any other desired page
	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting the SuccessHandler function...")
	defer logger.DualLog.Println("Exiting the SuccessHandler function.")

	data := map[string]interface{}{
		"ContentTemplateName": "success",
	}

	RenderTemplateWithData(w, "base.gohtml", data)
}
