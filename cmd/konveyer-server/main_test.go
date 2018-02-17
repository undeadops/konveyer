package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/undeadops/konveyer"
)

func TestMain(m *testing.M) {
	r := konveyer.Initialize()
	code := m.Run()
	fmt.Printf(r.MongoUri)
	//konveyer.Shutdown(r)
	os.Exit(code)
}

func TestIndexHandler(t *testing.T) {
	t.Parallel()

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Index(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	log.Printf("Looking at response...: %s", w.Body.Bytes())
}
