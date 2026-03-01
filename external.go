package trace

import . "github.com/infrago/base"

func Write(span Span) {
	module.Write(span)
}

func RegisterDriver(name string, driver Driver) {
	module.RegisterDriver(name, driver)
}

func RegisterConfig(name string, cfg Config) {
	module.RegisterConfig(name, cfg)
}

func RegisterConfigs(configs Configs) {
	module.RegisterConfigs(configs)
}

func Stats() Map {
	return module.Stats()
}
