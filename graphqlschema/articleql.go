package graphqlschema

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
	"github.com/spf13/viper"
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
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					article, ok := p.Source.(database.Article)
					if !ok {
						return nil, fmt.Errorf("expected type database.Article but got %T", p.Source)
					}
					baseURL := viper.GetString("storage.baseURL")
					imageURL := fmt.Sprintf("%s%s", baseURL, article.Image)
					return imageURL, nil
				},
			},
			"preview": &graphql.Field{
				Type: graphql.String,
			},
			"text": &graphql.Field{ // Make sure this field is included
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
		"text": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		title, _ := params.Args["title"].(string)
		image, _ := params.Args["image"].(string)
		preview, _ := params.Args["preview"].(string)
		text, _ := params.Args["text"].(string)

		newArticleID, err := database.CreateArticle(title, image, preview, text) // Add the 'text' parameter here
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

var updateArticleMutationField = &graphql.Field{
	Type: ArticleType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"title": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"image": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"preview": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"text": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		id, _ := params.Args["id"].(int)
		title, _ := params.Args["title"].(string)
		image, _ := params.Args["image"].(string)
		preview, _ := params.Args["preview"].(string)
		text, _ := params.Args["text"].(string)

		// Change the following line to handle both returned values
		updatedArticle, err := database.UpdateArticle(int64(id), title, image, preview, text) // Add the 'text' parameter here
		if err != nil {
			return nil, err
		}

		return updatedArticle, nil
	},
}
