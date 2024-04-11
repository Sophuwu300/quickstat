package main

import (
	"os"
)

type MEM struct {
	Total  Bytes
	Free   Bytes
	Avail  Bytes
	Buffer Bytes
}
type numSeeker struct {
	i int
	b []byte
}

func (n *numSeeker) Init(b []byte) {
	n.b = b
	n.i = 0
}
func (n *numSeeker) End() bool {
	return n.i >= len(n.b)
}
func (n *numSeeker) Seek() {
	n.i++
}
func (n *numSeeker) GetByte() uint8 {
	return uint8(n.b[n.i])
}
func (n *numSeeker) IsNum() bool {
	e := n.GetByte()
	return e >= 48 && e <= 57
}
func (n *numSeeker) SeekToNum() {
	for n.Seek(); !n.End(); n.Seek() {
		if n.IsNum() {
			break
		}
	}
}
func (n *numSeeker) GetNum() uint64 {
	var num uint64 = 0
	for n.SeekToNum(); !n.End(); n.Seek() {
		if n.IsNum() {
			num = num*uint64(10) + uint64(uint8(n.GetByte())-48)
		} else {
			return num
		}
	}
	return num
}
func (n *numSeeker) GetNums() []uint64 {
	nums := make([]uint64, 0)
	for ; !n.End(); n.Seek() {
		nums = append(nums, n.GetNum())
	}
	return nums
}

func (m *MEM) update(done chan bool) {

	b := make([]byte, 140)
	f, _ := os.Open("/proc/meminfo")
	_, err := f.Read(b)
	if err != nil {
		return
	}
	f.Close()
	var seeker numSeeker
	seeker.Init(b)
	m.Total = Bytes(seeker.GetNum())
	m.Free = Bytes(seeker.GetNum())
	m.Avail = Bytes(seeker.GetNum())
	m.Buffer = Bytes(seeker.GetNum() + seeker.GetNum())
	done <- true
}