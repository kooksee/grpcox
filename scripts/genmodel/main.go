package main

import (
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:        "./internal/models/query",
		FieldNullable:  true,
		FieldCoverable: true,
		Mode:           gen.WithQueryInterface | gen.WithDefaultQuery | gen.WithoutContext,
	})

	g.ApplyBasic(
	//&models.Action{},
	//&models.Endpoint{},
	)

	g.Execute()
}
