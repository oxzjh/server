package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var (
	ReturnError bool

	DefaultResponse = ResponseBytes([]byte{'{', '}'})
)

type IResponse interface {
	Write(http.ResponseWriter)
}

type ResponseBytes []byte

func (rb ResponseBytes) Write(w http.ResponseWriter) {
	w.Write(rb)
}

type ResponseString string

func (rs ResponseString) Write(w http.ResponseWriter) {
	w.Write([]byte(rs))
}

type ResponseMap map[string]any

func (rm ResponseMap) Write(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(rm)
}

type responseContent struct {
	contentType string
	content     []byte
}

func (rc *responseContent) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", rc.contentType)
	w.Write(rc.content)
}

func NewContent(contentType string, content []byte) IResponse {
	return &responseContent{contentType, content}
}

type responseStatus struct {
	status  int
	content string
}

func (rs *responseStatus) Write(w http.ResponseWriter) {
	w.WriteHeader(rs.status)
	w.Write([]byte(rs.content))
}

func NewStatus(status int, content string) IResponse {
	return &responseStatus{status, content}
}

type responseError struct {
	Code int    `json:"code"`
	Err  string `json:"err,omitempty"`
}

func (re *responseError) Write(w http.ResponseWriter) {
	if EnableLog {
		log.Println(re.Err)
	}
	if !ReturnError {
		re.Err = ""
	}
	json.NewEncoder(w).Encode(re)
}

func NewError(code int, err string) IResponse {
	return &responseError{code, err}
}

type responseJson struct {
	data any
}

func (rj *responseJson) Write(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(rj.data)
}

func NewJson(data any) IResponse {
	return &responseJson{data}
}

type responseFile struct {
	r    *http.Request
	file string
}

func (rf *responseFile) Write(w http.ResponseWriter) {
	http.ServeFile(w, rf.r, rf.file)
}

func NewFile(r *http.Request, file string) IResponse {
	return &responseFile{r, file}
}

type responsePipe struct {
	rc io.ReadCloser
}

func (rp *responsePipe) Write(w http.ResponseWriter) {
	io.Copy(w, rp.rc)
	rp.rc.Close()
}

func NewPipe(rc io.ReadCloser) IResponse {
	return &responsePipe{rc}
}

type responseCustom struct {
	callback func(http.ResponseWriter)
}

func (rc *responseCustom) Write(w http.ResponseWriter) {
	rc.callback(w)
}

func NewCuston(callback func(http.ResponseWriter)) IResponse {
	return &responseCustom{callback}
}
