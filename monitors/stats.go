package monitors

import (
	"log"
	_ "bufio"
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
	Client       *ssh.Client
}

func NewStatsMonitor(client *ssh.Client) *Stats {
	s := new(Stats)
	s.Client = client
	return s
}

func (s *Stats) getUptime() (string, error) {
	uptime, err := utils.RunCommand(s.Client, "/bin/cat /proc/uptime")
	if err != nil {
		log.Fatal(err)
	}

	result := ""
	parts := strings.Fields(uptime)
	if len(parts) == 2 {
		var upsecs float64
		upsecs, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		result = time.Duration(upsecs * 1e9).String()
	}

	return result, err
}

func (s *Stats) Run() string {
	output, _ := s.getUptime()
	return output
}

func RunUptimeMonitor (client *ssh.Client) string {
	uptime, err := utils.RunCommand(client, "/bin/cat /proc/uptime")
	if err != nil {
		log.Fatal(err)
	}

	result := ""
	parts := strings.Fields(uptime)
	if len(parts) == 2 {
		var upsecs float64
		upsecs, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		result = time.Duration(upsecs * 1e9).String()
	}

	return result	
}
