package handlers

import (
	"log"
	"net/http"
)

type Article struct {
	Title   string
	Image   string
	Preview string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("IndexHandler called")

	// Retrieve data from the database
	articles := []Article{
		{Title: "The Whimsical Dragon's Guide to the Meaning of Life", Image: "/static/images/dragon1.jpg", Preview: "Once upon a time, in a land far, far away..."},
		{Title: "Dragon Tales and the Meaning of Life", Image: "/static/images/dragon2.jpg", Preview: "Dragons are known for their powerful and often misunderstood nature..."},
		{Title: "The Treasure of Dragon Island and the Meaning of Life", Image: "/static/images/dragon3.jpg", Preview: "Deep in the heart of the ocean lay a mysterious island known only to the most daring of dragons..."},
	}

	log.Printf("Rendering index with data: %#v", articles)

	// Render the HTML template with the data
	renderTemplateWithData(w, "index.html", articles)
}
