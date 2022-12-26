package config

// Section defines interface used by configuration sections
type Section interface {
	IsValid() error
	LoadFromEnv()
}

// LoadOverrides load config parameters override from .env file
// and validates the configuration options
func LoadOverrides(configs []Section) error {

	for _, section := range configs {
		section.LoadFromEnv()
		if err := section.IsValid(); err != nil {
			return err
		}
	}
	return nil
}
