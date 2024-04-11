package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type float float32

type CPU struct {
	Load float
	MHz  float
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
	if Is(err) {
		return
	}
	var ns numSeeker
	var Ipart, Fpart, n float
	for _, v := range bytes.Split(b, []byte("\n")) {
		Ipart = 0
		Fpart = 0
		if bytes.HasPrefix(v, []byte("cpu MHz")) {
			ns.Init(v)
			Ipart = float(ns.GetNum())
			Fpart = float(ns.GetNum())
			for Fpart >= 1 {
				Fpart /= 10
			}
			c.MHz += Ipart + Fpart
			n++
		}
	}
	c.MHz /= n
}
func (c *CPU) loadUsage() {
	readStat := func(n *[4]float) bool {

		b := make([]byte, 100)
		f, err := os.Open("/proc/stat")
		if Is(err) {
			return true
		}
		_, err = f.Read(b)
		if Is(err) {
			return true
		}
		f.Close()
		for i, j := 6, 0; j < 4; i++ {
			if b[i] >= 48 && b[i] <= 57 {
				n[j] = n[j]*10 + float(b[i]) - 48
			} else if b[i] == ' ' {
				j++
			}
		}
		return false
	}
	var a, b [4]float
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
	c.loadMHz()
	c.loadTemp()
	c.loadUsage()
}
func (c *CPU) GHzStr() string {
	return fmt.Sprintf("%.2f GHz", c.MHz/1000)
}
func (c *CPU) LoadStr() string {
	return fmt.Sprintf("%.1f %c", c.Load*100, '%')
}
func (c *CPU) TempStr() string {
	if c.Temp == -100 {
		return ""
	}
	return fmt.Sprintf("%d C", c.Temp)
}
func (c *CPU) String() string {
	if CONFIG.Json {
		b, _ := json.Marshal(c)
		return string(b)
	}
	return fmt.Sprintf("%6.6s %8.8s %5.5s", c.LoadStr(), c.GHzStr(), c.TempStr())
}