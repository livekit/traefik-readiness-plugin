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

//go:build !linux

package hwstats

import (
	"os"
	"runtime"
)

type nullStatCPUMonitor struct{}

func (p *nullStatCPUMonitor) getCPUIdle() (float64, error) {
	return float64(runtime.NumCPU()), nil
}

func (p *nullStatCPUMonitor) numCPU() float64 {
	return float64(runtime.NumCPU())
}

func newPlatformCPUMonitor() (platformCPUMonitor, error) {
	os.Stderr.WriteString("CPU monitoring unsupported on current platform. Server capacity management will be disabled")

	return &nullStatCPUMonitor{}, nil
}
