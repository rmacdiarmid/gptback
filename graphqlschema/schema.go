package graphqlschema

import (
	"fmt"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
)

var Schema graphql.Schema

func init() {
	var err error
	Schema, err = graphql.NewSchema(SchemaConfig)
	if err != nil {
		panic(err)
	}
}

var articleType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Article",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"image": &graphql.Field{
				Type: graphql.String,
			},
			"preview": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// Define the FrontendLog type
var frontendLogType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "FrontendLog",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"message": &graphql.Field{
				Type: graphql.String,
			},
			"timestamp": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// Define the frontend log resolvers
var createFrontendLogField = &graphql.Field{
	Type: frontendLogType,
	Args: graphql.FieldConfigArgument{
		"message": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"timestamp": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		message, _ := params.Args["message"].(string)
		timestamp, _ := params.Args["timestamp"].(string)

		parsedTimestamp, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return nil, err
		}

		newFrontendLogID, err := database.InsertFrontendLog(database.FrontendLog{Message: message, Timestamp: parsedTimestamp})
		if err != nil {
			return nil, err
		}

		newFrontendLog, err := database.GetFrontendLogByID(strconv.FormatInt(newFrontendLogID, 10))
		if err != nil {
			return nil, err
		}

		return newFrontendLog, nil
	},
}

var readFrontendLogField = &graphql.Field{
	Type: frontendLogType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.Int, // Remove graphql.NewNonNull()
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		// Check if the ID argument is provided
		if id, ok := params.Args["id"].(int); ok {
			frontendLog, err := database.GetFrontendLogByID(strconv.Itoa(id))
			if err != nil {
				return nil, err
			}
			return frontendLog, nil
		} else {
			// Return all frontend logs if the ID argument is not provided
			frontendLogs, err := database.GetAllFrontendLogs()
			if err != nil {
				return nil, err
			}
			return frontendLogs, nil
		}
	},
}

var updateFrontendLogField = &graphql.Field{
	Type: frontendLogType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"message": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"timestamp": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		message, ok := params.Args["message"].(string)
		if !ok {
			return nil, fmt.Errorf("message should be a string")
		}

		timestamp, ok := params.Args["timestamp"].(string)
		if !ok {
			return nil, fmt.Errorf("timestamp should be a string")
		}

		id, _ := params.Args["id"].(int)
		currentLogEntry, err := database.GetFrontendLogByID(strconv.Itoa(id))
		if err != nil {
			return nil, err
		}

		if message != "" {
			currentLogEntry.Message = message
		}
		if timestamp != "" {
			parsedTimestamp, err := time.Parse(time.RFC3339, timestamp)
			if err != nil {
				return nil, err
			}
			currentLogEntry.Timestamp = parsedTimestamp
		}
		err = database.UpdateFrontendLogByID(strconv.Itoa(id), currentLogEntry)
		if err != nil {
			return nil, err
		}

		updatedLogEntry, err := database.GetFrontendLogByID(strconv.Itoa(id))
		if err != nil {
			return nil, err
		}

		return updatedLogEntry, nil
	},
}

var deleteFrontendLogField = &graphql.Field{
	Type: graphql.Boolean,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		id, _ := params.Args["id"].(int)

		err := database.DeleteFrontendLogByID(strconv.Itoa(id))
		if err != nil {
			return nil, err
		}

		return true, nil
	},
}

var Query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"article": &graphql.Field{
			Type:        articleType,
			Description: "Get single article by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					article, err := database.ReadArticle(int64(id))
					if err != nil {
						return nil, err
					}
					fmt.Printf("Resolver: article: %+v\n", article) // Add this line
					return article, nil
				}
				return nil, nil
			},
		},
		"articles": &graphql.Field{
			Type:        graphql.NewList(articleType),
			Description: "List of articles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				articles, err := database.GetArticles()
				if err != nil {
					return nil, err
				}
				return articles, nil
			},
		},
		"frontendLog": readFrontendLogField,
	},
})

var createArticleMutationField = &graphql.Field{
	Type: articleType,
	Args: graphql.FieldConfigArgument{
		"title": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"image": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"preview": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		title, _ := params.Args["title"].(string)
		image, _ := params.Args["image"].(string)
		preview, _ := params.Args["preview"].(string)

		newArticleID, err := database.CreateArticle(title, image, preview)
		if err != nil {
			return nil, err
		}

		newArticle, err := database.ReadArticle(newArticleID)
		if err != nil {
			return nil, err
		}

		return newArticle, nil
	},
}

var Mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createArticle": createArticleMutationField,
		"updateArticle": &graphql.Field{
			Type:        articleType,
			Description: "Update an existing article",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"title": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"image": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"preview": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				title, _ := p.Args["title"].(string)
				image, _ := p.Args["image"].(string)
				preview, _ := p.Args["preview"].(string)

				article, err := database.UpdateArticle(int64(id), title, image, preview)
				if err != nil {
					return nil, err
				}
				return article, nil
			},
		},
		"deleteArticle": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete an article by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)

				err := database.DeleteArticle(int64(id))
				if err != nil {
					return nil, err
				}
				return true, nil
			},
		},
		"createFrontendLog": createFrontendLogField,
		"updateFrontendLog": updateFrontendLogField,
		"deleteFrontendLog": deleteFrontendLogField,
	},
})

// Update the SchemaConfig to include the Mutation
var SchemaConfig = graphql.SchemaConfig{
	Query:    Query,
	Mutation: Mutation,
}
