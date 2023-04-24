package internal

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rmacdiarmid/GPTSite/logger"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
)

func CreateFrontendLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var logEntry database.FrontendLog
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err != nil {
		logger.DualLog.Printf("Error decoding frontend log: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = database.InsertFrontendLog(logEntry)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Read a frontend log by ID
func ReadFrontendLogHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	id := mux.Vars(r)["id"]

	// Fetch the log from the database
	logEntry, err := database.GetFrontendLogByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Return the log entry as JSON
	json.NewEncoder(w).Encode(logEntry)
}

// Update a frontend log by ID
func UpdateFrontendLogHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	id := mux.Vars(r)["id"]

	// Decode the log entry from the request body
	var updatedLogEntry database.FrontendLog
	err := json.NewDecoder(r.Body).Decode(&updatedLogEntry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the log entry in the database
	err = database.UpdateFrontendLogByID(id, updatedLogEntry)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return a success status
	w.WriteHeader(http.StatusOK)
}

// Delete a frontend log by ID
func DeleteFrontendLogHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	id := mux.Vars(r)["id"]

	// Delete the log entry from the database
	err := database.DeleteFrontendLogByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return a success status
	w.WriteHeader(http.StatusOK)
}
