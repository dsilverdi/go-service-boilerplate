package server

import (
	"encoding/json"
	"go-service-boilerplate/pkg/rest"
	"net/http"
)

func (svc *AppRestService) HelloHandler(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(&rest.Response{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    "Hello World Handler",
	})
}
