package config

type Config struct {
	DatabasePath string
	ServerPort   string
	// Add other configuration fields as needed
}

func LoadConfig() (*Config, error) {
	// TODO: Implement configuration loading
	return &Config{
		DatabasePath: "words.db",
		ServerPort:   ":8080",
	}, nil
}
