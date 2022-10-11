package models

type Request struct {
	ID               uint              `json:"id"`
	Name             string            `json:"name"`
	Metadata         map[string]string `json:"metadata"`
	RawRequest       string            `json:"raw_request"`
	ResponseHtml     string            `json:"response_html"`
	SchemaProtoHtml  string            `json:"schema_proto_html"`
	SelectedFunction string            `json:"selected_function"`
	SelectedService  string            `json:"selected_service"`
	Client           *Client           `json:"client"`
}
