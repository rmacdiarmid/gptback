// graphqlschema/log.go

package graphqlschema

import (
	"github.com/graphql-go/graphql"
)

var FrontendLogType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FrontendLog",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"message": &graphql.Field{
			Type: graphql.String,
		},
		// Add other fields as needed
	},
})
