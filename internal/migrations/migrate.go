package migrations

import (
	"github.com/pubgo/grpcox/internal/migrations/internal"
	"github.com/pubgo/lava/core/migrates"
)

func Migrations() []migrates.Migrate {
	return []migrates.Migrate{
		internal.M0001,
	}
}
