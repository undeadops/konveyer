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

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"github.com/go-chi/valve"

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
	var syncSecs = flag.Int("sync_seconds", 300, "Time between Syncs")
	var gitWebhook = flag.Bool("use-webhook", false, "Use Webhook instead of timed Sync")

	fmt.Println("Staring app...")
	//Gather environment variables that start with CI_
	//var cluster Cluster

	flag.Parse()

	fmt.Println("Initializing Git Repo")

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

	repo := repo.New(*repoFlag, sshKey, gitPath, syncTime, mutex)

	stopChan := make(chan struct{})
	stoppedChan := make(chan struct{})

	env := &Env{repo: repo, repoPath: *pathFlag, stop: stopChan, stopped: stoppedChan}

	// Our graceful valve shut-off package to manage code preemption and
	// shutdown signaling.
	valv := valve.New()
	baseCtx := valv.Context()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// HTTP service running in this program as well. The valve context is set
	// as a base context on the server listener at the point where we instantiate
	// the server - look lower.
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

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
			fmt.Println("shutting down..")

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
				fmt.Println("not all connections done")
			case <-ctx.Done():
			}

			// shutdown git
			fmt.Println("Closing out Git")
			close(stopChan)
			<-stoppedChan
			fmt.Println("Git Closed")
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
		r.Put("/image/{imageId}", env.SetDeployApp)   // PUT /deployment/namespace/appname/image/master-1234abf
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

// SetDeployApp - PUT /v1/deployment/<namespace>/<appname>
func (env *Env) SetDeployApp(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	appname := chi.URLParam(r, "appname")

	render.Render(w, r, &App{Namespace: namespace, AppName: appname})
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
