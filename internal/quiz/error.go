package quiz

import "errors"

// errShort means the hit could not be turned into a fair question.
var errShort = errors.New("snippet too thin")
