package config

import (
	"github.com/pkg/errors"
)

// Options ...
type Options struct {
	Server   string
	Database string
}

// Validate ...
// Check CLI arguments and encapsulate them into an Options if ok.
func Validate(server, db string) (Options, error) {
	if server == "" || db == "" {
		return Options{}, errors.New("server and db arguments are required")
	}
	return Options{
		Server:   server,
		Database: db,
	}, nil
}
