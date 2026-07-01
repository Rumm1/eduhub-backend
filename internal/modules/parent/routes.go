package parent

import "net/http"

func RegisterRoutes(mux *http.ServeMux, basePath string, handler *Handler) {
	if mux == nil {
		return
	}
	if basePath == "" {
		basePath = "/parent"
	}
	if handler == nil {
		handler = NewHandler(nil)
	}
	mux.HandleFunc("GET "+basePath, handler.List)
}
