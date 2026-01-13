package domain

import "errors"

// ErrNotImplemented は未実装エラーです.
var ErrNotImplemented = errors.New("not implemented: need repository owner information")
