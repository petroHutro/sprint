package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
	"strconv"
)

type DataReqAPI struct {
	URL string `json:"url"`
}

type DataRespAPI struct {
	Result string `json:"result"`
}

func HandlerPostAPI(w http.ResponseWriter, r *http.Request, baseAddress, file string, db *storage.StorageBase) {
	var buf bytes.Buffer
	var data DataReqAPI

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Error("PostAPI not body :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		logger.Error("PostAPI not byte to json :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if data.URL == "" {
		logger.Error("PostAPI not url :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.Header.Get("User_id"))
	if err != nil {
		logger.Error("bad user id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = db.LongToShort(r.Context(), data.URL, file, id)
	statusCode := http.StatusCreated
	if err != nil {
		var repErr *storage.RepError
		if errors.As(err, &repErr) && repErr.Repetition {
			statusCode = http.StatusConflict
			logger.Error("long already db :%v", err)
		} else {
			logger.Error("cannot convert long to short :%v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	short, _ := db.GetShort(r.Context(), data.URL) //!!!!!!!!!!!!!!!!!!!!
	dataResp := DataRespAPI{Result: baseAddress + "/" + short}
	resp, err := json.Marshal(dataResp)
	if err != nil {
		logger.Error("PostAPI not json to byte :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}
