package handlers

import (
	"encoding/json"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
	"strconv"
)

func HandlerGetUrls(w http.ResponseWriter, r *http.Request, db *storage.StorageBase) {
	userID := r.Header.Get("User_id")
	id, _ := strconv.Atoi(userID)

	urls, err := db.GetUrls(r.Context(), id)
	if err != nil {
		logger.Error("cannot get urls: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var dataResp []storage.Urls

	for _, url := range urls {
		dataResp = append(dataResp, storage.Urls{
			Long:  url.Long,
			Short: url.Short,
		})
	}

	resp, err := json.Marshal(dataResp)
	if err != nil {
		logger.Error("PostAPI not json to byte :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
