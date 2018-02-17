package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/undeadops/konveyer"
	"github.com/undeadops/konveyer/models"
)

// GetProjects - GET /projects
func GetProjects(k *konveyer.Runtime, w http.ResponseWriter, r *http.Request) error {
	message, err := models.ListProjects(k)
	if err != nil {
		return StatusError{http.StatusInternalServerError, fmt.Errorf("Error Creating List of Projects")}
	}
	respondWithJSON(w, http.StatusOK, message)
	return nil
}

// CreateProjects - POST /projects
func CreateProjects(k *konveyer.Runtime, w http.ResponseWriter, r *http.Request) error {
	var p models.Project
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return err
	}
	defer r.Body.Close()

	if err := models.CreateProject(k, p); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return err
	}
	respondWithJSON(w, http.StatusCreated, "Created")
	return nil
}

// DescribeProject - GET /project/{projectName}
func DescribeProject(k *konveyer.Runtime, w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	log.Printf("Variables: %s\n", params)
	if params["projectName"] == "" {
		return StatusError{http.StatusNotFound, fmt.Errorf("Not Found")}
	}

	message, err := models.DescribeProject(k, params["projectName"])
	if err != nil {
		return StatusError{http.StatusInternalServerError, fmt.Errorf("Invalid Project")}
	}
	respondWithJSON(w, http.StatusOK, message)
	return nil
}
