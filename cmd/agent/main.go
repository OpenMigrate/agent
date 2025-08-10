package main

import (
	"agent/internal"
	"agent/logger"
)
func main(){
  cfg := internal.LoadConfig();
  logger.Init(logger.Options{Verbose:cfg.Verbose, File:cfg.LogPath, JSON:false});
  log := logger.Get();

  metadata, err := internal.ScanSystem()
  if err != nil {
    log.Error(err.Error())
  }

  err = internal.SendDiscover(*cfg, metadata)
  if err != nil {
    log.Error(err.Error())
  }

  err = logger.Close();
  if err != nil {
    log.Error(err.Error())
  }
}
