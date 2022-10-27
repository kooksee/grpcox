package main

import (
	"github.com/pubgo/grpcox/internal/bootstrap"
	"github.com/pubgo/grpcox/internal/migrations"
	"github.com/pubgo/lava"

	"github.com/pubgo/lava/cmds/migratecmd"
)

func main() {
	bootstrap.Init()
	lava.Run(
		migratecmd.New(migrations.Migrations()),
	)
}
