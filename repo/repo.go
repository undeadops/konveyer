package repo

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
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
	Logger   *logrus.Logger
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
		r.Logger.WithFields(logrus.Fields{"git_path": r.Path}).Error("Error Attaching to Repo")
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
			repo.Logger.WithFields(
				logrus.Fields{"git_path": repo.Path},
			).Info("Starting Git Sync")
			repo.PullRepo()
		case <-stop:
			repo.Logger.Info("Stoping Sync")
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
	repo.Logger.WithFields(
		logrus.Fields{"git_path": repo.Path},
	).Info("Attaching/Pulling Git Repo")
	repo.Mutex.Lock()

	// change to nuke existing, always starting with external git repo as source of truth
	
	if _, err := os.Stat(repo.Path); os.IsNotExist(err) {
		err = os.MkdirAll(repo.Path, os.FileMode(0775))
		if err != nil {
			return errors.New("Error: Unable to Create Code Path, " + err.Error())
		}
		sshAuth, err := ssh.NewPublicKeysFromFile("git", repo.SSHKey, "")
		if err != nil {
			return errors.New("Error: Problem Loading SSH Key, " + err.Error())
		}
		gitOpts := git.CloneOptions{
			URL:  repo.URL,
			Auth: sshAuth,
			//ReferenceName: plumbing.ReferenceName(repo.Branch),
			//SingleBranch:  true,
		}

		r, err := git.PlainClone(repo.Path, false, &gitOpts)
		if err != nil {
			repo.Logger.WithFields(
				logrus.Fields{"git_path": repo.Path},
			).Error("Problem Cloning Repo")
		}
		w, err := r.Worktree()
		if err != nil {
			repo.Logger.WithFields(
				logrus.Fields{"git_path": repo.Path},
			).Error("Unable to set worktree on repo")
		}
		cloneOpts := git.PullOptions{
			Auth: sshAuth,
			//ReferenceName: plumbing.ReferenceName(repo.Branch),
			//SingleBranch:  true,
		}
		err = w.Pull(&cloneOpts)
		if err != nil {
			repo.Logger.WithFields(
				logrus.Fields{"git_path": repo.Path, "error": err},
			).Error("There was an error Pulling Repo, it may be ok?")
		}
		repo.Clone = r
	} else {
		repo.Logger.WithFields(
			logrus.Fields{"git_path": repo.Path},
		).Info("Attaching to existing Repo")
		r, err := git.PlainOpen(repo.Path)
		if err != nil {
			return errors.New("Unable to attach to existing git repo, " + err.Error())
		}
		repo.Clone = r
	}

	w, err := repo.Clone.Worktree()
	if err != nil {
		return errors.New("Unable to clone worktree, " + err.Error())
	}
	sshAuth, _ := ssh.NewPublicKeysFromFile("git", repo.SSHKey, "")
	cloneOpts := git.PullOptions{
		Auth: sshAuth,
		//ReferenceName: plumbing.ReferenceName(repo.Branch),
		//SingleBranch:  true,
		Force: true,
	}

	// Pull the latest updates from Git, before making changes.
	if err := w.Pull(&cloneOpts); err != nil {
		h, _ := repo.Clone.Head()
		if err == git.NoErrAlreadyUpToDate {
			repo.Logger.WithFields(
				logrus.Fields{"git_path": repo.Path, "git_hash": h},
			).Info("Repo already up-to date!")
		} else {
			return errors.New("Unable to pull repo update")
		}
	}
	h, _ := repo.Clone.Head()
	repo.Logger.WithFields(
		logrus.Fields{"git_path": repo.Path, "git_hash": h},
	).Info("Sucess!")
	return nil
}

// PullRepo - Pull Latest changes
func (repo *Repo) PullRepo() error {
	sshAuth, _ := ssh.NewPublicKeysFromFile("git", repo.SSHKey, "")

	cloneOpts := git.PullOptions{
		Auth: sshAuth,
		//ReferenceName: plumbing.ReferenceName(repo.Branch),
		//SingleBranch:  true,
		Force: true,
	}
	defer repo.Mutex.Unlock()
	repo.Mutex.Lock()
	w, err := repo.Clone.Worktree()
	if err != nil {
		return errors.New("Unable to clone worktree, " + err.Error())
	}

	// Pull the latest updates from Git, before making changes.
	if err := w.Pull(&cloneOpts); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			repo.Logger.WithFields(
				logrus.Fields{"git_path": repo.Path},
			).Info("Repo already up-to date!")
		} else {
			return errors.New("Unable to pull repo update")
		}
	}
	h, _ := repo.Clone.Head()
	repo.Logger.WithFields(
		logrus.Fields{"git_path": repo.Path, "git_hash": h},
	).Info("Latest Git Pull")
	return err
}

// PushRepo - Push latest changes back up
func (repo *Repo) PushRepo() error {
	sshAuth, _ := ssh.NewPublicKeysFromFile("git", repo.SSHKey, "")

	pushOpts := git.PushOptions{
		Auth: sshAuth,
	}
	// Assumption is being made that whatever calls PushRepo already has a mutex lock
	// Maybe context passing is the correct way for passing around the lock?
	// defer repo.Mutex.Unlock()
	// repo.Mutex.Lock()

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

	repo.Logger.WithFields(
		logrus.Fields{"git_path": repo.Path, "git_hash": commit},
	).Info("Commited")
	return nil
}
