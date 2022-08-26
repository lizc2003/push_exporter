package logic

import (
	"encoding/json"
	"net/http"
	"time"
)

func NewPushHandler(mgr *MetricMgr) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "body empty", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var metrics []*MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "decode body fail", http.StatusBadRequest)
			return
		}

		now := time.Now().Unix()
		for _, m := range metrics {
			mgr.AddMetric(m, now)
		}

		w.Write([]byte("OK"))
	}
}
