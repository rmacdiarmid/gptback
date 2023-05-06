package internal

import (
	"net/http"

	"github.com/rmacdiarmid/gptback/logger"
	"github.com/rmacdiarmid/gptback/pkg/database"
)

func TaskListHandler(w http.ResponseWriter, r *http.Request) {
	// Log message for starting TaskListHandler function
	logger.DualLog.Println("Starting TaskListHandler function...")
	defer logger.DualLog.Println("All tasks read successfully.")

	// Get all tasks from the database
	tasks, err := database.Rgithub.com/rmacdiarmid/gptback
		logger.DualLog.Printf("Error reading all tasks: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Tasks": tasks,
	}

	RenderTemplateWithData(w, "base.gohtml", "taskListContent", data)

	// Log message for successful completion of TaskListHandler function
	logger.DualLog.Println("TaskListHandler function completed successfully.")
}
