package mongo

import (
	"net/url"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/undeadops/konveyer/pkg"
)

// Session - MongoDB session struct
type Session struct {
	session *mgo.Session
}

// NewSession - Return MongoDB Session
func NewSession(config *root.Config) (*Session, error) {
	m, _ := url.Parse(config.MongoURI)

	AuthUserName := m.User.Username()
	AuthPassword, _ := m.User.Password()
	Hosts := m.Host
	AuthDatabase := strings.TrimPrefix(m.Path, "/")

	di := &mgo.DialInfo{
		Addrs:    []string{Hosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}
	session, err := mgo.DialWithInfo(di)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return &Session{session}, err
}

func (s *Session) Copy() *mgo.Session {
	return s.session.Copy()
}

func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

func (s *Session) DropDatabase(db string) error {
	if s.session != nil {
		return s.session.DB(db).DropDatabase()
	}
	return nil
}
