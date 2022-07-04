package app

import (
	"github.com/gusaul/grpcox/internal/proto/demov1pb"
	"github.com/gusaul/grpcox/web/ace"
	"github.com/gusaul/grpcox/web/jsutil"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Request struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	RawRequest       string `json:"raw_request"`
	ResponseHTML     string `json:"response_html"`
	SchemaProtoHTML  string `json:"schema_proto_html"`
	SelectedFunction string `json:"selected_function"`
	SelectedService  string `json:"selected_service"`
	ServerTarget     string `json:"server_target"`
}

type Home struct {
	app.Compo
	data           map[string]*Request
	picker         *jsutil.FilePicker
	chooseService  bool
	chooseFunction bool
	bodyRequest    bool
	response       bool
	target         string
	curReq         *Request

	tables       map[string]bool
	tableHidden  bool
	editor       *ace.Ace
	expanded     bool
	expanded1    bool
	saveExpanded bool
	services     []string
	input        string
	inputDesc    string
	output       string
	cc           demov1pb.TransportClient
	curSrv       string
	methods      []string
	curMth       string
	allRequests  []*Request
	reqFilter    string
}

func (c *Home) OnAppUpdate(ctx app.Context) {
	if ctx.AppUpdateAvailable() {
		ctx.Reload()
	}
}
