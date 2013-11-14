package wsmaterials

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
)

type jsonpResponseWriter struct {
	writer   http.ResponseWriter
	callback string
}

func (j *jsonpResponseWriter) Header() http.Header {
	return j.writer.Header()
}

func (j *jsonpResponseWriter) WriteHeader(status int) {
	j.writer.WriteHeader(status)
}

func (j *jsonpResponseWriter) Write(bytes []byte) (int, error) {
	if j.callback != "" {
		bytes = []byte(fmt.Sprintf("%s(%s)", j.callback, bytes))
	}
	return j.writer.Write(bytes)
}

func newJsonpResponseWriter(httpWriter http.ResponseWriter, callback string) *jsonpResponseWriter {
	jsonpResponseWriter := new(jsonpResponseWriter)
	jsonpResponseWriter.writer = httpWriter
	jsonpResponseWriter.callback = callback
	return jsonpResponseWriter
}

func JsonpFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	callback := req.Request.FormValue("callback")
	jsonpResponseWriter := newJsonpResponseWriter(resp.ResponseWriter, callback)
	resp.ResponseWriter = jsonpResponseWriter
	chain.ProcessFilter(req, resp)
}
