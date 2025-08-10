package internal

import "flag"

type Config struct {
	ServerUrl string
	Verbose   bool
	LogPath   string
	JSON      bool   // Optional: JSON logging format
}

// LoadConfig initializes config with command line flags
func LoadConfig() *Config {
	config := &Config{}
	
  flag.StringVar(&config.ServerUrl, "server-url", "http://localhost:4000", "Replicate Server URL")
	flag.BoolVar(&config.Verbose, "v", false, "Enable verbose logging (info level)")
	flag.BoolVar(&config.Verbose, "debug", false, "Enable debug logging (debug level)")
	flag.StringVar(&config.LogPath, "out", "", "Log file path (empty for stdout)")
	flag.BoolVar(&config.JSON, "json", false, "Use JSON log format")
	
	flag.Parse()
	return config
}
