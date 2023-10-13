package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
)

type DataReqAPI struct {
	URL string `json:"url"`
}

type DataRespAPI struct {
	Result string `json:"result"`
}

func HandlerPostAPI(w http.ResponseWriter, r *http.Request, baseAddress, file string, db *storage.StorageBase) {
	// if r.URL.Path != "/api/shorten" || utils.ValidContentType(r.Header.Get("Content-Type"), "application/json") != nil {
	// 	logger.Log.Error("PostAPI not Path or not Content-Type")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	var buf bytes.Buffer
	var data DataReqAPI

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("PostAPI not body :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		logger.Log.Error("PostAPI not byte to json :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if data.URL == "" {
		logger.Log.Error("PostAPI not url :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db.LongToShort(r.Context(), data.URL, file)
	dataResp := DataRespAPI{Result: baseAddress + "/" + db.GetShort(r.Context(), data.URL)}
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
