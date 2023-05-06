package loggerfrontend

import (
	"encoding/json"
	"net/http"

	"github.com/rmacdiarmid/gptback/logger"
	"github.com/rmacdiarmid/gptback/pkg/database"
)

func Log(w http.ResponseWriter, r *http.Request) {
	var logEntry database.FrontendLog
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.InsertFrontendLog(logEntry)
	if err != nil {
		logger.DualLog.Printf("Error inserting frontend log: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
