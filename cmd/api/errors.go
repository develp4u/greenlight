package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	// 응답 본문에 에러코드와 에러메세지를 작성한다.
	err := app.writeJson(w, status, env, nil)
	if err != nil {
		app.logError(r, err) // error log
		w.WriteHeader(500)   // error http status code
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	// 서버에 문제가 발생하여 요청을 처리할 수 없습니다.
	message := "The server encountered a problem and could not process your request."
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	// 요청된 리소스를 찾을 수 없습니다.
	message := "The requested resource could not be found."
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {

	// 이 리소스에 대해 POST 메서드는 지원되지 않습니다.
	message := fmt.Sprintf("The %s method is not supported for this resource.", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
