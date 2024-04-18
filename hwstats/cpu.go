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

package hwstats

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/atomic"
)

// This object returns cgroup quota aware cpu stats. On other systems than Linux,
// it falls back to full system stats

type platformCPUMonitor interface {
	getCPUIdle() (float64, error)
	numCPU() float64
}

type CPUStats struct {
	idleCPUs atomic.Float64
	platform platformCPUMonitor

	idleCallback func(idle float64)
	closeChan    chan struct{}
}

func NewCPUStats(idleUpdateCallback func(idle float64)) (*CPUStats, error) {
	p, err := newPlatformCPUMonitor()
	if err != nil {
		return nil, err
	}

	c := &CPUStats{
		platform:     p,
		idleCallback: idleUpdateCallback,
		closeChan:    make(chan struct{}),
	}

	go c.monitorCPULoad()

	return c, nil
}

func (c *CPUStats) GetCPUIdle() float64 {
	return c.idleCPUs.Load()
}

func (c *CPUStats) NumCPU() float64 {
	return c.platform.numCPU()
}

func (c *CPUStats) Stop() {
	close(c.closeChan)
}

func (c *CPUStats) monitorCPULoad() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeChan:
			return
		case <-ticker.C:
			idle, err := c.platform.getCPUIdle()
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("failed retrieving CPU idle: %v\n", err))
				continue
			}

			c.idleCPUs.Store(idle)
		}
	}
}
