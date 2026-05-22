package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Interceptor func(b []byte, w http.ResponseWriter)

var (
	ReturnError bool

	DefaultResponse = ResponseBytes([]byte{'{', '}'})
)

type IResponse interface {
	Write(http.ResponseWriter, Interceptor)
}

type ResponseBytes []byte

func (r ResponseBytes) Write(w http.ResponseWriter, interceptor Interceptor) {
	if interceptor != nil {
		interceptor(r, w)
	} else {
		w.Write(r)
	}
}

type ResponseString string

func (r ResponseString) Write(w http.ResponseWriter, interceptor Interceptor) {
	if interceptor != nil {
		interceptor([]byte(r), w)
	} else {
		w.Write([]byte(r))
	}
}

type ResponseMap map[string]any

func (r ResponseMap) Write(w http.ResponseWriter, interceptor Interceptor) {
	if interceptor != nil {
		b, _ := json.Marshal(r)
		interceptor(b, w)
	} else {
		json.NewEncoder(w).Encode(r)
	}
}

type responseContent struct {
	contentType string
	content     []byte
}

func (r *responseContent) Write(w http.ResponseWriter, interceptor Interceptor) {
	w.Header().Set("Content-Type", r.contentType)
	if interceptor != nil {
		interceptor(r.content, w)
	} else {
		w.Write(r.content)
	}
}

func NewContent(contentType string, content []byte) IResponse {
	return &responseContent{contentType, content}
}

type responseStatus struct {
	status  int
	content string
}

func (r *responseStatus) Write(w http.ResponseWriter, _ Interceptor) {
	w.WriteHeader(r.status)
	w.Write([]byte(r.content))
}

func NewStatus(status int, content string) IResponse {
	return &responseStatus{status, content}
}

type responseError struct {
	Code int    `json:"code"`
	Err  string `json:"err,omitempty"`
}

func (r *responseError) Write(w http.ResponseWriter, interceptor Interceptor) {
	if EnableLog {
		log.Println(r.Err)
	}
	if !ReturnError {
		r.Err = ""
	}
	if interceptor != nil {
		b, _ := json.Marshal(r)
		interceptor(b, w)
	} else {
		json.NewEncoder(w).Encode(r)
	}
}

func NewError(code int, err string) IResponse {
	return &responseError{code, err}
}

type responseJson struct {
	data any
}

func (r *responseJson) Write(w http.ResponseWriter, interceptor Interceptor) {
	if interceptor != nil {
		b, _ := json.Marshal(r.data)
		interceptor(b, w)
	} else {
		json.NewEncoder(w).Encode(r.data)
	}
}

func NewJson(data any) IResponse {
	return &responseJson{data}
}

type responseFile struct {
	req  *http.Request
	file string
}

func (r *responseFile) Write(w http.ResponseWriter, _ Interceptor) {
	http.ServeFile(w, r.req, r.file)
}

func NewFile(r *http.Request, file string) IResponse {
	return &responseFile{r, file}
}

type responsePipe struct {
	rc io.ReadCloser
}

func (r *responsePipe) Write(w http.ResponseWriter, _ Interceptor) {
	io.Copy(w, r.rc)
	r.rc.Close()
}

func NewPipe(rc io.ReadCloser) IResponse {
	return &responsePipe{rc}
}

type responseCustom struct {
	callback func(http.ResponseWriter)
}

func (r *responseCustom) Write(w http.ResponseWriter, _ Interceptor) {
	r.callback(w)
}

func NewCuston(callback func(http.ResponseWriter)) IResponse {
	return &responseCustom{callback}
}
