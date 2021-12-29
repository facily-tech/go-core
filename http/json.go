package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/facily-tech/go-core/log"
	"github.com/pkg/errors"
)

// ToJSON transform a struct given into json string.
func ToJSON(ctx context.Context, logger log.Logger, input interface{}, w http.ResponseWriter, status int) {
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(ctx, "can't encode input to json object", log.Any("input", input))
	}
	w.WriteHeader(status)
	if n, err := fmt.Fprintf(w, "%s", buffer.String()); err != nil {
		logger.Error(ctx, "can't write json to writer",
			log.Any("bytes wrote", n),
			log.Error(err),
		)
	}
}

// FromJSON turns a json string into a struct by output given.
func FromJSON(logger log.Logger, output interface{}, r io.Reader) error {
	content, err := io.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "unable to read body to decode")
	}
	if err := json.Unmarshal(content, output); err != nil {
		return errors.Wrapf(err, "unable to decode body: %s", content)
	}

	return nil
}
