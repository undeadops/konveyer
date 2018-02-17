package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/undeadops/konveyer"
)

var (
	runtime *konveyer.Runtime
)

func init() {
	runtime = konveyer.Initialize()
}

func TestGetProjectsHandler(t *testing.T) {
	t.Parallel()

	r, err := http.NewRequest("GET", "/project", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	GetProjects(runtime, w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	fmt.Printf("Looking at response...: %s", w.Body.Bytes())
}

// func TestCreateProjectsHandler(t *testing.T) {
// 	t.Parallel()

// 	r, err := http.NewRequest("PUT", "/project", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	w := httptest.NewRecorder()

// 	CreateProjects(runtime, w, r)

// 	assert.Equal(t, http.StatusCreated, w.Code)
// 	fmt.Printf("Looking at response....: %s", w.Body.Bytes())
// }
