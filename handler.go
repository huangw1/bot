package bot

import (
	"github.com/pkg/errors"
	"regexp"
)

type RespondFunc func(Message) error

type responseHandler struct {
	regexp *regexp.Regexp
	run    RespondFunc
}

func newHandler(expr string, fun RespondFunc) (responseHandler, error) {
	h := responseHandler{run: fun}
	var err error
	h.regexp, err = regexp.Compile(expr)
	if err != nil {
		return h, errors.Wrap(err, "invalid regular expression")
	}
	return h, nil
}
