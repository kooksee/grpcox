package bootstrap

import (
	"github.com/pubgo/lava/clients/orm"
)

type Config struct {
	Db *orm.Cfg `yaml:"orm"`
}
