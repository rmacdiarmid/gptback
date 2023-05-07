package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/rmacdiarmid/gptback/graphqlschema"
	"github.com/rmacdiarmid/gptback/logger"
	"github.com/rmacdiarmid/gptback/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Initialize the logger
	logger.InitLogger(os.Stdout)

	// Initialize the in-memory database for testing
	database.InitDB(":memory:")

	graphqlschema.InitSchema()
	os.Exit(m.Run())
}

func convertID(id interface{}) (int, error) {
	switch v := id.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("unexpected type %T for ID", v)
	}
}

func TestGraphQL(t *testing.T) {
	// Initialize your database here with test data
	_, err := database.InitDB(":memory:")
	assert.Nil(t, err, "Failed to initialize database")

	articleID, err := database.CreateArticle("Test title", "Test image", "Test preview", "Test text")
	assert.Nil(t, err, "Failed to create test article")

	// GraphQL query to fetch the test data
	query := fmt.Sprintf(`
	{
		article(id: %d) {
			id
			title
			image
			preview
			text
		}
	}
`, articleID)

	params := graphql.Params{Schema: graphqlschema.Schema, RequestString: query}
	result := graphql.Do(params)

	// Check if there are any errors
	assert.Empty(t, result.Errors, "GraphQL query returned errors")

	// Extract the response data
	data := result.Data.(map[string]interface{})
	article := data["article"].(map[string]interface{})

	// Convert the ID
	id, err := convertID(article["id"])
	if err != nil {
		t.Errorf("Error converting ID: %v", err)
	}
	article["id"] = id

	// Compare the expected output with the actual result
	expected := map[string]interface{}{
		"article": map[string]interface{}{
			"id":      1,
			"title":   "Test title",
			"image":   "Test image",
			"preview": "Test preview",
			"text":    "Test text",
		},
	}

	assert.Equal(t, expected, data, "GraphQL query result doesn't match expected output")

	// Cleanup: Remove test data from the database
	err = database.DeleteArticle(int64(articleID))
	assert.Nil(t, err, "Failed to delete test article")
}
func TestGraphQLArticlesQuery(t *testing.T) {
	// Insert articles into the database
	articles := []database.Article{
		{Title: "Test Article 1", Image: "test-image-1.jpg", Preview: "This is test preview of article 1", Text: "This is test article 1"},
		{Title: "Test Article 2", Image: "test-image-2.jpg", Preview: "This is test preview of article 2", Text: "This is test article 2"},
		{Title: "Test Article 3", Image: "test-image-3.jpg", Preview: "This is test preview of article 3", Text: "This is test article 3"},
	}

	for _, article := range articles {
		articleID, err := database.CreateArticle(article.Title, article.Image, article.Preview, article.Text)
		if err != nil {
			logger.DualLog.Printf("Error creating test article: %v", err)
			t.FailNow()
		}
		t.Cleanup(func() {
			_ = database.DeleteArticle(articleID)
		})
	}

	// Define the articles query
	query := `
		{
			articles {
				id
				title
				image
				preview
				text
			}
		}
	`

	params := graphql.Params{Schema: graphqlschema.Schema, RequestString: query}
	result := graphql.Do(params)

	if len(result.Errors) > 0 {
		logger.DualLog.Printf("Unexpected errors: %v", result.Errors)
		t.FailNow()
	}

	// Extract the response data
	data := result.Data
	articlesData := data.(map[string]interface{})["articles"].([]interface{})

	if len(articlesData) != len(articles) {
		logger.DualLog.Printf("Expected %d articles, got %d", len(articles), len(articlesData))
		t.FailNow()
	}

	// Compare the response data with the original data
	for i, articleData := range articlesData {
		articleMap := articleData.(map[string]interface{})

		// Convert the ID
		id, err := convertID(articleMap["id"])
		if err != nil {
			t.Errorf("Error converting ID: %v", err)
		}
		articleMap["id"] = id

		if articleMap["title"] != articles[i].Title ||
			articleMap["image"] != articles[i].Image ||
			articleMap["preview"] != articles[i].Preview ||
			articleMap["text"] != articles[i].Text {
			t.Errorf("Article data does not match: expected %+v, got %+v", articles[i], articleMap)
		}
	}
}
func TestGraphQLCreateArticleMutation(t *testing.T) {
	mutation := `
        mutation {
            createArticle(title: "New Test Article", image: "new-test-image.jpg", preview: "This is a new test article", text: "This is a new test article") {
                id
                title
                image
                preview
				text
            }
        }
    `
	params := graphql.Params{Schema: graphqlschema.Schema, RequestString: mutation}
	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		logger.DualLog.Printf("Failed to execute GraphQL mutation: %v", result.Errors)
		t.FailNow()
	}

	expected := map[string]interface{}{
		"title":   "New Test Article",
		"image":   "new-test-image.jpg",
		"preview": "This is a new test article",
		"text":    "This is a new test article",
	}

	data, ok := result.Data.(map[string]interface{})
	if !ok {
		logger.DualLog.Printf("Failed to type assert result.Data to map[string]interface{}")
		t.FailNow()
	}

	actual := data["createArticle"].(map[string]interface{})

	// Convert the ID
	id, err := convertID(actual["id"])
	if err != nil {
		t.Errorf("Error converting ID: %v", err)
	}
	actual["id"] = id

	articleID := id
	t.Cleanup(func() {
		_ = database.DeleteArticle(int64(articleID))
	})
	delete(actual, "id")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GraphQL mutation result doesn't match expected output\nexpected: %v\nactual: %v", expected, actual)
	}
}

func TestGraphQLUpdateArticleMutation(t *testing.T) {
	articleID, err := database.CreateArticle("Test title", "Test image", "Test preview", "Test text")
	assert.Nil(t, err, "Failed to create test article")
	t.Cleanup(func() {
		_ = database.DeleteArticle(int64(articleID)) // Convert the int to int64
	})

	mutation := fmt.Sprintf(`
		mutation {
			updateArticle(id: %d, title: "Updated Test Article", image: "updated-test-image.jpg", preview: "This is an updated test article", text: "This is an updated test article") {
				id
				title
				image
				preview
				text
			}
		}
	`, articleID)

	params := graphql.Params{Schema: graphqlschema.Schema, RequestString: mutation}
	result := graphql.Do(params)

	assert.Empty(t, result.Errors, "GraphQL mutation returned errors")

	// Handle int, int64, and float64 for the "id" field
	dataMap, ok := result.Data.(map[string]interface{})
	assert.True(t, ok, "Failed to assert result.Data as map[string]interface{}")

	idValue, ok := dataMap["updateArticle"].(map[string]interface{})["id"].(int)
	if ok {
		dataMap["updateArticle"].(map[string]interface{})["id"] = int64(idValue)
	} else {
		idValueFloat, ok := dataMap["updateArticle"].(map[string]interface{})["id"].(float64)
		if ok {
			dataMap["updateArticle"].(map[string]interface{})["id"] = int64(idValueFloat)
		}
	}

	expected := map[string]interface{}{
		"updateArticle": map[string]interface{}{
			"id":      int64(articleID), // Convert the int to int64
			"title":   "Updated Test Article",
			"image":   "updated-test-image.jpg",
			"preview": "This is an updated test article",
			"text":    "This is an updated test article",
		},
	}

	assert.Equal(t, expected, dataMap, "GraphQL mutation result doesn't match expected output")
}

func TestGraphQLDeleteArticleMutation(t *testing.T) {
	articleID, err := database.CreateArticle("Test title", "Test image", "Test preview", "Test text")
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
