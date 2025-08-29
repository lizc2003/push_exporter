package logic

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func NewPushHandler(mgr *MetricMgr) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "body empty", http.StatusBadRequest)
			return
		}

		data, err := io.ReadAll(req.Body)
		if err != nil || len(data) == 0 {
			http.Error(w, "read body fail", http.StatusBadRequest)
			return
		}

		var metrics []*MetricValue
		if data[0] == '[' {
			err = json.Unmarshal(data, &metrics)
			if err != nil {
				http.Error(w, "decode body fail", http.StatusBadRequest)
				return
			}
		} else {
			m := &MetricValue{}
			err = json.Unmarshal(data, m)
			if err != nil {
				http.Error(w, "decode body fail", http.StatusBadRequest)
				return
			}
			metrics = []*MetricValue{m}
		}

		now := time.Now().Unix()
		for _, m := range metrics {
			mgr.AddMetric(m, now)
		}

		w.Write([]byte("OK"))
	}
}
