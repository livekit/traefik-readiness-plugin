// Copyright 2023 LiveKit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package traefik_readiness_plugin

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/livekit/traefik-readiness-plugin/hwstats"
)

// Config the plugin configuration.
type Config struct {
	ReadyPath     string  `json:"ready_path,omitempty"`
	ReadyCPULimit float64 `json:"ready_cpu_limit,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		ReadyPath:     "/ready",
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
	cpuStats, err := hwstats.NewCPUStats()
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
		os.Stdout.WriteString(fmt.Sprintf("readiness plugin recieved request: %v\n", req.URL.Path))

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

		rw.Write([]byte("\n"))
		rw.Write([]byte(fmt.Sprintf("Num CPUs: %v\n", p.cpuStats.NumCPU())))
		rw.Write([]byte(fmt.Sprintf("CPU Load: %v / %v\n", cpuLoad, p.readyCPULimit)))

		return
	}

	p.next.ServeHTTP(rw, req)
}
