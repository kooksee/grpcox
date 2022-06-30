package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pubgo/xerror"
	"honnef.co/go/js/dom/v2"
	fetch "marwan.io/wasm-fetch"

	"github.com/gusaul/grpcox/web/jsutil"
)

func (h *Home) getValidTarget() string {
	//t = $('#server-target').val().trim();
	//    if (t == "") {
	//        return target;
	//    }
	//
	//    ts = t.split("://");
	//    if (ts.length > 1) {
	//        $('#server-target').val(ts[1]);
	//        return ts[1];
	//    }
	//    return ts[0];
	return ""
}

//

func (h *Home) serviceClose(addr string) bool {
	rsp, err := fetch.Fetch(`active/close/`+addr, &fetch.Opts{
		Method: fetch.MethodDelete,
	})
	xerror.Panic(err)
	var data = new(struct {
		Data struct {
			Success bool `json:"success"`
		} `json:"data"`
	})
	xerror.Panic(json.Unmarshal(rsp.Body, data))
	return data.Data.Success
}

func (h *Home) listActiveSrv() []string {
	rsp, err := fetch.Fetch(`active/get`, &fetch.Opts{
		Method: fetch.MethodGet,
	})
	xerror.Panic(err)

	var data = new(struct {
		Data []string `json:"data"`
	})
	xerror.Panic(json.Unmarshal(rsp.Body, data))
	return data.Data
}

func (h *Home) listServices(target string) []string {
	rsp, err := fetch.Fetch(fmt.Sprintf(`/server/%s/services?restart=1`, target), &fetch.Opts{
		Method: fetch.MethodGet,
	})
	xerror.Panic(err)
	var data = new(struct {
		Data []string `json:"data"`
	})
	xerror.Panic(json.Unmarshal(rsp.Body, data))
	return data.Data
}

type invokeFuncRsp struct {
	Data struct {
		Timer  string `json:"timer"`
		Result string `json:"result"`
	} `json:"data"`
}

func (h *Home) invokeFunc(target, fn string, body string) *invokeFuncRsp {
	rsp, err := fetch.Fetch("/server/"+target+"/function/"+fn+"/invoke", &fetch.Opts{
		Method:  fetch.MethodPost,
		Body:    strings.NewReader(body),
		Headers: map[string]string{
			//"ContentType": "application/json",
			//"use_tls":  "false",
			//"Metadata": "",
			//metadata: s,22
			//use_tls: undefined
			//	content-type: application/x-www-form-urlencoded; charset=UTF-8
		},
	})
	xerror.Panic(err)
	app.Log(rsp.OK)

	var data *invokeFuncRsp
	xerror.Panic(json.Unmarshal(rsp.Body, &data))
	return data
}

func (h *Home) functionDescribe(target, selected string) (schema string, template string) {
	rsp, err := fetch.Fetch("server/"+target+"/function/"+selected+"/describe", &fetch.Opts{
		Method: fetch.MethodGet,
	})
	xerror.Panic(err)

	var data = new(struct {
		Data struct {
			Schema   string `json:"schema"`
			Template string `json:"template"`
		} `json:"data"`
	})
	xerror.Panic(json.Unmarshal(rsp.Body, data))
	return data.Data.Schema, data.Data.Template
}

func (h *Home) listFuncs(target, selected string) []string {
	rsp, err := fetch.Fetch("server/"+target+"/service/"+selected+"/functions", &fetch.Opts{
		Method: fetch.MethodGet,
	})
	xerror.Panic(err)
	var data = new(struct {
		Data []string `json:"data"`
	})
	xerror.Panic(json.Unmarshal(rsp.Body, data))
	return data.Data
}

func (h *Home) setReqResData(req *Request) {
	// $('#server-target').val(data.server_target);
	//    target = data.server_target;
	//    $("#select-service").html(new Option(data.selected_service, data.selected_service, true, true));
	//    $('#choose-service').show();
	//    $("#select-function").html(new Option(data.selected_function.substr(data.selected_service.length), data.selected_function, true, true));
	//    $('#choose-function').show();
	//    generate_editor(data.raw_request);
	//    $('#body-request').show();
	//    $('#schema-proto').html(data.schema_proto_html);
	//    $('#json-response').html(data.response_html);
	//    $('#response').show();
	h.req = req
}

func (h *Home) deleteRequest(id string) {
	rsp, err := fetch.Fetch(fmt.Sprintf(`/api/request/%s`, id), &fetch.Opts{
		Method: fetch.MethodDelete,
	})
	xerror.Panic(err)
	app.Log(rsp.OK)
}

func (h *Home) updateRequest(req *Request) {
	var dt, err = json.Marshal(req)
	xerror.Panic(err)

	rsp, err := fetch.Fetch(fmt.Sprintf(`/api/request/%s`, req.ID), &fetch.Opts{
		Method:  fetch.MethodPut,
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    bytes.NewReader(dt),
	})
	xerror.Panic(err)
	app.Log(rsp.OK)
}

func (h *Home) insertRequest(req *Request) {
	var dt, err = json.Marshal(req)
	xerror.Panic(err)

	rsp, err := fetch.Fetch(fmt.Sprintf(`/api/request`), &fetch.Opts{
		Method:  fetch.MethodPost,
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    bytes.NewReader(dt),
	})
	xerror.Panic(err)
	app.Log(rsp.OK)
}

func (h *Home) getAllRequestKey() []*Request {
	rsp, err := fetch.Fetch(fmt.Sprintf(`/api/requests`), &fetch.Opts{
		Method: fetch.MethodGet,
	})
	xerror.Panic(err)

	var req map[string][]*Request
	xerror.Panic(json.Unmarshal(rsp.Body, &req))
	return req["data"]
}

