// internal/db/errors.go
package db

import "errors"

var ErrNotInitialized = errors.New("db pool not initialized")
