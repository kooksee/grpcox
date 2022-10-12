package projectsrv

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pubgo/grpcox/internal/models"
	"gorm.io/gorm"
)

var _ Service = (*serviceImpl)(nil)

type serviceImpl struct {
	db *gorm.DB
}

func (s *serviceImpl) ListProject(ctx context.Context, req *ListProjectReq) (*ListProjectRsp, error) {
	var rsp = new(ListProjectRsp)
	return rsp, s.db.WithContext(ctx).Find(&rsp.Projects).Error
}

func (s *serviceImpl) SaveProject(ctx context.Context, req *SaveProjectReq) (*SaveProjectRsp, error) {
	if err := s.db.WithContext(ctx).Save(req.Project).Error; err != nil {
		return nil, err
	}
	return &SaveProjectRsp{Project: req.Project}, nil
}

func (s *serviceImpl) GetProject(ctx context.Context, req *GetProjectReq) (*GetProjectRsp, error) {
	var rsp = new(GetProjectRsp)
	return rsp, s.db.WithContext(ctx).Model(models.Project{ID: req.ID}).First(&rsp.Project).Error
}

func (s *serviceImpl) DeleteProject(ctx context.Context, req *DeleteProjectReq) (*DeleteProjectRsp, error) {
	var rsp = new(DeleteProjectRsp)
	var err = s.db.WithContext(ctx).Model(models.Project{ID: req.ID}).First(&rsp.Project).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return rsp, err
}

func (s *serviceImpl) UpdateProject(ctx context.Context, req *UpdateProjectReq) (*UpdateProjectRsp, error) {
	var rsp = new(UpdateProjectRsp)
	err := s.db.WithContext(ctx).Model(models.Project{ID: req.Project.ID}).First(&rsp.Project).Error
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(rsp.Project, req.Project, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, err
	}

	return rsp, s.db.WithContext(ctx).Save(rsp.Project).Error
}
