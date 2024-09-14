package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) readIDFrom(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	//b, err := json.Marshal(data)
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	defer r.Body.Close()

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			app.logger.Info("syntaxError!!")
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			app.logger.Info("ErrUnexpectedEOF!!")
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			app.logger.Info("unmarshalTypeError!!")
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrent JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrent JSON (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			app.logger.Info("io.EOF!!")
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			app.logger.Info("Unknown field!!")
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			app.logger.Info("maxBytesError!!")
			return fmt.Errorf("body must not be larger than %s bytes", humanize.Comma(maxBytesError.Limit))

		case errors.As(err, &invalidUnmarshalError):
			app.logger.Info("invalidUnmarshalError!!")
			panic(err)
		default:
			app.logger.Info("default!!")
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		app.logger.Info("io.EOF!!")
		return errors.New("body must only contains a single JSON value")
	}

	return nil
}
