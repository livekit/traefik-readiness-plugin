package traefikreadinessplugin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/livekit/protocol/utils/hwstats"
)

// Config the plugin configuration.
type Config struct {
	ReadyPath     string  `json:"ready_path,omitempty"`
	ReadyCPULimit float64 `json:"ready_cpu_limit,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		ReadyPath:     "/ping/ready",
		ReadyCPULimit: 0.8,
	}
}

type Readiness struct {
	next     http.Handler
	name     string
	cpuStats *hwstats.CPUStats

	readyPath     string
	readyCPULimit float64
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	cpuStats, err := hwstats.NewCPUStats(nil)
	if err != nil {
		return nil, err
	}

	r := &Readiness{
		next:     next,
		name:     name,
		cpuStats: cpuStats,

		readyPath:     config.ReadyPath,
		readyCPULimit: config.ReadyCPULimit,
	}

	return r, nil
}

func (p *Readiness) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == p.readyPath {
		var cpuLoad float64
		cpuIdle := p.cpuStats.GetCPUIdle()
		if cpuIdle > 0 {
			cpuLoad = 1 - (cpuIdle / p.cpuStats.NumCPU())
		}

		if cpuLoad > p.readyCPULimit {
			rw.WriteHeader(http.StatusNotAcceptable)
			rw.Write([]byte("Not ready: CPU Limit reached\n"))
		} else {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("Ready\n"))
		}

		rw.Write([]byte("/n"))
		rw.Write([]byte(fmt.Sprintf("Num CPUs: %v\n", p.cpuStats.NumCPU())))
		rw.Write([]byte(fmt.Sprintf("CPU Load: %v / %v\n", cpuLoad, p.readyCPULimit)))

		return
	}

	p.next.ServeHTTP(rw, req)
}
