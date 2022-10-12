package models

import "gorm.io/gorm/schema"

var _ schema.Tabler = new(Project)

type Project struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Target string `json:"target"`
}

func (p Project) TableName() string {
	return "sys_projects"
}
