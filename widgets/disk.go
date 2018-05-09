package widgets

import (
	"fmt"
	"time"

	"github.com/cjbassi/gotop/utils"
	ui "github.com/cjbassi/termui"
	psDisk "github.com/shirou/gopsutil/disk"
)

type Disk struct {
	*ui.Gauge
	fs       string // which filesystem to get the disk usage of
	interval time.Duration

	prevRead  uint64
	prevWrite uint64
}

func NewDisk() *Disk {
	self := &Disk{
		Gauge:    ui.NewGauge(),
		fs:       "/",
		interval: time.Second,
	}
	self.Label = "Disk Usage"

	self.update()

	ticker := time.NewTicker(self.interval)
	go func() {
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

func (self *Disk) update() {
	self.updateUsedPercent()
	self.updateIOUsage()
}

func (self *Disk) updateUsedPercent() {
	usage, _ := psDisk.Usage(self.fs)
	self.Percent = int(usage.UsedPercent)
	self.Description = fmt.Sprintf(" (%dGB free)", int(utils.BytesToGB(usage.Free)))
}

func (self *Disk) updateIOUsage() {
	ret, _ := psDisk.IOCounters("/dev/sda")
	data := ret["sda"]
	curRead, curWrite := data.ReadBytes, data.WriteBytes
	if self.prevRead != 0 { // if this isn't the first update
		readRecent := curRead - self.prevRead
		writeRecent := curWrite - self.prevWrite

		readFloat, unitRead := utils.ConvertBytes(readRecent)
		writeFloat, unitWrite := utils.ConvertBytes(writeRecent)
		readRecent, writeRecent = uint64(readFloat), uint64(writeFloat)
		self.Read = fmt.Sprintf("R/s: %3d %2s", readRecent, unitRead)
		self.Write = fmt.Sprintf("W/s: %3d %2s", writeRecent, unitWrite)
	} else {
		self.Read = fmt.Sprintf("R/s: %3d %2s", 0, "B")
		self.Write = fmt.Sprintf("W/s: %3d %2s", 0, "B")
	}
	self.prevRead, self.prevWrite = curRead, curWrite
}
