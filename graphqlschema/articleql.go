package graphqlschema

import (
	"github.com/graphql-go/graphql"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
)

var ArticleType = graphql.NewObject(
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

var createArticleMutationField = &graphql.Field{
	Type: ArticleType,
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
