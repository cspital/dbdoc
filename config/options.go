package config

import (
	"github.com/pkg/errors"
)

// Credentials ...
type Credentials struct {
	Username string
	Password string
}

// Options ...
type Options struct {
	Server      string
	Database    string
	Credentials *Credentials
}

// Validate ...
// Check CLI arguments and encapsulate them into an Options if ok.
func Validate(server, db, user, pass string) (Options, error) {
	if server == "" || db == "" {
		return Options{}, errors.New("server and db arguments are required")
	}
	opts := Options{
		Server:   server,
		Database: db,
	}

	if user != "" {
		opts.Credentials = &Credentials{
			Username: user,
			Password: pass,
		}
	}

	return opts, nil
}
