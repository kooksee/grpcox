package projectsrv

import "github.com/pubgo/grpcox/internal/models"

type ListProjectReq struct {
}

type ListProjectRsp struct {
	Projects []*models.Project
}

type SaveProjectReq struct {
	Project *models.Project
}

type SaveProjectRsp struct {
	Project *models.Project
}

type GetProjectReq struct {
	ID uint
}

type GetProjectRsp struct {
	Project *models.Project
}

type DeleteProjectReq struct {
	ID uint
}

type DeleteProjectRsp struct {
	Project *models.Project
}

type UpdateProjectReq struct {
	Project *models.Project
}

type UpdateProjectRsp struct {
	Project *models.Project
}
