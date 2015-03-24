package golgtm

type Config struct {
	FontBasePath string
}

var config *Config

func init() {
	config = NewConfig()
}

func NewConfig() *Config {
	return &Config{
		FontBasePath: "/Library/Fonts/",
	}
}
