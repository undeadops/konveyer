package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
	"encoding/json"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"github.com/go-chi/valve"
	"github.com/sirupsen/logrus"

	"github.com/undeadops/konveyer/repo"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

// Env - Runtime Env
type Env struct {
	repo     *repo.Repo
	repoPath string
	stop     chan struct{}
	stopped  chan struct{}
}

func main() {
	var sshkeyFlag = flag.String("ssh_key", "/app/konveyer.key", "Path to SSH Key")
	var pathFlag = flag.String("repo_path", "/code", "Path to store repo data")
	var repoFlag = flag.String("repo", "", "Git Repo as Source of truth")
	var syncSecs = flag.Int("syncsec", 300, "Time between Syncs")
	var gitWebhook = flag.Bool("use-webhook", false, "Use Webhook instead of timed Sync")

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		// disable, as we set our own
		FullTimestamp: true,
	}

	logger.Info("Starting Up....")
	//Gather environment variables that start with CI_
	//var cluster Cluster

	flag.Parse()

	logger.Info("Initializing Git Repo")

	sshKey := func(r *repo.Repo) {
		r.SSHKey = *sshkeyFlag
	}

	gitPath := func(r *repo.Repo) {
		r.Path = *pathFlag
	}

	syncTime := func(r *repo.Repo) {
		r.SyncTime = *syncSecs
	}

	mutex := func(r *repo.Repo) {
		r.Mutex = &sync.RWMutex{}
	}

    l := func(r *repo.Repo) {
		r.Logger = logger
	}

	repo := repo.New(*repoFlag, sshKey, gitPath, syncTime, mutex, l)

	stopChan := make(chan struct{})
	stoppedChan := make(chan struct{})

	env := &Env{repo: repo, repoPath: *pathFlag, stop: stopChan, stopped: stoppedChan}

	// Our graceful valve shut-off package to manage code preemption and
	// shutdown signaling.
	valv := valve.New()
	baseCtx := valv.Context()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Setup the logger backend using sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/sirupsen/logrus


	// HTTP service running in this program as well. The valve context is set
	// as a base context on the server listener at the point where we instantiate
	// the server - look lower.
	r := chi.NewRouter()
	r.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.RequestID,
		NewStructuredLogger(logger),
		middleware.URLFormat,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sup"))
	})

	// r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {

	// 	valve.Lever(r.Context()).Open()
	// 	defer valve.Lever(r.Context()).Close()

	// 	select {
	// 	case <-valve.Lever(r.Context()).Stop():
	// 		fmt.Println("valve is closed. finish up..")

	// 	case <-time.After(5 * time.Second):
	// 		// The above channel simulates some hard work.
	// 		// We want this handler to complete successfully during a shutdown signal,
	// 		// so consider the work here as some background routine to fetch a long running
	// 		// search query to find as many results as possible, but, instead we cut it short
	// 		// and respond with what we have so far. How a shutdown is handled is entirely
	// 		// up to the developer, as some code blocks are preemptable, and others are not.
	// 		time.Sleep(5 * time.Second)
	// 	}

	// 	w.Write([]byte(fmt.Sprintf("all done.\n")))
	// })

	// RESTy routes for "articles" resource
	// r.Route("/articles", func(r chi.Router) {
	// 	r.With(paginate).Get("/", ListArticles)
	// 	r.Post("/", CreateArticle)       // POST /articles
	// 	r.Get("/search", SearchArticles) // GET /articles/search

	// 	r.Route("/{articleID}", func(r chi.Router) {
	// 		r.Use(ArticleCtx)            // Load the *Article on the request context
	// 		r.Get("/", GetArticle)       // GET /articles/123
	// 		r.Put("/", UpdateArticle)    // PUT /articles/123
	// 		r.Delete("/", DeleteArticle) // DELETE /articles/123
	// 	})

	// 	// GET /articles/whats-up
	// 	r.With(ArticleCtx).Get("/{articleSlug:[a-z-]+}", GetArticle)
	// })

	r.Mount("/v1", env.v1Router())

	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if *routes {
		// fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/go-chi/chi",
			Intro:       "Welcome to the chi/_examples/rest generated docs.",
		}))
		return
	}

	// Start git sync
	if *gitWebhook == false {
		go repo.Sync(stopChan, stoppedChan)
	}

	srv := http.Server{Addr: ":5000", Handler: chi.ServerBaseContext(baseCtx, r)}

	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Info("Shutting Down..")
			// first valv
			valv.Shutdown(20 * time.Second)

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			// start http shutdown
			srv.Shutdown(ctx)

			// verify, in worst case call cancel via defer
			select {
			case <-time.After(21 * time.Second):
				logger.Info("Not all connections are closed")
			case <-ctx.Done():
			}

			// shutdown git
			logger.Info("Closing out Git repo")
			close(stopChan)
			<-stoppedChan

			logger.Info("Git Closed")
		}
	}()

	srv.ListenAndServe()
}

