package graphqlschema

import (
	"github.com/graphql-go/graphql"
	"github.com/rmacdiarmid/gptback/internal"
	"github.com/rmacdiarmid/gptback/pkg/database"
)

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"userId": &graphql.Field{
			Type: graphql.Int,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var RegisterInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "RegisterInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"email": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"password": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"passwordConfirmation": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var LoginInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "LoginInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"email": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"password": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var RegisterMutation = &graphql.Field{
	Type: UserType,
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(RegisterInputType),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		input := p.Args["input"].(map[string]interface{})
		newUser := database.NewUser{
			Email:                input["email"].(string),
			Password:             input["password"].(string),
			PasswordConfirmation: input["passwordConfirmation"].(string),
		}
		userID, err := internal.RegisterUser(newUser)
		if err != nil {
			return nil, err
		}
		user, err := database.GetUserByID(userID)
		return user, err
	},
}

var LoginMutation = &graphql.Field{
	Type: graphql.String,
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(LoginInputType),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		input := p.Args["input"].(map[string]interface{})
		token, err := internal.LoginUser(input)
		return token, err
	},
}