func (h *Home) getRequest(id string) *Request {
	rsp, err := fetch.Fetch(fmt.Sprintf(`/api/request/%s`, id), &fetch.Opts{
		Method: fetch.MethodGet,
	})
	xerror.Panic(err)

	var req map[string]*Request
	xerror.Panic(json.Unmarshal(rsp.Body, &req))
	return req["data"]
}

func (h *Home) search(ctx app.Context, e app.Event) {
	var ee = jsutil.Event(e)
	var li = dom.GetWindow().Document().QuerySelectorAll(".request-list")
	for i := range li {
		if strings.Contains(strings.ToLower(li[i].GetAttribute("request-name")), ee.CurrentTarget().NodeValue()) {
			li[i].SetAttribute("style", "display: ")
		} else {
			li[i].SetAttribute("style", "display: none")
		}
	}
}

func (h *Home) resetReqResData() {
	//var doc = dom.GetWindow().Document()
	//doc.GetElementByID("choose-service").SetAttribute("hidden", "")
	//    target = "";
	//    $('#choose-service').hide();
	//    $('#choose-function').hide();
	//    $('#body-request').hide();
	//    $('#response').hide();
}

//setServerTargetActive

func (h *Home) setServerTargetActive() {
	var dd = dom.GetWindow().Document().QuerySelectorAll("[for='server-target']")
	for i := range dd {
		dd[i].Class().Add("active")
	}
}

func (h *Home) removeRequestSelectedClass() {
	var dd = dom.GetWindow().Document().QuerySelectorAll(".request-list")
	for i := range dd {
		dd[i].Class().Remove("selected")
	}
}

func (h *Home) updateRequestView(id string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		fmt.Println("updateRequestView", id)
		var req = h.data[id]
		fmt.Println(h.data)
		//app.JSValue(e).Set("hidden", true)
		h.resetReqResData()
		h.setReqResData(req)
		h.removeRequestSelectedClass()
		h.setServerTargetActive()
		h.Update()
	}
}

func (h *Home) getReqResData() *Request {
	var serverTarget = dom.GetWindow().Document().GetElementByID("server-target").GetAttribute("value")
	var selectService = dom.GetWindow().Document().GetElementByID("select-service").GetAttribute("value")
	var selectFunction = dom.GetWindow().Document().GetElementByID("select-function").GetAttribute("value")
	var responseHTML = dom.GetWindow().Document().GetElementByID("json-response").InnerHTML()
	var schemaProtoHTML = dom.GetWindow().Document().GetElementByID("schema-proto").InnerHTML()
	var editor = app.Window().Get("ace").Call("edit", "editor").Call("getValue").String()
	return &Request{
		ServerTarget:     serverTarget,
		SelectedService:  selectService,
		SelectedFunction: selectFunction,
		RawRequest:       editor,
		ResponseHTML:     responseHTML,
		SchemaProtoHTML:  schemaProtoHTML,
	}
}

func (h *Home) getActiveRequestListName() *Request {
	var elems = doc.QuerySelectorAll(".request-list")
	for i := range elems {
		if elems[i].Class().Contains("selected") {
			return h.data[elems[i].GetAttribute("request-name")]
		}
	}
	return nil
}

func (h *Home) showModalSaveAsRequest() func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		var req = h.getReqResData()
		doc.GetElementByID("input-request-name").SetAttribute("value", req.SelectedFunction)
		app.Script().Text(`$('#saveRequest').modal('toggle');`)
	}
}

func (h *Home) showModalSaveRequest() func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		var req = h.getReqResData()
		var activeRequestName = h.getActiveRequestListName()
		if activeRequestName == nil {
			fmt.Println(activeRequestName)
			doc.GetElementByID("input-request-name").SetAttribute("value", req.SelectedFunction)
			//doc.GetElementByID("saveRequest").SetAttribute("value", req.SelectedFunction)
			app.Script().Text(`$('#saveRequest').modal('toggle');`)
		} else {
			req.ID = activeRequestName.ID
			req.Name = activeRequestName.Name
			h.updateRequest(req)
		}
	}
}

func (h *Home) saveRequest() func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		var requestName = app.Window().GetElementByID("input-request-name").String()
		if requestName == "" {
			jsutil.Alert("request name is require")
		} else {
			var req = h.getReqResData()
			req.Name = requestName
			h.insertRequest(req)
		}
	}
}

func (h *Home) removeRequestEvent(id string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		h.removeRequest(id)
	}
}

func (h *Home) removeRequest(id string) {
	rsp, err := fetch.Fetch(fmt.Sprintf(`/api/request/%s`, id), &fetch.Opts{
		Method: fetch.MethodDelete,
	})
	xerror.Panic(err)
	app.Log(rsp)
}

func (h *Home) renderRequestList() {
	var ul = dom.GetWindow().Document().GetElementByID("request-list")
	ul.SetInnerHTML("")
	var nameList = h.getAllRequestKey()
	for _, req := range nameList {
		var node = dom.GetWindow().Document().CreateElement("li")
		node.Class().Add("list-group-item")
		node.Class().Add("request-list")
		node.SetAttribute("request-name", req.Name)
		node.AddEventListener("click", true, func(event dom.Event) {
			fmt.Println(event.CurrentTarget().ID())
		})
		app.Log(app.HTMLString(app.A().Title("Delete this request").Class("delete-request").OnClick(h.removeRequestEvent(req.ID)).Body(
			app.I().Class("fa fa-times"))))
		node.SetInnerHTML(`
       <a title="Delete this request" class="delete-request" onclick="removeRequest('${item.id}')"><i class="fa fa-times"></i></a>
       <p class="one-long-line request" onclick="updateRequestView('${item.id}',this)">${item.name}</p>
`)
		ul.AppendChild(node)
	}
}
