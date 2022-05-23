package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gusaul/grpcox/web/jsutil"
	"honnef.co/go/js/dom/v2"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pubgo/xerror"
	fetch "marwan.io/wasm-fetch"
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

func (h *Home) generateEditor(content string) bool {
	if h.editor != nil {
		h.editor.Call("setValue", content)
		return true
	}

	var edi = app.Window().Get("ace").Call("edit", "editor")
	edi.Call("setOptions", map[string]interface{}{"maxLines": app.Window().Get("Infinity").Int()})
	edi.Get("renderer").Call("setScrollMargin", 10, 10, 10, 10)
	edi.Call("setTheme", "ace/theme/github")
	edi.Get("session").Call("setMode", "ace/mode/json")
	edi.Get("renderer").Call("setShowGutter", false)
	h.editor = edi
	return false
}

func (h *Home) serviceClose() {
	//var ip = $(this).siblings("span").text();
	var ip = ""
	rsp, err := fetch.Fetch(`active/close/`+ip, &fetch.Opts{
		Method: fetch.MethodDelete,
	})
	xerror.Panic(err)
	app.Log(rsp.OK)
	//	 $('[data-toggle="tooltip"]').tooltip('hide');
	//            if(res.data.success) {
	//                $parent.remove();
	//                updateCountNum();
	//            }

}

func (h *Home) serviceGet() {
	rsp, err := fetch.Fetch(`active/get`, &fetch.Opts{
		Method: fetch.MethodGet,
	})
	xerror.Panic(err)
	app.Log(rsp.OK)
	//	  $(".connections .title span").html(res.data.length);
	//            $(".connections .nav").html("");
	//            res.data.forEach(function(item){
	//                $list = $("#conn-list-template").clone();
	//                $list.find(".ip").html(item);
	//                $(".connections .nav").append($list.html());
	//            });
	//            refreshToolTip();
	//	console.warn("Failed to update active connections", thrownError)
}

func (h *Home) getServices() {
	h.resetReqResData()
	h.removeRequestSelectedClass()
	var target = h.getValidTarget()
	//	use_tls = "false";
	//    var restart = "0"
	//    if($('#restart-conn').is(":checked")) {
	//        restart = "1"
	//    }
	//    if($('#use-tls').is(":checked")) {
	//        use_tls = "true"
	//    }
	rsp, err := fetch.Fetch(fmt.Sprintf(`"server/"+%s+"/services?restart="+%s`, target, "1"), &fetch.Opts{
		Method: fetch.MethodGet,
	})
	xerror.Panic(err)
	app.Log(rsp.OK)
	//	 $("#select-service").html(new Option("Choose Service", ""));
	//            $.each(res.data, (_, item) => $("#select-service").append(new Option(item, item)));
	//            $('#choose-service').show();

}

func (h *Home) invokeFunc(target, fn string) bool {
	//ctxArr = [];
	//    $(".ctx-metadata-input-field").each(function(index, val){
	//        ctxArr.push($(val).text())
	//    });
	//	var func = $('#select-function').val();
	//    if (func == "") {
	//        return false;
	//    }

	rsp, err := fetch.Fetch("server/"+target+"/function/"+fn+"/invoke", &fetch.Opts{
		Method: fetch.MethodPost,
		Body:   nil,
		//	dataType: "json",

	})
	xerror.Panic(err)
	app.Log(rsp.OK)
	//	$("#json-response").html(PR.prettyPrintOne(res.data.result));
	//            $("#timer-resp span").html(res.data.timer);
	//            $('#response').show();
	//	 $('#response').hide();
	//            xhr.setRequestHeader('use_tls', use_tls);
	//            if(ctxUse) {
	//                xhr.setRequestHeader('Metadata', ctxArr);
	//            }
	//            $(this).html("Loading...");
	//            show_loading();
	return true
}

func (h *Home) selectFunction(target, selected string) bool {
	rsp, err := fetch.Fetch("server/"+target+"/function/"+selected+"/describe", &fetch.Opts{
		Method: fetch.MethodGet,
	})
	xerror.Panic(err)
	app.Log(rsp.OK)

	//generate_editor(res.data.template);
	//$("#schema-proto").html(PR.prettyPrintOne(res.data.schema));
	//$('#body-request').show();
	//xhr.setRequestHeader('use_tls', use_tls);
	return true
}

func (h *Home) selectService(target, selected string) bool {
	rsp, err := fetch.Fetch("server/"+target+"/service/"+selected+"/functions", &fetch.Opts{
		Method: fetch.MethodGet,
	})
	xerror.Panic(err)
	app.Log(rsp.OK)

	//	$("#select-function").html(new Option("Choose Method", ""));
	//            $.each(res.data, (_, item) => $("#select-function").append(new Option(item.substr(selected.length) , item)));
	//            $('#choose-function').show();
	//	  xhr.setRequestHeader('use_tls', use_tls);
	return true
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
