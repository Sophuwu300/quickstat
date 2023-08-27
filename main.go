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

func errCheck(err error, errFunc func() string) bool {
	if err == nil {
		return false
	}
	print("  An error occurred: " + errFunc())
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		print(err.Error() + "\n")
	}
	return true
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
	errFunc := func() string {
		prctMem = -1
		prctBuff = -1
		prctFree = -1
		return "Unable to read memory.\n"
	}

	b := make([]byte, 84)
	f, err := os.Open("/proc/meminfo")
	if errCheck(err, errFunc) {
		return
	}
	_, err = f.Read(b)
	if errCheck(err, errFunc) {
		return
	}
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

func findCPUTypeNo(path string, fileName string, comp string) (string, error) {
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

func getCPUTemp() {
	errFunc := func() string {
		cpuTemp = -1
		return "Unable to read CPU temperature.\n"
	}
	var (
		b         = make([]byte, 2)
		thrmPath  = "/sys/class/thermal/thermal_zone"
		hwMonPath = "/sys/class/hwmon/hwmon"
	)
	nStr, err := findCPUTypeNo(thrmPath, "/type", "x86_pkg_temp")
	if err != nil {
		nStr, err = findCPUTypeNo(hwMonPath, "/name", "k10temp")
		if errCheck(err, errFunc) {
			return
		}
		nStrHW := nStr
		nStr, err = findCPUTypeNo(hwMonPath+nStrHW+"/temp", "_label", "Tdie")
		if errCheck(err, errFunc) {
			return
		}
		b, err = os.ReadFile(hwMonPath + nStrHW + "/temp" + nStr + "_input")
	} else {
		b, err = os.ReadFile(thrmPath + nStr + "/temp")
	}
	if errCheck(err, errFunc) {
		return
	}
	cpuTemp = (int(b[0])-48)*10 + int(b[1]) - 48
}

func getCPU() {
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
	time.Sleep(1 * time.Second)
	if readStat(&b) {
		return
	}
	cpuLoad = int(100 * ((b[0] + b[1] + b[2]) - (a[0] + a[1] + a[2])) / ((b[0] + b[1] + b[2] + b[3]) - (a[0] + a[1] + a[2] + a[3])))
}

func getMHz() {
	b, err := os.ReadFile("/proc/cpuinfo")

	if errCheck(err, func() string {
		cpuMHz = -1
		return "Unable to read cpu MHz.\n"
	}) {
		return
	}

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

	if n != 0 {
		for i := 0; i < 3; i++ {
			if n%10 >= 5 {
				n += 10
			}
			n /= 10
		}
	}

	n = int(f) + n
	if n == 0 {
		return "0"
	}
	for n > 0 {
		s = string(n%10+48) + s
		n /= 10
	}

	return s
}

func printNWithErr(in interface{}) string {
	var n float64
	switch v := in.(type) {
	case int:
		n = float64(v)
		break
	case float64:
		n = v
		break
	default:
		return "Err"
	}
	if n < 0 {
		return "Err"
	}
	return round(n)
}

func main() {
	go getMem()
	go getCPUTemp()
	go getMHz()
	getCPU()
	println("   Mem:  ", printNWithErr(prctMem), "% USED  ", printNWithErr(prctBuff), "% BUFF  ", printNWithErr(prctFree), "% FREE   CPU:  ", printNWithErr(cpuLoad), "%  ", printNWithErr(cpuTemp), "C  ", printNWithErr(cpuMHz), "MHz")
}