func (env *Env) v1Router() chi.Router {
	r := chi.NewRouter()
	r.Route("/event/sync", func(r chi.Router) {
		r.Post("/", env.TriggerRepoSync) // POST /event/sync
	})
	r.Route("/deployment/{namespace}/{appname}", func(r chi.Router) {
		r.Get("/", env.GetDeployApp)                  // GET /deployment/namespace/appname
		r.Put("/image", env.SetDeployImage)   // PUT /deployment/namespace/appname/image -d '{ "container_name": "konveyer", "image": "undeadops/konveyer:master" }
		r.Put("/annotations", env.SetAnnotations)     // PUT /deployment/namespace/appname -d '{ []map[string]string }'
		r.Patch("/annotations", env.PatchAnnotations) // PATCH /deployment/namespace/app -d '{ []map[string]string }'
	})

	return r
}

// TriggerRepoSync - GET /v1/event/sync
func (env *Env) TriggerRepoSync(w http.ResponseWriter, r *http.Request) {
	err := env.repo.PullRepo()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	render.Render(w, r, &App{})
}

// GetDeployApp - GET /v1/deployment/<namespace>/<appname>
func (env *Env) GetDeployApp(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	appname := chi.URLParam(r, "appname")

	containers, err := env.repo.GetDeploymentImage(namespace, appname)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	render.Render(w, r, &App{Namespace: namespace, AppName: appname, Containers: containers})
}

// SetDeployImage - PUT /v1/deployment/<namespace>/<appname>
//curl -X POST -d '{"image":"undeadops/konveyer:master-1234abc", "container":"konveyer"}' http://localhost:3000/v1/deployment/<namespace>/<appname>
func (env *Env) SetDeployImage(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	appname := chi.URLParam(r, "appname")
 
	var d DeployImage 

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	//deployment, err := env.repo.GetDeploymentImage(namespace, appname)
	//if err != nil {
	//	render.Render(w, r, ErrInvalidRequest(err))
	//	return
	//}

	// var msg map[string]interface {}
	// msg["old_image"] = deployment[d.ContainerName]
	// msg["new_image"] = d.Image

	//LogEntrySetFields(r, msg)
	// TODO: Handle multiple containers somehow... 
	err := env.repo.SetDeploymentImage(namespace, appname, d.Image)
	if err != nil {
		render.Render(w, r, ErrServerUnable(err))
		return
	}
	containers, err := env.repo.GetDeploymentImage(namespace, appname)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	//render.Render(w, r, NewArticleResponse(article))

	render.Render(w, r, &App{Namespace: namespace, AppName: appname, Containers: containers})
}

// DeployImage - Modifications to make to deployment
type DeployImage struct {
	Namespace      string `json:"namespace"`
	AppName        string `json:"app_name"`
	ContainerName  string `json:"container_name"`
	Image          string `json:"image"`
}


// SetAnnotations - PUT /
func (env *Env) SetAnnotations(w http.ResponseWriter, r *http.Request) {
	// namespace := chi.URLParam(r, "namespace")
	// appname := chi.URLParam(r, "appname")

	response := make(map[string]string)
	response["message"] = "Set Annotations"
	render.JSON(w, r, response)
}

// PatchAnnotations - PATCH /
func (env *Env) PatchAnnotations(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "Patched Annotations"
	render.JSON(w, r, response)
}

// App = Application metadata
type App struct {
	Namespace  string            `json:"name_space"`
	AppName    string            `json:"app_name"`
	Containers map[string]string `json:"containers"`
}

// Render - Render App Struct for JSON consumption?  something like that
func (a *App) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	//rd.Elapsed = 10
	return nil
}


//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render - Render an Error Response Object
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest - Data Payload is incorrect
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrRender - Error Rendering response
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

// ErrServerUnable - Error, Server unable to complete the request, check server logs
func ErrServerUnable(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 501,
		StatusText:     "Server Unable to complete the request.",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound - Resources not found
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

