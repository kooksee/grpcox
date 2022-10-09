package internal

import (
	"github.com/pubgo/grpcox/internal/models"
	"github.com/pubgo/lava/core/migrates"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = (*Action)(nil)

type Action struct {
	Id      uint32
	Code    string  `gorm:"size:32;not null;unique"`
	Type    string  `gorm:"size:8;not null"`
	Name    string  `gorm:"size:64;not null"`
	ResType *string `gorm:"size:8"`
}

func (a Action) TableName() string {
	return new(models.Action).TableName()
}

type Endpoint struct {
	Id          uint32
	Protocol    string `gorm:"not null;uniqueIndex:url"`
	Method      string `gorm:"not null;uniqueIndex:url"`
	Path        string `gorm:"not null;uniqueIndex:url"`
	Service     string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description *string
	ActionId    uint32  `gorm:"not null"`
	Action      *Action `gorm:"foreignkey:ActionId"`
	ParentId    *uint32
}

func (a Endpoint) TableName() string {
	return new(models.Endpoint).TableName()
}

func M0001() *migrates.Migration {
	return &migrates.Migration{
		ID: "0001_init",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().CreateTable(&Action{}, &Endpoint{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Action{}, &Endpoint{})
		},
	}
}
