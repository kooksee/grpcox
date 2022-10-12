package bootstrap

import (
	"github.com/pubgo/dix/di"
	"github.com/pubgo/lava/clients/orm"
	"github.com/pubgo/lava/config"

	_ "github.com/pubgo/lava/clients/orm/driver/sqlite"
)

func Init() {
	di.Provide(func(c config.Config) Config {
		return config.Decode[Config](c)
	})

	di.Provide(orm.New)
}
