package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/rmacdiarmid/GPTSite/graphqlschema"
	"github.com/rmacdiarmid/GPTSite/pkg/database" // Add this line to import the schema

	// Replace the previous import with this one
	"github.com/stretchr/testify/assert"
)

func TestGraphQL(t *testing.T) {
	// Initialize your database here with test data
	_, err := database.InitDB(":memory:")
	assert.Nil(t, err, "Failed to initialize database")

	articleID, err := database.CreateArticle("Test title", "Test image", "Test preview")
	assert.Nil(t, err, "Failed to create test article")

	// GraphQL query to fetch the test data
	query := `
	{
		article(id: 1) {
			id
			title
			image
			preview
		}
	}
`

	fmt.Printf("Schema: %+v\n", graphqlschema.Schema)

	params := graphql.Params{Schema: graphqlschema.Schema, RequestString: query}
	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		t.Logf("Errors: %v", result.Errors)
	}
	t.Logf("Result: %+v", result)
	t.Logf("Result data: %v", result.Data)

	// Add this line to print the result
	fmt.Printf("Result: %+v\n", result)

	// Check if there are any errors
	assert.Empty(t, result.Errors, "GraphQL query returned errors")

	// Compare the expected output with the actual result
	expected := map[string]interface{}{
		"article": map[string]interface{}{
			"id":      1,
			"image":   "Test image",
			"preview": "Test preview",
			"title":   "Test title",
		},
	}

	fmt.Printf("Result data: %+v\n", result.Data)
	assert.Equal(t, expected, result.Data, "GraphQL query result doesn't match expected output")

	// Cleanup: Remove test data from the database
	err = database.DeleteArticle(articleID)
	assert.Nil(t, err, "Failed to delete test article")
}

func TestGraphQLArticlesQuery(t *testing.T) {
	// Insert articles into the database
	articles := []database.Article{
		{Title: "Test Article 1", Image: "test-image-1.jpg", Preview: "This is test article 1"},
		{Title: "Test Article 2", Image: "test-image-2.jpg", Preview: "This is test article 2"},
		{Title: "Test Article 3", Image: "test-image-3.jpg", Preview: "This is test article 3"},
	}

	for _, article := range articles {
		_, err := database.CreateArticle(article.Title, article.Image, article.Preview)
		if err != nil {
			t.Fatalf("Error creating test article: %v", err)
		}
	}

	// Define the articles query
	query := `
		{
			articles {
				id
				title
				image
				preview
			}
		}
	`

	params := graphql.Params{Schema: graphqlschema.Schema, RequestString: query}
	result := graphql.Do(params)

	if len(result.Errors) > 0 {
		t.Fatalf("Unexpected errors: %v", result.Errors)
	}

	// Extract the response data
	data := result.Data
	articlesData := data.(map[string]interface{})["articles"].([]interface{})

	if len(articlesData) != len(articles) {
		t.Fatalf("Expected %d articles, got %d", len(articles), len(articlesData))
	}

	// Compare the response data with the original data
	for i, articleData := range articlesData {
		articleMap := articleData.(map[string]interface{})
		if articleMap["title"] != articles[i].Title ||
			articleMap["image"] != articles[i].Image ||
			articleMap["preview"] != articles[i].Preview {
			t.Errorf("Article data does not match: expected %+v, got %+v", articles[i], articleMap)
		}
	}
}

func TestGraphQLCreateArticleMutation(t *testing.T) {
	mutation := `
        mutation {
            createArticle(title: "New Test Article", image: "new-test-image.jpg", preview: "This is a new test article") {
                id
                title
                image
                preview
            }
        }
    `

	params := graphql.Params{Schema: graphqlschema.Schema, RequestString: mutation}
	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		t.Fatalf("Failed to execute GraphQL mutation: %v", result.Errors)
	}

	expected := map[string]interface{}{
		"title":   "New Test Article",
		"image":   "new-test-image.jpg",
		"preview": "This is a new test article",
	}

	data, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Failed to type assert result.Data to map[string]interface{}")
	}

	actual := data["createArticle"].(map[string]interface{})
	delete(actual, "id")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GraphQL mutation result doesn't match expected output\nexpected: %v\nactual: %v", expected, actual)
	}
}

func TestGraphQLUpdateArticleMutation(t *testing.T) {
	articleID, err := database.CreateArticle("Test title", "Test image", "Test preview")
	assert.Nil(t, err, "Failed to create test article")

	mutation := fmt.Sprintf(`
		mutation {
			updateArticle(id: %d, title: "Updated Test Article", image: "updated-test-image.jpg", preview: "This is an updated test article") {
				id
				title
				image
				preview
			}
		}
	`, articleID)

	params := graphql.Params{Schema: graphqlschema.Schema, RequestString: mutation}
	result := graphql.Do(params)

	assert.Empty(t, result.Errors, "GraphQL mutation returned errors")

	expected := map[string]interface{}{
		"updateArticle": map[string]interface{}{
			"id":      6,
			"title":   "Updated Test Article",
			"image":   "updated-test-image.jpg",
			"preview": "This is an updated test article",
		},
	}

	assert.Equal(t, expected, result.Data, "GraphQL mutation result doesn't match expected output")
}

func TestGraphQLDeleteArticleMutation(t *testing.T) {
	articleID, err := database.CreateArticle("Test title", "Test image", "Test preview")
	assert.Nil(t, err, "Failed to create test article")

	mutation := fmt.Sprintf(`
		mutation {
			deleteArticle(id: %d)
		}
	`, articleID)

	params := graphql.Params{Schema: graphqlschema.Schema, RequestString: mutation}
	result := graphql.Do(params)

	assert.Empty(t, result.Errors, "GraphQL mutation returned errors")

	expected := map[string]interface{}{
		"deleteArticle": true,
	}

	assert.Equal(t, expected, result.Data, "GraphQL mutation result doesn't match expected output")
}
