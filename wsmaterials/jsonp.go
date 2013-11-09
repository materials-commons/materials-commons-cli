package wsmaterials

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
)

type JsonpResponseWriter struct {
	writer   http.ResponseWriter
	callback string
}

func (j *JsonpResponseWriter) Header() http.Header {
	return j.writer.Header()
}

func (j *JsonpResponseWriter) WriteHeader(status int) {
	j.writer.WriteHeader(status)
}

func (j *JsonpResponseWriter) Write(bytes []byte) (int, error) {
	if j.callback != "" {
		bytes = []byte(fmt.Sprintf("%s(%s)", j.callback, bytes))
	}
	return j.writer.Write(bytes)
}

func NewJsonpResponseWriter(httpWriter http.ResponseWriter, callback string) *JsonpResponseWriter {
	jsonpResponseWriter := new(JsonpResponseWriter)
	jsonpResponseWriter.writer = httpWriter
	jsonpResponseWriter.callback = callback
	return jsonpResponseWriter
}

func JsonpFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	callback := req.Request.FormValue("callback")
	jsonpResponseWriter := NewJsonpResponseWriter(resp.ResponseWriter, callback)
	resp.ResponseWriter = jsonpResponseWriter
	chain.ProcessFilter(req, resp)
}
