package types

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Channel string `json:"X-Goog-Channel-ID"`
	MessageNumber int64 `json:"X-Goog-Message-Number"`
	ResourceID string `json:"X-Goog-Resource-ID"`
	ResourceState string `json:"X-Goog-Resource-State"`
	ResourceLink string `json:"X-Goog-Resource-URI"`
	Changed string `json:"X-Goog-Changed"`
	Expiration string `json:"X-Goog-Channel-Expiration"`
	Token string `json:"X-Goog-Channel-Token"`
}

type ErrorResponse struct {
	Status int    `json:"status"`
}

func (r *Request) Unmarshal(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (e *ErrorResponse) Send(w http.ResponseWriter) {
	_ = json.NewEncoder(w).Encode(e)
}