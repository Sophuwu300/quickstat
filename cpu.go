package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type CPU struct {
	Load float64
	MHz  float64
	Temp int
}

func (cpu *CPU) loadTemp() {
	if cpu.Temp == -100 {
		return
	}
	var findCPUTypeNo = func(path string, fileName string, comp string) (string, error) {
		var b = make([]byte, 100)
		var i int = 0
		if fileName == "_label" {
			i = 1
		}
		var err error = nil
		for ; err == nil && string(b[:len(comp)]) != comp; i++ {
			b, err = os.ReadFile(path + string(i+48) + fileName)
		}
		if err != nil {
			return "", err
		}
		return string(i + 48 - 1), nil
	}
	Err := func(err error) bool {
		if err == nil {
			return false
		}
		cpu.Temp = -100
		return true
	}
	var (
		b         = make([]byte, 2)
		thrmPath  = "/sys/class/thermal/thermal_zone"
		hwMonPath = "/sys/class/hwmon/hwmon"
	)
	nStr, err := findCPUTypeNo(thrmPath, "/type", "x86_pkg_temp")
	if err != nil {
		nStr, err = findCPUTypeNo(hwMonPath, "/name", "k10temp")
		if Err(err) {
			return
		}
		nStrHW := nStr
		nStr, err = findCPUTypeNo(hwMonPath+nStrHW+"/temp", "_label", "Tdie")
		if Err(err) {
			return
		}
		b, err = os.ReadFile(hwMonPath + nStrHW + "/temp" + nStr + "_input")
	} else {
		b, err = os.ReadFile(thrmPath + nStr + "/temp")
	}
	if Err(err) {
		return
	}
	cpu.Temp = (int(b[0])-48)*10 + int(b[1]) - 48
}
func (c *CPU) loadMHz() {
	b, err := os.ReadFile("/proc/cpuinfo")

	if errCheck(err, func() string {
		cpuMHz = -1
		return "Unable to read cpu MHz.\n"
	}) {
		return
	}

	var tmp, j int
	c.MHz = 0
	for i := 0; i < len(b); i++ {
		if string(b[i:i+7]) == "cpu MHz" {
			tmp = 0
			for ; i < len(b) && b[i] != 10 && b[i] != '.'; i++ {
				if b[i] >= 48 && b[i] <= 57 {
					tmp = tmp*10 + int(b[i]) - 48
				}
			}
			cpuMHz += tmp
			j++
		}
	}
	c.MHz /= j
}
func (c *CPU) loadUsage() {
	readStat := func(n *[4]float64) bool {
		errFunc := func() string {
			cpuLoad = -1
			return "Unable to read cpu load.\n"
		}

		b := make([]byte, 100)
		f, err := os.Open("/proc/stat")
		if errCheck(err, errFunc) {
			return true
		}
		_, err = f.Read(b)
		if errCheck(err, errFunc) {
			return true
		}
		f.Close()
		for i, j := 6, 0; j < 4; i++ {
			if b[i] >= 48 && b[i] <= 57 {
				n[j] = n[j]*10 + float64(b[i]) - 48
			} else if b[i] == ' ' {
				j++
			}
		}
		return false
	}
	var a, b [4]float64
	if readStat(&a) {
		return
	}
	time.Sleep(CONFIG.Time)
	if readStat(&b) {
		return
	}
	c.Load = ((b[0] + b[1] + b[2]) - (a[0] + a[1] + a[2])) / ((b[0] + b[1] + b[2] + b[3]) - (a[0] + a[1] + a[2] + a[3]))
}
func (c *CPU) update() {
	go c.loadMHz()
	go c.loadTemp()
	c.loadUsage()
}
func (c *CPU) GHz() float64 {
	return c.MHz / 1000
}
func (c *CPU) TempStr() string {
	if c.Temp == -100 {
		return ""
	}
	return fmt.Sprintf("%d °C", c.Temp)
}
func (c *CPU) String() string {
	if CONFIG.Json {
		b, _ := json.Marshal(c)
		return string(b)
	}
	return fmt.Sprintf("%.1f %c %.2f GHz %s\n", c.Load*100, '%', c.GHz(), c.TempStr())
}
func Cpu() CPU {
	var cpu CPU
	cpu.update()
	return cpu
}