package handlers

import (
	"encoding/json"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
	"strconv"
)

type RespGetUrls struct {
	Short string `json:"short_url"`
	Long  string `json:"original_url"`
}

func HandlerGetUrls(w http.ResponseWriter, r *http.Request, baseAddress string, db *storage.StorageBase) {
	_, err := r.Cookie("Authorization")
	if err != nil {
		logger.Error("cookies do not contain a token: %v", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID := r.Header.Get("User_id")
	id, _ := strconv.Atoi(userID)

	urls, err := db.GetUrls(r.Context(), id)
	if err != nil {
		logger.Error("cannot get urls: %v", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var dataResp []RespGetUrls

	for _, url := range urls {
		dataResp = append(dataResp, RespGetUrls{
			Long:  url.Long,
			Short: baseAddress + "/" + url.Short,
		})
	}

	resp, err := json.Marshal(dataResp)
	if err != nil {
		logger.Error("PostAPI not json to byte :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
