package repo

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// Repo - Git Repository Settings
type Repo struct {
	URL      string `json:"repo_url"`
	Path     string `json:"repo_path"`
	SSHKey   string `json:"repo_ssh_key"`
	Branch   string `json:"repo_branch"`
	SyncTime int    `json:"sync_time"`
	Clone    *git.Repository
	Mutex    *sync.RWMutex
}

// New - Create New Git Repo Object
func New(url string, options ...func(*Repo)) *Repo {

	r := Repo{
		Branch: "Master",
		Path:   "/code",
		SSHKey: "",
		URL:    url,
	}

	for _, option := range options {
		option(&r)
	}
	err := r.attachRepo()
	if err != nil {
		fmt.Println("Error Attaching Repo")
		panic(err)
	}
	return &r
}

// Sync - Timer for pulling Latest git changes
func (repo *Repo) Sync(stop, stopped chan struct{}) {
	defer close(stopped)

	ticker := time.NewTicker(time.Duration(repo.SyncTime) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("GIT: Starting Sync")
			repo.PullRepo()
		case <-stop:
			fmt.Println("tick: caller has told us to stop")
			return
		}
	}
}

// // MkRoot - Create Root direcotry for git clone
// func (repo *Repo) MkRoot() error {
// 	var err error
// 	if _, err = os.Stat(repo.Path); os.IsNotExist(err) {
// 		err = os.Mkdir(repo.Path, os.FileMode(0775))
// 	}
// 	return err
// }

// func getRepo(repoPath string) ({

// }

// AttachRepo - Create Clone of Git Repo
func (repo *Repo) attachRepo() error {
	defer repo.Mutex.Unlock()
	fmt.Println("Path is: ", repo.Path)
	repo.Mutex.Lock()
	if _, err := os.Stat(repo.Path); os.IsNotExist(err) {
		fmt.Println("Repo Base Directory")
		err = os.MkdirAll(repo.Path, os.FileMode(0775))
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Clone Repo")
		sshAuth, err := ssh.NewPublicKeysFromFile("git", repo.SSHKey, "")
		if err != nil {
			fmt.Println("Error: Problem Loading SSH Key")
		}
		gitOpts := git.CloneOptions{
			URL:  repo.URL,
			Auth: sshAuth,
			//ReferenceName: plumbing.ReferenceName(repo.Branch),
			//SingleBranch:  true,
		}

		r, err := git.PlainClone(repo.Path, false, &gitOpts)
		if err != nil {
			fmt.Println("Error: Unable to Clone repo, ", err)
		}
		w, err := r.Worktree()
		if err != nil {
			fmt.Println("Error: Unable to Clone repo, ", err)
		}
		cloneOpts := git.PullOptions{
			Auth: sshAuth,
			//ReferenceName: plumbing.ReferenceName(repo.Branch),
			//SingleBranch:  true,
		}
		err = w.Pull(&cloneOpts)
		if err != nil {
			fmt.Println("There was an Error with pull but it may be ok...", err)
		}
		repo.Clone = r
	} else {
		fmt.Println("Attempting to open Repo")
		r, err := git.PlainOpen(repo.Path)
		if err != nil {
			fmt.Println(err.Error())
			// Not sure if I need checking above...
		}
		repo.Clone = r
	}

	w, err := repo.Clone.Worktree()
	if err != nil {
		fmt.Println("Error: Unable to Clone repo, ", err)
	}
	sshAuth, _ := ssh.NewPublicKeysFromFile("git", repo.SSHKey, "")
	cloneOpts := git.PullOptions{
		Auth: sshAuth,
		//ReferenceName: plumbing.ReferenceName(repo.Branch),
		//SingleBranch:  true,
	}

	// Pull the latest updates from Git, before making changes.
	if err := w.Pull(&cloneOpts); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Printf("GIT: %s is already up-to-date.\n", repo.Path)
		} else {
			fmt.Printf("Error: %s\n", err)
		}
	}
	fmt.Println("Success! We can has git repo")
	return nil
}

// PullRepo - Pull Latest changes
func (repo *Repo) PullRepo() error {
	sshAuth, _ := ssh.NewPublicKeysFromFile("git", repo.SSHKey, "")

	cloneOpts := git.PullOptions{
		Auth: sshAuth,
		//ReferenceName: plumbing.ReferenceName(repo.Branch),
		//SingleBranch:  true,
	}
	defer repo.Mutex.Unlock()
	repo.Mutex.Lock()
	w, err := repo.Clone.Worktree()
	if err != nil {
		fmt.Println("Error: Unable to attach worktree, ", err)
	}

	// Pull the latest updates from Git, before making changes.
	if err := w.Pull(&cloneOpts); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Printf("GIT: %s is already up-to-date.\n", repo.Path)
		} else {
			fmt.Printf("Error: %s\n", err)
		}
	}
	h, err := repo.Clone.Head()
	if err != nil {
		fmt.Printf("Error capturing head...\n")
	}
	fmt.Printf("head looks like: %s", h)
	return err
}

// PushRepo - Push latest changes back up
func (repo *Repo) PushRepo() error {
	sshAuth, _ := ssh.NewPublicKeysFromFile("git", repo.SSHKey, "")

	pushOpts := git.PushOptions{
		Auth: sshAuth,
	}
	defer repo.Mutex.Unlock()
	repo.Mutex.Lock()

	w, err := repo.Clone.Worktree()
	if err != nil {
		return errors.New("Error: Unable to attach worktree")
	}

	commit, err := w.Commit("Commit Made By Konveyer", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Konveyer",
			Email: "konveyer@konveyer.sh",
			When:  time.Now(),
		},
	})

	if err != nil {
		return errors.New("Error: Unable to Commit Changed deployment image")
	}

	err = repo.Clone.Push(&pushOpts)
	if err != nil {
		return errors.New("Error: There was an error pushing Repo")
	}

	fmt.Printf("Commited: Commit ID: %s", commit)
	return nil
}
