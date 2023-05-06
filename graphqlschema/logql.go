package graphqlschema

import (
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/rmacdiarmid/gptback/pkg/database"
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
		"timestamp": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var CreateFrontendLogField = &graphql.Field{
	Type: FrontendLogType,
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

var ReadFrontendLogField = &graphql.Field{
	Type: FrontendLogType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		if id, ok := params.Args["id"].(int); ok {
			frontendLog, err := database.GetFrontendLogByID(strconv.Itoa(id))
			if err != nil {
				return nil, err
			}
			return frontendLog, nil
		} else {
			frontendLogs, err := database.GetAllFrontendLogs()
			if err != nil {
				return nil, err
			}
			return frontendLogs, nil
		}
	},
}

var UpdateFrontendLogField = &graphql.Field{
	Type: FrontendLogType,
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
		id, _ := params.Args["id"].(int)
		message, _ := params.Args["message"].(string)
		timestamp, _ := params.Args["timestamp"].(string)

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

var DeleteFrontendLogField = &graphql.Field{
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
