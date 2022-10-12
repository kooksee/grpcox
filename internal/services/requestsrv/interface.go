package requestsrv

import "context"

type Service interface {
	ListRequest(ctx context.Context, req *ListRequestReq) (*ListRequestRsp, error)
	SaveRequest(ctx context.Context, req *SaveRequestReq) (*SaveRequestRsp, error)
	GetRequest(ctx context.Context, req *GetRequestReq) (*GetRequestRsp, error)
	DeleteRequest(ctx context.Context, req *DeleteRequestReq) (*DeleteRequestRsp, error)
	UpdateRequest(ctx context.Context, req *UpdateRequestReq) (*UpdateRequestRsp, error)
}
