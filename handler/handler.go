package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pubgo/grpcox/internal/services/grpcproxy"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Handler hold all handler methods
type Handler struct {
	g *grpcproxy.serviceImpl
}

// InitHandler Constructor
func InitHandler() *Handler {
	return &Handler{
		g: grpcproxy.New(),
	}
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	body := new(bytes.Buffer)
	err := indexHTML.Execute(body, make(map[string]string))
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(body.Bytes())
}

func (h *Handler) getActiveConns(w http.ResponseWriter, r *http.Request) {
	response(w, h.g.GetActiveConns(context.TODO()))
}

func (h *Handler) closeActiveConns(w http.ResponseWriter, r *http.Request) {
	host := chi.URLParam(r, "host")
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Host"))
		return
	}

	err := h.g.CloseActiveConns(strings.Trim(host, " "))
	if err != nil {
		writeError(w, err)
		return
	}
	response(w, map[string]bool{"success": true})
}

func (h *Handler) getLists(w http.ResponseWriter, r *http.Request) {
	host := chi.URLParam(r, "host")
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Host"))
		return
	}

	service := chi.URLParam(r, "serv_name")

	useTLS, _ := strconv.ParseBool(r.Header.Get("use_tls"))
	restart, _ := strconv.ParseBool(r.FormValue("restart"))

	res, err := h.g.GetResource(context.Background(), host, !useTLS, restart)
	if err != nil {
		writeError(w, err)
		return
	}

	result, err := res.List(service)
	if err != nil {
		writeError(w, err)
		return
	}

	h.g.Extend(host)
	response(w, result)
}

// getListsWithProto handling client request for service list with proto
func (h *Handler) getListsWithProto(w http.ResponseWriter, r *http.Request) {
	host := chi.URLParam(r, "host")
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Host"))
		return
	}

	service := chi.URLParam(r, "serv_name")

	useTLS, _ := strconv.ParseBool(r.Header.Get("use_tls"))
	restart, _ := strconv.ParseBool(r.FormValue("restart"))

	// limit upload file to 5mb
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		writeError(w, err)
		return
	}

	// convert uploaded files to list of Proto struct
	files := r.MultipartForm.File["protos"]
	protos := make([]grpcproxy.Proto, 0, len(files))
	for _, file := range files {
		fileData, err := file.Open()
		if err != nil {
			writeError(w, err)
			return
		}
		defer fileData.Close()

		content, err := ioutil.ReadAll(fileData)
		if err != nil {
			writeError(w, err)
		}

		protos = append(protos, grpcproxy.Proto{
			Name:    file.Filename,
			Content: content,
		})
	}

	res, err := h.g.GetResourceWithProto(context.Background(), host, !useTLS, restart, protos)
	if err != nil {
		writeError(w, err)
		return
	}

	result, err := res.List(service)
	if err != nil {
		writeError(w, err)
		return
	}

	h.g.Extend(host)
	response(w, result)
}

func (h *Handler) describeFunction(w http.ResponseWriter, r *http.Request) {
	host := chi.URLParam(r, "host")
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Host"))
		return
	}

	funcName := chi.URLParam(r, "func_name")
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Func Name"))
		return
	}

	useTLS, _ := strconv.ParseBool(r.Header.Get("use_tls"))

	res, err := h.g.GetResource(context.Background(), host, !useTLS, false)
	if err != nil {
		writeError(w, err)
		return
	}

	// get param
	result, _, err := res.Describe(funcName)
	if err != nil {
		writeError(w, err)
		return
	}
	match := reGetFuncArg.FindStringSubmatch(result)
	if len(match) < 2 {
		writeError(w, fmt.Errorf("Invalid Func Type"))
		return
	}

	// describe func
	result, template, err := res.Describe(match[1])
	if err != nil {
		writeError(w, err)
		return
	}

	type desc struct {
		Schema   string `json:"schema"`
		Template string `json:"template"`
	}

	h.g.Extend(host)
	response(w, desc{
		Schema:   result,
		Template: template,
	})

}

func (h *Handler) invokeFunction(w http.ResponseWriter, r *http.Request) {
	host := chi.URLParam(r, "host")
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Host"))
		return
	}

	funcName := chi.URLParam(r, "func_name")
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Func Name"))
		return
	}

	useTLS, _ := strconv.ParseBool(r.Header.Get("use_tls"))

	res, err := h.g.GetResource(context.Background(), host, !useTLS, false)
	if err != nil {
		writeError(w, err)
		return
	}

	// context metadata
	metadataHeader := r.Header.Get("Metadata")
	metadataArr := strings.Split(metadataHeader, ",")

	// construct array of string with "key: value" form to satisfy grpcurl MetadataFromHeaders
	var metadata []string
	var metadataStr string
	for i, m := range metadataArr {
		i += 1
		if isEven := i%2 == 0; isEven {
			metadataStr = metadataStr + m
			metadata = append(metadata, metadataStr)
			metadataStr = ""
			continue
		}
		metadataStr = fmt.Sprintf("%s:", m)
	}

	// get param
	result, timer, err := res.Invoke(context.Background(), metadata, funcName, r.Body)
	if err != nil {
		writeError(w, err)
		return
	}

	type invRes struct {
		Time   string `json:"timer"`
		Result string `json:"result"`
	}

	h.g.Extend(host)
	response(w, invRes{
		Time:   timer.String(),
		Result: result,
	})
}

func (h *Handler) listRequest(w http.ResponseWriter, r *http.Request) {
	var dataList []map[string]interface{}
	if err := list(db, bucketName, func(data *map[string]interface{}) {
		dataList = append(dataList, *data)
	}); err != nil {
		writeError(w, err)
		return
	}
	response(w, dataList)
}

func (h *Handler) saveRequest(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, err)
		return
	}

	id := genID()
	data["id"] = id
	if err := set(db, bucketName, id, data); err != nil {
		writeError(w, err)
		return
	}

	response(w, data)
}

func (h *Handler) delRequest(w http.ResponseWriter, r *http.Request) {
	var name = chi.URLParam(r, "name")
	if err := del(db, bucketName, name); err != nil {
		writeError(w, err)
		return
	}

	response(w, nil)
}

func (h *Handler) updateRequest(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, err)
		return
	}

	var name = chi.URLParam(r, "name")
	if err := set(db, bucketName, name, data); err != nil {
		writeError(w, err)
		return
	}

	response(w, nil)
}

func (h *Handler) getRequest(w http.ResponseWriter, r *http.Request) {
	var name = chi.URLParam(r, "name")
	var data map[string]interface{}
	if err := get(db, bucketName, name, &data); err != nil {
		writeError(w, err)
		return
	}

	response(w, data)
}

func (h *Handler) delAllRequest(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := db.Update(func(tx *bolt.Tx) error { return tx.DeleteBucket([]byte(bucketName)) }); err != nil {
		writeError(w, err)
		return
	}

	response(w, data)
}

func (h *Handler) downloadAllRequest(w http.ResponseWriter, r *http.Request) {
	var dataList []map[string]interface{}
	if err := list(db, bucketName, func(data *map[string]interface{}) {
		dataList = append(dataList, *data)
	}); err != nil {
		writeError(w, err)
		return
	}
	responseFile(w, dataList)
}

func genID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
