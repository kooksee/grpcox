package projectsrv

import "context"

type Service interface {
	ListProject(ctx context.Context, req *ListProjectReq) (*ListProjectRsp, error)
	SaveProject(ctx context.Context, req *SaveProjectReq) (*SaveProjectRsp, error)
	GetProject(ctx context.Context, req *GetProjectReq) (*GetProjectRsp, error)
	DeleteProject(ctx context.Context, req *DeleteProjectReq) (*DeleteProjectRsp, error)
	UpdateProject(ctx context.Context, req *UpdateProjectReq) (*UpdateProjectRsp, error)
}
