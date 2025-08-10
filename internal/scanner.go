package internal

import (
  "bytes"
  "fmt"
  "os"
  "os/exec"
  "runtime"
  "strconv"
  "strings"
  "time"
)

type DiskInfo struct {
  Name string `json:"name"`
  Size string `json:"size"`
  Type string `json:"type"`
}

type Metadata struct {
  Hostname        string     `json:"hostname"`
  OS              string     `json:"os"`
  Arch            string     `json:"arch"`
  NumCPU          int        `json:"num_cpu"`
  Kernel          string     `json:"kernel"`
  Uptime          string     `json:"uptime"`
  TotalMemoryMB   uint64     `json:"total_memory_mb"`
  TotalDiskSizeGB string     `json:"total_disk_size_gb"`
  Disks           []DiskInfo `json:"disks"`
  MountedCount    int        `json:"mounted_count"`
  TimestampUTC    string     `json:"timestamp_utc"`
}

// Scans entire OS, parms include
// Hostname
// OS type
// Architecture
// CPU
// Kernel
// Uptime
// Memory & Disk
func ScanSystem() (*Metadata, error) {
  hostname, _ := os.Hostname()
  kernel, _ := exec.Command("uname", "-r").Output()
  uptime, _ := exec.Command("uptime", "-p").Output()

  totalMemMB, _ := getTotalMemory()
  totalDiskGB, disks, mountedCount, _ := getDiskInfo()

  return &Metadata{
    Hostname:        hostname,
    OS:              runtime.GOOS,
    Arch:            runtime.GOARCH,
    NumCPU:          runtime.NumCPU(),
    Kernel:          string(bytes.TrimSpace(kernel)),
    Uptime:          string(bytes.TrimSpace(uptime)),
    TotalMemoryMB:   totalMemMB,
    TotalDiskSizeGB: totalDiskGB,
    Disks:           disks,
    MountedCount:    mountedCount,
    TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
  }, nil
}

func getTotalMemory() (uint64, error) {
  switch runtime.GOOS {
  case "linux":
    out, err := exec.Command("grep", "MemTotal", "/proc/meminfo").Output()
    if err != nil {
      return 0, err
    }
    fields := strings.Fields(string(out))
    if len(fields) < 2 {
      return 0, fmt.Errorf("unexpected meminfo format")
    }
    kb, _ := strconv.ParseUint(fields[1], 10, 64)
    return kb / 1024, nil // MB

  case "darwin":
    out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
    if err != nil {
      return 0, err
    }
    bytesVal, _ := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
    return bytesVal / 1024 / 1024, nil // MB
  }

  return 0, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
}


func getDiskInfo() (totalGB string, disks []DiskInfo, mountedCount int, err error) {
  switch runtime.GOOS {
  case "linux":
    out, err := exec.Command("lsblk", "-b", "-o", "NAME,SIZE,ROTA,MOUNTPOINT").Output()
    if err != nil {
      return "", nil, 0, err
    }
    var totalBytes uint64
    lines := strings.Split(strings.TrimSpace(string(out)), "\n")[1:]
    for _, line := range lines {
      fields := strings.Fields(line)
      if len(fields) < 3 {
        continue
      }
      name := fields[0]
      sizeBytes, _ := strconv.ParseUint(fields[1], 10, 64)
      diskType := "SSD"
      if fields[2] == "1" {
        diskType = "HDD"
      }
      if len(fields) > 3 && fields[3] != "" {
        mountedCount++
      }
      totalBytes += sizeBytes
      disks = append(disks, DiskInfo{
        Name: name,
        Size: fmt.Sprintf("%.2f GB", float64(sizeBytes)/1e9),
        Type: diskType,
      })
    }
    return fmt.Sprintf("%.2f", float64(totalBytes)/1e9), disks, mountedCount, nil
  }

  return "", nil, 0, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
}

