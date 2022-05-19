package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Init - routes initialization
func Init(router *chi.Mux) {
	h := InitHandler()
	router.HandleFunc("/index", h.index)

	router.Route("/server/{host}", func(r chi.Router) {
		r.Get("/services", CorsHandler(h.getLists))
		r.Post("/services", CorsHandler(h.getListsWithProto))
		r.Get("/service/{serv_name}/functions", CorsHandler(h.getLists))
		r.Get("/function/{func_name}/describe", CorsHandler(h.describeFunction))
		r.Post("/function/{func_name}/invoke", CorsHandler(h.invokeFunction))
	})

	// get list of active connection
	router.Get("/active/get", CorsHandler(h.getActiveConns))
	// close active connection
	router.Delete("/active/close/{host}", CorsHandler(h.closeActiveConns))

	router.Get("/api/request/{name}", CorsHandler(h.getRequest))
	router.Post("/api/request", CorsHandler(h.saveRequest))
	router.Put("/api/request/{name}", CorsHandler(h.updateRequest))
	router.Delete("/api/request/{name}", CorsHandler(h.delRequest))
	router.Get("/api/requests", CorsHandler(h.listRequest))
	router.Get("/api/requests:download", CorsHandler(h.downloadAllRequest))

	assetsPath := "index"
	router.Mount("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(assetsPath+"/css/"))))
	router.Mount("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(assetsPath+"/js/"))))
	router.Mount("/font/", http.StripPrefix("/font/", http.FileServer(http.Dir(assetsPath+"/font/"))))
	router.Mount("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(assetsPath+"/img/"))))
}

func CorsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Headers", "use_tls")
			return
		}

		h.ServeHTTP(w, r)
	}
}
