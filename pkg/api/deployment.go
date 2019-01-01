package api

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/undeadops/konveyer/pkg"
)

type deploymentRouter struct {
	deploymentService root.DeploymentService
	//auth *authHelper
}

// NewDeploymentRouter - Create Routes for dealing with Deployment Manifests
func NewDeploymentRouter(u root.DeploymentService) chi.Router {
	deploymentRouter := deploymentRouter{u}

	r := chi.NewRouter()
	r.Route("/{namespace}/{appname}", func(r chi.Router) {
		r.Get("/", deploymentRouter.getDeployApp) // GET /v1/deployment/namespace/appname/
		//r.Post("/", deploymentRouter.CreateDeployApp)              // POST /v1/deployment/namespace/appname/
		//r.Put("/image", deploymentRouter.SetDeployImage)           // PUT /v1/deployment/namespace/appname/image -d '{ "container_name": "konveyer", "image": "undeadops/konveyer:master" }
		//r.Put("/annotations", deploymentRouter.SetAnnotations)     // PUT /v1/deployment/namespace/appname -d '{ []map[string]string }'
		//r.Patch("/annotations", deploymentRouter.PatchAnnotations) // PATCH /v1/deployment/namespace/app -d '{ []map[string]string }'
	})
	return r
}

// // DeploymentRouter - Handle Deployment manifests in Kubernetes
// func (a *API) DeploymentRouter() chi.Router {
// 	r := chi.NewRouter()
// 	r.Route("/{namespace}/{appname}", func(r chi.Router) {
// 		r.Get("/", a.GetDeployApp)                  // GET /v1/deployment/namespace/appname/
// 		r.Post("/", a.CreateDeployApp)              // POST /v1/deployment/namespace/appname/
// 		r.Put("/image", a.SetDeployImage)           // PUT /v1/deployment/namespace/appname/image -d '{ "container_name": "konveyer", "image": "undeadops/konveyer:master" }
// 		r.Put("/annotations", a.SetAnnotations)     // PUT /v1/deployment/namespace/appname -d '{ []map[string]string }'
// 		r.Patch("/annotations", a.PatchAnnotations) // PATCH /v1/deployment/namespace/app -d '{ []map[string]string }'
// 	})
// 	return r
// }

// GetDeployApp - Display a currently deployed app
func (d *deploymentRouter) getDeployApp(w http.ResponseWriter, r *http.Request) {

	namespace := chi.URLParam(r, "namespace")
	appname := chi.URLParam(r, "appname")

	//render.Render(w, r, NewArticleResponse(article))
	err, deployment := d.deploymentService.GetDeployment(appname, namespace)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Json(w, http.StatusOK, deployment)
}

// func(ur *userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
//   vars := mux.Vars(r)
//   username := vars["username"]

//   err, user := ur.userService.GetUserByUsername(username)
//   if err != nil {
//     Error(w, http.StatusNotFound, err.Error())
//     return
//   }

//   Json(w, http.StatusOK, user)
// }

// func(ur* userRouter) createUserHandler(w http.ResponseWriter, r *http.Request) {
//   err, user := decodeUser(r)
//   if err != nil {
//     Error(w, http.StatusBadRequest, "Invalid request payload")
//     return
//   }

//   err = ur.userService.CreateUser(&user)
//   if err != nil {
//     Error(w, http.StatusInternalServerError, err.Error())
//     return
//   }

//   Json(w, http.StatusOK, err)
// }

// // CreateDeployApp - Create a new App Deployment
// func (a *API) CreateDeployApp(w http.ResponseWriter, r *http.Request) {
// 	namespace := chi.URLParam(r, "namespace")
// 	appname := chi.URLParam(r, "appname")

// 	var d db.Deployment

// 	d.Namespace = namespace
// 	d.App = appname
// 	d.Image = "undeadops/webby:latest"

// 	// if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
// 	// 	render.Render(w, r, ErrInvalidRequest(err))
// 	// 	return
// 	// }

// 	if err := a.db.CreateDeployment(d); err != nil {
// 		render.Render(w, r, ErrServerUnable(err))
// 		return
// 	}

// 	sendJSON(w, http.StatusCreated, map[string]string{"message": "successfully created"})
// }

// // SetDeployImage - Update Applications Deployed Image
// func (a *API) SetDeployImage(w http.ResponseWriter, r *http.Request) {
// 	namespace := chi.URLParam(r, "namespace")
// 	appname := chi.URLParam(r, "appname")

// 	var d DeployImage

// 	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
// 		render.Render(w, r, ErrInvalidRequest(err))
// 		return
// 	}

// 	//deployment, err := env.repo.GetDeploymentImage(namespace, appname)
// 	//if err != nil {
// 	//	render.Render(w, r, ErrInvalidRequest(err))
// 	//	return
// 	//}

// 	// var msg map[string]interface {}
// 	// msg["old_image"] = deployment[d.ContainerName]
// 	// msg["new_image"] = d.Image

// 	//LogEntrySetFields(r, msg)
// 	// TODO: Handle multiple containers somehow...
// 	// err := env.repo.SetDeploymentImage(namespace, appname, d.Image)
// 	// if err != nil {
// 	// 	render.Render(w, r, ErrServerUnable(err))
// 	// 	return
// 	// }
// 	// containers, err := env.repo.GetDeploymentImage(namespace, appname)
// 	// if err != nil {
// 	// 	render.Render(w, r, ErrInvalidRequest(err))
// 	// 	return
// 	// }

// 	render.Status(r, http.StatusCreated)
// 	//render.Render(w, r, NewArticleResponse(article))

// 	containers := make(map[string]string)
// 	containers["bar"] = "patched something"
// 	render.Render(w, r, &App{Namespace: namespace, AppName: appname, Containers: containers})
// }

// func (a *API) SetAnnotations(w http.ResponseWriter, r *http.Request) {

// }

// func (a *API) PatchAnnotations(w http.ResponseWriter, r *http.Request) {

// }

// // DeployImage - Modifications to make to deployment
// type DeployImage struct {
// 	Namespace     string `json:"namespace"`
// 	AppName       string `json:"app_name"`
// 	ContainerName string `json:"container_name"`
// 	Image         string `json:"image"`
// }

// // App = Application metadata
// type App struct {
// 	Namespace  string            `json:"name_space"`
// 	AppName    string            `json:"app_name"`
// 	Containers map[string]string `json:"containers"`
// }

// // Render - Render App Struct for JSON consumption?  something like that
// func (a *App) Render(w http.ResponseWriter, r *http.Request) error {
// 	// Pre-processing before a response is marshalled and sent across the wire
// 	//rd.Elapsed = 10
// 	return nil
// }
