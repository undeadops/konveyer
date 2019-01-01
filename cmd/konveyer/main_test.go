package main

import (
	"github.com/undeadops/konveyer/pkg"
)

type mockDB struct{}

func (mdb *mockDB) GetDeployments() ([]*root.Deployment, error) {
	depl := make([]*root.Deployment, 0)
	depl = append(depl, &root.Deployment{"2018-12-08 12:32:32", "2018-12-08 13:34:12", "foobar", "default", ""})
	depl = append(depl, &root.Deployment{"2018-11-08 09:32:32", "2018-12-04 11:34:12", "spidey", "default", ""})
	return depl, nil
}

// func TestGetDeployments(t *testing.T) {
// 	req := http.
// }
