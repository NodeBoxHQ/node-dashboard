package services

import (
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/services/config"
	"github.com/NodeboxHQ/node-dashboard/services/dusk"
	"github.com/NodeboxHQ/node-dashboard/services/linea"
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
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" alt="icon">
                                    <g fill="#ffffff"><path d="m21.25 12.75c.42 0 .75-.34.75-.75 0-.42-.33-.75-.75-.75h-1.25v-2.2h1.25c.42 0 .75-.33.75-.75 0-.41-.33-.75-.75-.75h-1.48c-.48-1.59-1.73-2.84-3.32-3.32v-1.48c0-.41-.34-.75-.75-.75s-.75.34-.75.75v1.25h-2.2v-1.25c0-.41-.34-.75-.75-.75s-.75.34-.75.75v1.25h-2.19v-1.25c0-.41-.34-.75-.75-.75-.42 0-.75.34-.75.75v1.48c-1.6.48-2.85 1.73-3.33 3.32h-1.48c-.41 0-.75.34-.75.75 0 .42.34.75.75.75h1.25v2.2h-1.25c-.41 0-.75.33-.75.75 0 .41.34.75.75.75h1.25v2.2h-1.25c-.41 0-.75.33-.75.75 0 .41.34.75.75.75h1.48c.47 1.59 1.73 2.84 3.33 3.32v1.48c0 .41.33.75.75.75.41 0 .75-.34.75-.75v-1.25h2.2v1.25c0 .41.33.75.75.75.41 0 .75-.34.75-.75v-1.25h2.19v1.25c0 .41.34.75.75.75s.75-.34.75-.75v-1.48c1.59-.48 2.84-1.73 3.32-3.32h1.48c.42 0 .75-.34.75-.75 0-.42-.33-.75-.75-.75h-1.25v-2.2zm-3.99 1.51c0 1.65-1.35 3-3 3h-4.52c-1.65 0-3-1.35-3-3v-4.52c0-1.65 1.35-3 3-3h4.52c1.65 0 3 1.35 3 3z"/><path d="m10.02 16.2483h3.97c1.25 0 2.27-1.01 2.27-2.27v-3.97c0-1.25002-1.01-2.27002-2.27-2.27002h-3.97c-1.25 0-2.27 1.01-2.27 2.27002v3.97c0 1.26 1.01 2.27 2.27 2.27z"/></g>
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
								<svg viewBox="0 0 1024.047 1068.467" xmlns="http://www.w3.org/2000/svg">
									<path fill="#ffffff" d="M0 207.515h1024v267.45l-22.108 1.07c-.297-.014-.644-.022-.994-.022-11.722 0-21.225 9.503-21.225 21.225s9.503 21.225 21.225 21.225c.35 0 .697-.008 1.042-.025l-.05.002 22.157 1.07v356.4H0V524.585l16.384-5.073c9.562-3.037 16.37-11.836 16.37-22.225s-6.807-19.188-16.205-22.18l-.166-.045L0 469.942zm977.456 46.545H46.546v183.994c19.836 12.543 32.816 34.36 32.816 59.206s-12.98 46.663-32.53 59.037l-.285.17V829.36h930.91V560.33c-25.15-10.017-42.863-33.6-44.167-61.472l-.006-.155v-2.048c.948-28.358 18.77-52.325 43.703-62.253l.47-.165zM222.023 385.318h313.25v224.35h-313.25zm266.705 46.546h-220.16v131.258h220.16zm0-46.546h313.25v224.35H488.73zm266.706 46.546h-220.16v131.258h220.16zm-575.768 420.77H133.12v-133.12h46.546zm177.804 0h-46.545v-133.12h46.545zm177.804 0H488.73v-133.12h46.544zm177.803 0h-46.545v-133.12h46.545zm177.804 0h-46.544v-133.12h46.545z"/>
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
									<svg fill="#ffffff" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="m8 16.5a1 1 0 1 0 1 1 1 1 0 0 0 -1-1zm4-14.5c-4 0-8 1.37-8 4v12c0 2.63 4 4 8 4s8-1.37 8-4v-12c0-2.63-4-4-8-4zm6 16c0 .71-2.28 2-6 2s-6-1.29-6-2v-3.27a13.16 13.16 0 0 0 6 1.27 13.16 13.16 0 0 0 6-1.27zm0-6c0 .71-2.28 2-6 2s-6-1.29-6-2v-3.27a13.16 13.16 0 0 0 6 1.27 13.16 13.16 0 0 0 6-1.27zm-6-4c-3.72 0-6-1.29-6-2s2.28-2 6-2 6 1.29 6 2-2.28 2-6 2zm-4 2.5a1 1 0 1 0 1 1 1 1 0 0 0 -1-1z"/></svg>
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
                            <div class="w-[45px] ml-2">
								<svg height="36" viewBox="0 0 12 12" width="36" xmlns="http://www.w3.org/2000/svg">
									<path d="m6 1c2.76142375 0 5 2.23857625 5 5s-2.23857625 5-5 5-5-2.23857625-5-5 2.23857625-5 5-5zm0 1c-2.209139 0-4 1.790861-4 4s1.790861 4 4 4 4-1.790861 4-4-1.790861-4-4-4zm-.5 1.5c.24545989 0 .44960837.17687516.49194433.41012437l.00805567.08987563v2h1.5c.27614237 0 .5.22385763.5.5 0 .24545989-.17687516.44960837-.41012437.49194433l-.08987563.00805567h-2c-.24545989 0-.44960837-.17687516-.49194433-.41012437l-.00805567-.08987563v-2.5c0-.27614237.22385763-.5.5-.5z" fill="#ffffff"/></svg>
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

func GetActivity(config *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.Node == "Linea" {
			activityTemplate := `
        		<div class="w-2/6 h-7 mt-5 rounded-full overflow-hidden relative m-0" hx-get="/metrics/activity" hx-trigger="every 1s" hx-swap="outerHTML" id="activity-bar">
            		<div class="absolute top-0 left-0 w-full z-0 h-full bg-progressBarBackgroundColor rounded-full"></div>
            		<div class="absolute top-0 left-0 h-full rounded-[10px] transition-[width] w-full z-10 bg-%s-500"></div>
            		<div class="items-center text-sm font-bold text-textColor absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-20"> %s - Node %s </div>
        		</div>
			`

			status := linea.NodeStatus()
			color := ""
			adjective := ""

			if status.Failure {
				color = "red"
				adjective = "Offline"
			} else if status.Syncing {
				color = "yellow"
				adjective = "Syncing"
			} else {
				color = "green"
				adjective = "Online"
			}

			activityTemplate = activityTemplate +
				`
			<script>
				tippy('#activity-bar', {
					content: '<b>Node</b> - %s <br> <b>Owner</b> - %s <br> <b>IPv4</b> - %s <br> <b>IPv6</b> - %s',
					allowHTML: true,
				})
			</script>
			`

			return c.SendString(fmt.Sprintf(activityTemplate, color, config.Node, adjective, config.Node, config.Owner, config.IPv4, config.IPv6))
		} else if config.Node == "Dusk" {
			activityTemplate := `
        		<div class="w-2/6 h-7 mt-5 rounded-full overflow-hidden relative m-0" hx-get="/metrics/activity" hx-trigger="every 1s" hx-swap="outerHTML" id="activity-bar">
            		<div class="absolute top-0 left-0 w-full z-0 h-full bg-progressBarBackgroundColor rounded-full"></div>
            		<div class="absolute top-0 left-0 h-full rounded-[10px] transition-[width] w-full z-10 bg-%s-500"></div>
            		<div class="items-center text-sm font-bold text-textColor absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-20"> %s - Node %s </div>
        		</div>
			`

			status := dusk.NodeStatus()
			color := ""
			adjective := ""

			if status.Failure {
				color = "red"
				adjective = "Offline"
			} else {
				color = "green"
				adjective = "Online"
			}

			activityTemplate = activityTemplate +
				`
			<script>
				tippy('#activity-bar', {
					content: '<b>Node</b> - %s <br> <b>Owner</b> - %s <br> <b>IPv4</b> - %s <br> <b>IPv6</b> - %s <br> <b>Height</b> - %d',
					allowHTML: true,
				})
			</script>
			`

			return c.SendString(fmt.Sprintf(activityTemplate, color, config.Node, adjective, config.Node, config.Owner, config.IPv4, config.IPv6, status.Height))
		} else {
			return c.SendString("")
		}
	}
}
