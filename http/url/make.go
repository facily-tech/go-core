package url

import (
	"net/url"
	"path"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

func Make(basePath, requestPath string, queryParameters interface{}) (string, error) {
	u, err := url.Parse(basePath)
	if err != nil {
		return "", errors.Wrap(err, "fail parsing requestPath")
	}

	u.Path = path.Join(u.Path, requestPath)
	v, err := query.Values(queryParameters)
	if err != nil {
		return "", err
	}

	u.RawQuery = v.Encode()
	return u.String(), nil
}
