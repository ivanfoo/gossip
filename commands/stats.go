package commands

import (
	"bufio"
	"strconv"
	"strings"
	"time"

	"github.com/ivanfoo/gossip/utils"
	"golang.org/x/crypto/ssh"
)

type FSInfo struct {
	MountPoint string
	Used       uint64
	Free       uint64
}

type Stats struct {
	Uptime       time.Duration
	Hostname     string
	Load1        string
	Load5        string
	Load10       string
	RunningProcs string
	TotalProcs   string
	MemTotal     uint64
	MemFree      uint64
	MemBuffers   uint64
	MemCached    uint64
	SwapTotal    uint64
	SwapFree     uint64
	FSInfos      []FSInfo
}

func getAllStats(client *ssh.Client, stats *Stats) {
	getUptime(client, stats)
	getHostname(client, stats)
	getLoad(client, stats)
	getMemInfo(client, stats)
	getFSInfo(client, stats)
}

func getUptime(client *ssh.Client, stats *Stats) (err error) {
	uptime, err := utils.RunCommand(client, "/bin/cat /proc/uptime")
	if err != nil {
		return
	}

	parts := strings.Fields(uptime)
	if len(parts) == 2 {
		var upsecs float64
		upsecs, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return
		}
		stats.Uptime = time.Duration(upsecs * 1e9)
	}

	return
}

func getHostname(client *ssh.Client, stats *Stats) (err error) {
	hostname, err := utils.RunCommand(client, "/bin/hostname -f")
	if err != nil {
		return
	}

	stats.Hostname = strings.TrimSpace(hostname)
	return
}

func getLoad(client *ssh.Client, stats *Stats) (err error) {
	line, err := utils.RunCommand(client, "/bin/cat /proc/loadavg")
	if err != nil {
		return
	}

	parts := strings.Fields(line)
	if len(parts) == 5 {
		stats.Load1 = parts[0]
		stats.Load5 = parts[1]
		stats.Load10 = parts[2]
		if i := strings.Index(parts[3], "/"); i != -1 {
			stats.RunningProcs = parts[3][0:i]
			if i+1 < len(parts[3]) {
				stats.TotalProcs = parts[3][i+1:]
			}
		}
	}

	return
}

func getMemInfo(client *ssh.Client, stats *Stats) (err error) {
	lines, err := utils.RunCommand(client, "/bin/cat /proc/meminfo")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 3 {
			val, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				continue
			}
			val *= 1024
			switch parts[0] {
			case "MemTotal:":
				stats.MemTotal = val
			case "MemFree:":
				stats.MemFree = val
			case "Buffers:":
				stats.MemBuffers = val
			case "Cached:":
				stats.MemCached = val
			case "SwapTotal:":
				stats.SwapTotal = val
			case "SwapFree:":
				stats.SwapFree = val
			}
		}
	}

	return
}

func getFSInfo(client *ssh.Client, stats *Stats) (err error) {
	lines, err := utils.RunCommand(client, "/bin/df -B1")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	flag := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		n := len(parts)
		dev := n > 0 && strings.Index(parts[0], "/dev/") == 0
		if n == 1 && dev {
			flag = 1
		} else if (n == 5 && flag == 1) || (n == 6 && dev) {
			i := flag
			flag = 0
			used, err := strconv.ParseUint(parts[2-i], 10, 64)
			if err != nil {
				continue
			}
			free, err := strconv.ParseUint(parts[3-i], 10, 64)
			if err != nil {
				continue
			}
			stats.FSInfos = append(stats.FSInfos, FSInfo{
				parts[5-i], used, free,
			})
		}
	}

	return
}
