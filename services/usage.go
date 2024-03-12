package services

import (
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	uptime2 "github.com/mackerelio/go-osstat/uptime"
	"io/ioutil"
	"math"
	"os"
	"syscall"
	"time"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type TotalDiskUsage struct {
	All  uint64 `json:"All"`
	Used uint64 `json:"Used"`
	Free uint64 `json:"Free"`
}

func DiskUsage(path string) (disk TotalDiskUsage) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

func GetCPUUsage(c *fiber.Ctx) error {
	cpuUsageTemplate := `
                <div class="cursor-pointer items-center py-2.5 px-5 border backdrop-blur-md border-cardBackgroundColor rounded-[20px] bg-cardBackgroundColor w-full shadow-md" hx-get="/metrics/cpu" hx-trigger="load" hx-swap="outerHTML">
                    <h3 class="mb-2.5 text-center text-cardTitleColor text-lg font-semibold">CPU</h3>
                    <div class="flex flex-col">
                        <div class="flex content-between items-center gap-1.5 justify-around flex-col">
                            <div class="w-[45px]">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="#ffffff" alt="icon">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m0-10.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.75c0 5.592 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.57-.598-3.75h-.152c-3.196 0-6.1-1.25-8.25-3.286Zm0 13.036h.008v.008H12v-.008Z" />
                                </svg>
                            </div>
                            <div class="pt-[5px] text-center font-normal text-textColor flex flex-col my-auto w-[90%%]"> %v </div>
                            <div class="text-cardSubBodyColor text-sm flex items-center justify-center"> [%%] </div>
                            <div class="w-full h-2.5 rounded-[10px] overflow-hidden relative m-0">
                                <div class="absolute top-0 left-0 w-full z-0 h-full bg-progressBarBackgroundColor rounded-[10px]"></div>
                                <div class="absolute top-0 left-0 h-full rounded-[10px] z-10 transition-[width] bg-progressBarFillColor" style="width: %v%%">
                            </div>
                        </div>
                    </div>
                </div>
	`

	before, err := cpu.Get()
	if err != nil {
		return c.SendString(fmt.Sprintf(cpuUsageTemplate, 0.0, 0.0))
	}
	time.Sleep(time.Second)

	after, err := cpu.Get()
	if err != nil {
		return c.SendString(fmt.Sprintf(cpuUsageTemplate, 0.0, 0.0))
	}

	total := float64(after.Total - before.Total)
	idle := float64(after.Idle - before.Idle)
	usage := 100 * (total - idle) / total

	return c.SendString(fmt.Sprintf(cpuUsageTemplate, math.Round(usage*100)/100, math.Round(usage*100)/100))
}

func GetRAMUsage(c *fiber.Ctx) error {
	ramUsageTemplate := `
       <div class="cursor-pointer items-center py-2.5 px-5 border backdrop-blur-md border-cardBackgroundColor rounded-[20px] bg-cardBackgroundColor w-full shadow-md" hx-get="/metrics/ram" hx-trigger="every 1s" hx-swap="outerHTML">
           <h3 class="mb-2.5 text-center text-cardTitleColor text-lg font-semibold">RAM</h3>
           <div class="flex flex-col">
                        <div class="flex content-between items-center gap-1.5 justify-around flex-col">
                            <div class="w-[45px]">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="#ffffff" alt="icon">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m0-10.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.75c0 5.592 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.57-.598-3.75h-.152c-3.196 0-6.1-1.25-8.25-3.286Zm0 13.036h.008v.008H12v-.008Z" />
                                </svg>
                            </div>
                            <div class="pt-[5px] text-center font-normal text-textColor flex flex-col my-auto w-[90%%]"> %v </div>
                            <div class="text-cardSubBodyColor text-sm flex items-center justify-center"> [%%] </div>
                            <div class="w-full h-2.5 rounded-[10px] overflow-hidden relative m-0">
                                <div class="absolute top-0 left-0 w-full z-0 h-full bg-progressBarBackgroundColor rounded-[10px]"></div>
                                <div class="absolute top-0 left-0 h-full rounded-[10px] z-10 transition-[width] bg-progressBarFillColor" style="width: %v%%">
                            </div>
                        </div>
           </div>
       </div>
	`

	memoryStat, err := memory.Get()

	if err != nil {
		return c.SendString(fmt.Sprintf(ramUsageTemplate, 0.0, 0.0))
	}

	usage := 100 * (float64(memoryStat.Used) / float64(memoryStat.Total))
	return c.SendString(fmt.Sprintf(ramUsageTemplate, math.Round(usage*100)/100, math.Round(usage*100)/100))
}

func GetDiskUsage(c *fiber.Ctx) error {
	diskUsageTemplate := `
       <div class="cursor-pointer items-center py-2.5 px-5 border backdrop-blur-md border-cardBackgroundColor rounded-[20px] bg-cardBackgroundColor w-full shadow-md" hx-get="/metrics/disk" hx-trigger="every 1s" hx-swap="outerHTML">
           <h3 class="mb-2.5 text-center text-cardTitleColor text-lg font-semibold">Disk Usage</h3>
           <div class="flex flex-col">
                        <div class="flex content-between items-center gap-1.5 justify-around flex-col">
                            <div class="w-[45px]">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="#ffffff" alt="icon">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m0-10.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.75c0 5.592 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.57-.598-3.75h-.152c-3.196 0-6.1-1.25-8.25-3.286Zm0 13.036h.008v.008H12v-.008Z" />
                                </svg>
                            </div>
                            <div class="pt-[5px] text-center font-normal text-textColor flex flex-col my-auto w-[90%%]"> %v </div>
                            <div class="text-cardSubBodyColor text-sm flex items-center justify-center"> [%%] </div>
                            <div class="w-full h-2.5 rounded-[10px] overflow-hidden relative m-0">
                                <div class="absolute top-0 left-0 w-full z-0 h-full bg-progressBarBackgroundColor rounded-[10px]"></div>
                                <div class="absolute top-0 left-0 h-full rounded-[10px] z-10 transition-[width] bg-progressBarFillColor" style="width: %v%%">
                            </div>
                        </div>
           </div>
       </div>
	`
	rootUsage := DiskUsage("/")
	totalUsage := TotalDiskUsage{
		All:  rootUsage.All,
		Used: rootUsage.Used,
		Free: rootUsage.Free,
	}

	if _, err := os.Stat("/mnt"); !os.IsNotExist(err) {
		mntFolders, err := ioutil.ReadDir("/mnt")
		if err == nil {
			for _, folder := range mntFolders {
				if !folder.IsDir() {
					continue
				}

				path := "/mnt/" + folder.Name()
				usage := DiskUsage(path)

				totalUsage.All += usage.All
				totalUsage.Used += usage.Used
				totalUsage.Free += usage.Free
			}
		}
	}

	if totalUsage.All > 0 {
		totalUsage.Used = uint64((float64(totalUsage.Used) / float64(totalUsage.All)) * 100)
	}

	return c.SendString(fmt.Sprintf(diskUsageTemplate, totalUsage.Used, totalUsage.Used))
}

func GetSystemUptime(c *fiber.Ctx) error {
	uptimeTemplate := `
       <div class="cursor-pointer items-center py-2.5 px-5 border backdrop-blur-md border-cardBackgroundColor rounded-[20px] bg-cardBackgroundColor w-full shadow-md" hx-get="/metrics/uptime" hx-trigger="every 1s" hx-swap="outerHTML">
           <h3 class="mb-2.5 text-center text-cardTitleColor text-lg font-semibold">Up Time</h3>
           <div class="flex flex-col">
                        <div class="flex content-between items-center gap-1.5 justify-around flex-col">
                            <div class="w-[45px]">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="#ffffff" alt="icon">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m0-10.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.75c0 5.592 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.57-.598-3.75h-.152c-3.196 0-6.1-1.25-8.25-3.286Zm0 13.036h.008v.008H12v-.008Z" />
                                </svg>
                            </div>
                            <div class="pt-[5px] text-center font-normal text-textColor flex flex-col my-auto w-[90%%]"> %s </div>
                        </div>
           </div>
       </div>
	`

	uptime, err := uptime2.Get()

	if err != nil {
		return c.SendString(fmt.Sprintf(uptimeTemplate, "0s"))
	}

	return c.SendString(fmt.Sprintf(uptimeTemplate, utils.SecondsToReadable(int(uptime.Seconds()))))
}
