package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"
)

var (
	reGetFuncArg *regexp.Regexp
	indexHTML    *template.Template
)

func init() {
	reGetFuncArg = regexp.MustCompile("\\( (.*) \\) returns")
	indexHTML = template.Must(template.New("index.html").Delims("{[", "]}").ParseFiles("index/index.html"))
}

// Response - Standar ajax Response
type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data"`
}

func writeError(w http.ResponseWriter, err error) {
	e, _ := json.Marshal(Response{
		Error: err.Error(),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(e)
}

func response(w http.ResponseWriter, data interface{}) {
	e, _ := json.Marshal(Response{
		Data: data,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(e)
}

func responseFile(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/octet-stream")
	//强制浏览器下载
	w.Header().Set("Content-Disposition", "attachment; filename=data.json")
	//浏览器下载或预览
	w.Header().Set("Content-Disposition", "inline;filename=data.json")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Cache-Control", "no-cache")

	w.WriteHeader(http.StatusOK)
	e, _ := json.Marshal(Response{
		Data: data,
	})
	w.Write(e)
}
