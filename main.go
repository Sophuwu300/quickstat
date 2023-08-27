package main

import (
	"os"
	"time"
)

var (
	cpuLoad,
	cpuMHz,
	cpuTemp int
	totalMem,
	usedMem,
	takenMem,
	freeMem,
	availMem,
	prctMem,
	prctBuff,
	prctFree float64
)

func errCheck(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func numFromLine(i int, b []byte, n *float64) int {
	for ; i < len(b); i++ {
		if b[i] == 10 {
			return i + 1
		}
		if b[i] >= 48 && b[i] <= 57 {
			*n = *n*10 + float64(b[i]) - 48
		}
	}
	return i
}

func getMem() {
	b := make([]byte, 84)
	f, err := os.Open("/proc/meminfo")
	errCheck(err)
	_, err = f.Read(b)
	errCheck(err)
	f.Close()

	i := numFromLine(0, b, &totalMem)
	totalMem /= 1000
	i = numFromLine(i, b, &freeMem)
	freeMem /= 1000
	_ = numFromLine(i, b, &availMem)
	availMem /= 1000

	usedMem = totalMem - freeMem
	takenMem = totalMem - availMem
	prctMem = takenMem * 100 / totalMem
	prctBuff = (usedMem - takenMem) * 100 / totalMem
	prctFree = freeMem * 100 / totalMem
}

func getCPUTemp() {

	isX86 := func(b []byte) bool {
		if len(b) < 12 {
			return false
		}
		for i, v := range []byte{120, 56, 54, 95, 112, 107, 103, 95, 116, 101, 109, 112} {
			if b[i] != v {
				return false
			}
		}
		return true
	}

	var (
		err error = nil
		i   int   = 0
		b         = make([]byte, 12)
	)
	for ; err == nil && !isX86(b); i++ {
		b, err = os.ReadFile("/sys/class/thermal/thermal_zone" + string(i+48) + "/type")
	}
	errCheck(err)
	b, err = os.ReadFile("/sys/class/thermal/thermal_zone" + string(i+48-1) + "/temp")
	errCheck(err)
	cpuTemp = (int(b[0])-48)*10 + int(b[1]) - 48
}

func getCPU() {
	readStat := func(n *[4]float64) {
		b := make([]byte, 100)
		f, err := os.Open("/proc/stat")
		errCheck(err)
		_, err = f.Read(b)
		errCheck(err)
		f.Close()
		for i, j := 6, 0; j < 4; i++ {
			if b[i] >= 48 && b[i] <= 57 {
				n[j] = n[j]*10 + float64(b[i]) - 48
			} else if b[i] == ' ' {
				j++
			}
		}
	}
	var a, b [4]float64
	readStat(&a)
	time.Sleep(1 * time.Second)
	readStat(&b)
	cpuLoad = int(100 * ((b[0] + b[1] + b[2]) - (a[0] + a[1] + a[2])) / ((b[0] + b[1] + b[2] + b[3]) - (a[0] + a[1] + a[2] + a[3])))
}

func getMHz() {
	b, err := os.ReadFile("/proc/cpuinfo")
	errCheck(err)
	var tmp, j int
	cpuMHz = 0
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
	cpuMHz /= j
}

func round(f float64) string {
	var n int = int(f*1000) % 1000
	var s string = ""

	for i := 0; i < 3; i++ {
		if n%10 >= 5 {
			n += 10
		}
		n /= 10
	}

	n = int(f) + n
	for n > 0 {
		s = string(n%10+48) + s
		n /= 10
	}

	return s
}

func main() {
	go getMem()
	//go getCPUTemp()
	go getMHz()
	getCPU()
	println("   Mem:  ", round(prctMem), "% USED  ", round(prctBuff), "% BUFF  ", round(prctFree), "% FREE   CPU:  ", cpuLoad, "%  ", cpuTemp, "C  ", cpuMHz, "MHz")
}