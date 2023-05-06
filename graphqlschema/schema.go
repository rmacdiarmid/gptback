package graphqlschema

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/rmacdiarmid/gptback/pkg/database"
)

var Schema graphql.Schema

// Change the init function to an exported function
func InitSchema() {
	var err error
	Schema, err = graphql.NewSchema(SchemaConfig)
	if err != nil {
		panic(err)
	}
}

var Query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"article": &graphql.Field{
			Type:        ArticleType,
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
			Type:        graphql.NewList(ArticleType),
			Description: "List of articles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				articles, err := database.GetArticles()
				if err != nil {
					return nil, err
				}
				return articles, nil
			},
		},
		"frontendLog": ReadFrontendLogField,
		"frontendLogs": &graphql.Field{
			Type:        graphql.NewList(FrontendLogType),
			Description: "List of frontend logs",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				frontendLogs, err := database.GetAllFrontendLogs()
				if err != nil {
					return nil, err
				}
				return frontendLogs, nil
			},
		},
	},
})

var Mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createArticle": createArticleMutationField,
		"updateArticle": &graphql.Field{
			Type:        ArticleType,
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
				"text": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				title, _ := p.Args["title"].(string)
				image, _ := p.Args["image"].(string)
				preview, _ := p.Args["preview"].(string)
				text, _ := p.Args["text"].(string)

				updatedArticle, err := database.UpdateArticle(int64(id), title, image, preview, text)
				if err != nil {
					return nil, err
				}
				return updatedArticle, nil
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
		"createFrontendLog": CreateFrontendLogField,
		"updateFrontendLog": UpdateFrontendLogField,
		"deleteFrontendLog": DeleteFrontendLogField,
	},
})

// Update the SchemaConfig to include the Mutation
var SchemaConfig = graphql.SchemaConfig{
	Query:    Query,
	Mutation: Mutation,
}
