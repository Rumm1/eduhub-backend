package organization

import "net/http"

func RegisterRoutes(mux *http.ServeMux, basePath string, handler *Handler) {
	if mux == nil {
		return
	}
	if basePath == "" {
		basePath = "/organization"
	}
	if handler == nil {
		handler = NewHandler(nil)
	}
	mux.HandleFunc("GET "+basePath, handler.List)
}
