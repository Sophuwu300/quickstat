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
}{false, false, 0, time.Second}

func Help() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	fmt.Println("  -j  Output in JSON format")
	fmt.Println("  -r  Repeat output")
	fmt.Printf("  -[%s] Print in unit\n", SI)
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
				if strings.ContainsAny(arg, "h?") {
					Help()
				}
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

func (b Bytes) Unit() string {
	n := float64(b)
	if SI[CONFIG.Unit] == '%' {
		n /= 100
	}
	for i := 1; i < CONFIG.Unit; i++ {
		n /= 1024
	}
	f := "%.0f"
	if n < 100 {
		f = "%.2f"
	}
	return fmt.Sprintf(f+" %c", n, SI[CONFIG.Unit])
}

func (b Bytes) String() string {
	if CONFIG.Json {
		json.Marshal(b * 1024)
	}
	return b.Unit()
}

type HWInfo struct {
	Time int64 `json:"UnixMilli"`
	CPU  CPU   `json:"CPU"`
	MEM  MEM   `json:"MEM"`
}

func (this *HWInfo) Update() {
	done := make(chan bool)
	go func() {
		this.CPU.update()
		done <- true
	}()
	go func() {
		this.MEM.update()
		done <- true
	}()
	<-done
	<-done
	this.Time = time.Now().UnixMilli()
}
func (i HWInfo) String() string {
	if CONFIG.Json {
		b, _ := json.Marshal(i)
		return string(b)
	}
	s := i.MEM.String()
	if SI[CONFIG.Unit] == '%' {
		s = i.MEM.Percent().String()
	}
	return fmt.Sprintf("MEM: %s | CPU: %s", s, i.CPU.String())
}

func Is(err error) bool {
	return err != nil
}

func main() {
	var hw HWInfo
	for do := true; do; do = CONFIG.Repeat {
		hw.Update()
		fmt.Println(hw)
	}
}