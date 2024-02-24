package enput

func WithConfig(config *Config) option {
	return func(e *Extension) {
		e.data.Config = config
	}
}
