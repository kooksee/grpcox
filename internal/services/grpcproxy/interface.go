package grpcproxy

import "context"

type Service interface {
	ListRequest(ctx context.Context, req *ListRequestReq) (*ListRequestRsp, error)
	SaveRequest(ctx context.Context, req *SaveRequestReq) (*SaveRequestRsp, error)
	GetRequest(ctx context.Context, req *GetRequestReq) (*GetRequestRsp, error)
	DeleteRequest(ctx context.Context, req *DeleteRequestReq) (*DeleteRequestRsp, error)
	UpdateRequest(ctx context.Context, req *UpdateRequestReq) (*UpdateRequestRsp, error)
	ListService(ctx context.Context, req *ListServiceReq) (*ListServiceRsp, error)
	ListMethod(ctx context.Context, req *ListMethodReq) (*ListMethodRsp, error)
	Invoke(ctx context.Context, req *InvokeReq) (*InvokeRsp, error)
}
