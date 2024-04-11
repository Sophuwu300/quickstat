package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const SI = "%KMG"

var CONFIG = struct {
	Json   bool
	Repeat bool
	Unit   int
	Time   time.Duration
}{false, false, 2, time.Second}

func Help() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	fmt.Println("  -j  Output in JSON format")
	fmt.Println("  -r  Repeat output")
	fmt.Printf("  -[%s] Use unit\n", SI)
	fmt.Println("  -t<n> Set sampling time to n seconds")
	os.Exit(0)
}
func parseTime(s string) time.Duration {
	var n int
	for s = s[strings.Index(s, "t")+1:]; len(s) > 0 && s[0] >= '0' && s[0] <= '9'; s = s[1:] {
		n = n*10 + int(s[0]) - '0'
	}
	return time.Duration(n) * time.Second
}
func init() {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if strings.Count(arg, "-") < 2 {
				if strings.Contains(arg, "j") {
					CONFIG.Json = true
				}
				if strings.Contains(arg, "r") {
					CONFIG.Repeat = true
				}
				for i, c := range SI {
					if strings.Contains(arg, string(c)) {
						CONFIG.Unit = i
					}
				}
				if strings.Contains(arg, "t") {
					CONFIG.Time = parseTime(arg)
				}
			} else {
				Help()
			}
		}
	}
}

type Bytes uint64

func (b Bytes) String() string {
	if CONFIG.Unit == 0 {

	}
	return fmt.Sprint(b)
}

type Net struct {
	Tx Bytes
	Rx Bytes
}
type HWInfo struct {
	CPU CPU
	MEM MEM
	NET Net
}

func (this *HWInfo) update() {
	done := make(chan bool)
	go this.CPU.update(done)
	go this.MEM.update(done)
	<-done
	<-done
	// i.NET.update()
}
func (i *HWInfo) String() string {
	if CONFIG.Json {
		b, _ := json.Marshal(i)
		return string(b)
	}
	return fmt.Sprintf("%s %s %s\n", i.CPU, i.MEM, i.NET)

}

func Is(err error) bool {
	return err != nil
}

func main() {
	var HW HWInfo
	HW.update()
	fmt.Println(HW)
	// go getCPUTemp()
	// go getMHz()
	// getCPU()
	// println("   Mem:  ", printNWithErr(prctMem), "% USED  ", printNWithErr(prctBuff), "% BUFF  ", printNWithErr(prctFree), "% FREE   CPU:  ", printNWithErr(cpuLoad), "%  ", printNWithErr(cpuTemp), "C  ", printNWithErr(cpuMHz), "MHz")
}