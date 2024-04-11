package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var BitRateUnit = CONFIG.Unit | 1

type BitRate int64

func (b BitRate) String() string {
	n := float(b)
	for i := 0; i <= CONFIG.Unit; i++ {
		n /= 1024
	}
	n /= float(CONFIG.Time.Seconds())
	return fmt.Sprintf("%.2f %cB/s", n, SI[BitRateUnit])
}

type Net struct {
	Tx BitRate
	Rx BitRate
}

func readDev() *Net {
	f, err := os.OpenFile("/proc/net/dev", os.O_RDONLY, 0)
	if Is(err) {
		return &Net{}
	}
	b := make([]byte, 1024)
	ln, err := f.Read(b)
	f.Close()
	if Is(err) {
		return &Net{}
	}
	var ns numSeeker
	var arr []uint64
	var rx, tx BitRate
	for _, v := range strings.Split(string(b[0:ln]), "\n") {
		arr = []uint64{}
		ns.Init([]byte(v))
		arr = ns.GetNums()
		if len(arr) < 16 {
			continue
		}
		tx += BitRate(int64(arr[8]))
		rx += BitRate(int64(arr[0]))
	}
	return &Net{tx, rx}
}
func (n *Net) update() {
	n2 := readDev()
	time.Sleep(CONFIG.Time)
	n = readDev()
	n.Tx = n.Tx - n2.Tx
	n.Rx = n.Rx - n2.Rx
}
func (n *Net) String() string {
	if CONFIG.Json {
		b, _ := json.Marshal(n)
		return string(b)
	}
	return fmt.Sprintf("TX: %s | RX: %s", n.Tx, n.Rx)
}