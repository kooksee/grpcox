package internal

import (
	"github.com/pubgo/lava/core/migrates"
	"gorm.io/gorm"
)

type Project struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Target string `json:"target"`
}

func (p Project) TableName() string {
	return "sys_projects"
}

func M0001() *migrates.Migration {
	return &migrates.Migration{
		ID: "0001_create_project_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().CreateTable(&Project{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Project{})
		},
	}
}
