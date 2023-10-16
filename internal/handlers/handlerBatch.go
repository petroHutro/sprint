package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
)

type DataReqBatch struct {
	ID   string `json:"correlation_id"`
	Long string `json:"original_url"`
}

type DataRespBatch struct {
	ID    string `json:"correlation_id"`
	Short string `json:"short_url"`
}

func HandlerPostBatch(w http.ResponseWriter, r *http.Request, baseAddress, file string, db *storage.StorageBase) {
	var buf bytes.Buffer
	var data []DataReqBatch

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("PostBatch not body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		logger.Log.Error("PostBatch not byte to json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var longs []string

	for _, item := range data {
		if item.ID == "" || item.Long == "" {
			logger.Log.Error("PostBatch not correlation_id or original_url: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		longs = append(longs, item.Long)
	}

	if err := db.SetAllDB(r.Context(), longs); err != nil {
		logger.Log.Error("cannot set all: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var dataResp []DataRespBatch

	for _, item := range data {
		dataResp = append(dataResp, DataRespBatch{
			ID:    item.ID,
			Short: baseAddress + "/" + db.GetShort(r.Context(), item.Long)})
	}
	resp, err := json.Marshal(dataResp)
	if err != nil {
		logger.Log.Error("PostAPI not json to byte :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
