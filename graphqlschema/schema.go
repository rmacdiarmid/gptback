package graphqlschema

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
)

var Schema graphql.Schema

func init() {
	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: Query,
	})
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
	},
})
