package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-chi/chi/middleware"

	"github.com/sirupsen/logrus"

	"github.com/undeadops/konveyer/pkg"
)


// API - Server Runtime Properties 
type API struct {
	log        *logrus.Entry
	logger     *logrus.Logger
	config     *root.Config
	port       string
	router     *chi.Mux
	// db         db.DataSource
	//payerProxy payerProxy
	// version    string
}

// type Server struct {
//   router *mux.Router
//   config *root.ServerConfig
// }

// NewAPI - Setup HTTP Server and Routes
func NewAPI(d root.DeploymentService, config *root.Config) *API {
	api := &API{
		log:        logrus.WithField("component", "api"),
		logger:     config.Logger,
		config:     config,
		port:       config.Port,
		// db:         db,
		// version:    version,
		router:     chi.NewRouter(),
	}

	
	api.router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.RequestID,
		NewStructuredLogger(api.logger),
		middleware.URLFormat,
		middleware.DefaultCompress,
		middleware.Recoverer,
	)

  	api.router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		Json(w, http.StatusOK, map[string]string{
			"version": "0.1",
			"application": "konveyer",
		})
  	})
	api.router.Mount("/v1/deployment", NewDeploymentRouter(d))
	return api
} 


// func(s *Server) Start() {
//   log.Println("Listening on port " + s.config.Port)
//   if err := http.ListenAndServe(s.config.Port, handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
//       log.Fatal("http.ListenAndServe: ", err)
//   }
// }

// Start - Startup WebServer
func (a *API) Start() error {
	a.log.Infof("Konveyer API started on: %s", a.port)
	return http.ListenAndServe(a.port, a.router)
}

// Server - Configuration 
// type Server struct {
// 	DB     
// 	Router chi.Router
// 	Logger *logrus.Logger
// }

// func (a *API) hello(w http.ResponseWriter, r *http.Request) {
// 	sendJSON(w, http.StatusOK, map[string]string{
// 		"version":     a.version,
// 		"application": "konveyer",
// 	})
// }

// func sendJSON(w http.ResponseWriter, status int, obj interface{}) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	encoder := json.NewEncoder(w)
// 	encoder.Encode(obj)
// }

