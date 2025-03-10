package config

import (
	"chickChirick/cmd/config"
	globalConfig "chickChirick/pkg/config"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	stockDBC      config.DataBaseConfigInterface
	MigrationPath string
}

func NewConfiguration() DatabaseConfig {
	dbc := config.PrepareDatabaseConfig()

	return DatabaseConfig{
		stockDBC:      &dbc,
		MigrationPath: viper.GetString(globalConfig.DbMigrationPatch),
	}
}

func (d *DatabaseConfig) Host() string {
	return d.stockDBC.Host()
}

func (d *DatabaseConfig) Port() int {
	return d.stockDBC.Port()
}

func (d *DatabaseConfig) Name() string {
	return d.stockDBC.Name()
}

func (d *DatabaseConfig) User() string {
	return d.stockDBC.User()
}

func (d *DatabaseConfig) Password() string {
	return d.stockDBC.Password()
}

func (d *DatabaseConfig) Timezone() string {
	return d.stockDBC.Timezone()
}

func (d *DatabaseConfig) SetHost(host string) {
	d.stockDBC.SetHost(host)
}

func (d *DatabaseConfig) SetPort(port int) {
	d.stockDBC.SetPort(port)
}

func (d *DatabaseConfig) SetName(name string) {
	d.stockDBC.SetName(name)
}

func (d *DatabaseConfig) SetUser(user string) {
	d.stockDBC.SetUser(user)
}

func (d *DatabaseConfig) SetPassword(password string) {
	d.stockDBC.SetPassword(password)
}

func (d *DatabaseConfig) SetTimezone(timezone string) {
	d.stockDBC.SetTimezone(timezone)
}
