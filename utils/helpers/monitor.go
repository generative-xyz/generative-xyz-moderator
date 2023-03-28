package helpers

import (
	"runtime"
)

type SysMonitor struct {
	Cpu int `json:"cpu"`
	Rot int `json:"routines"`
	Memory uint64 `json:"memory"`
}

func System() SysMonitor {
	mem := &runtime.MemStats{}
 
	cpu := runtime.NumCPU()
	rot := runtime.NumGoroutine()
		// Byte
	runtime.ReadMemStats(mem)
	
	return SysMonitor{
		Cpu: cpu,
		Rot: rot,
		Memory: mem.Alloc,
	}
}
