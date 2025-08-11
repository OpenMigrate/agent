package internal

import (
  "encoding/gob"
  "errors"
  "flag"
  "os"
)

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

func UpdateState(info DiscoverInfo) error {
  state := ".state"
  f, err := os.OpenFile(state, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
  if err != nil {
    return err
  }

  encoder := gob.NewEncoder(f)
  if err := encoder.Encode(info); err != nil {
    f.Close()
    return err
  }
  if err := f.Sync(); err != nil {
    return err
  }
  return f.Close()
}

func GetState() (info DiscoverInfo, err error ) {
  f, err := os.Open(".state")
  if errors.Is(err, os.ErrNotExist) {
    return info, nil
  }
  if err != nil {
    return 
  }
  defer f.Close()

  dec := gob.NewDecoder(f)
  if err = dec.Decode(&info); err != nil{
    return 
  }

  return
}
