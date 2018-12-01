package api


type Server struct {
	repo     *repo.Repo
	repoPath string
	stop     chan struct{}
	stopped  chan struct{}
}

