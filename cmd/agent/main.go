package main

import (
  "agent/internal"
  "agent/logger"
  "log/slog"
)
func main(){
  cfg := internal.LoadConfig();
  logger.Init(logger.Options{Verbose:cfg.Verbose, File:cfg.LogPath, JSON:false});
  log := logger.Get();

  if err := discover(log, cfg); err != nil {
    log.Error("discovery failed", "err", err)
  }

  err := logger.Close();
  if err != nil {
    log.Error(err.Error())
  }
}

func discover(log *slog.Logger, cfg *internal.Config)(err error){
  // Fetching system infor for sync
  metadata, err := internal.ScanSystem()
  if err != nil {
    return
  }

  // Checking for discovery status from state 
  state, err := internal.GetState()
  if err != nil {
    log.Warn("no local state, will register", "err", err)
  } else if state.Id != "" {
    if err := internal.ValidateDiscovery(*cfg, state.Id); err == nil {
      log.Info("valid state found, skipping registration", "id", state.Id)
      return nil
    }
    log.Warn("state invalid, re-registering", "id", state.Id)
  }


  // Discovering as no local state found
  info, err := internal.SendDiscover(*cfg, metadata)
  if err != nil {
    return
  }
  log.Info("Discovery Done", "uuid", info.Id)

  if err = internal.UpdateState(info); err != nil {
    return
  }

  log.Info("State Saved", "body", info)
  return nil

}
