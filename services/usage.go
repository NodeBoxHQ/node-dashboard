package services

import (
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/services/babylon"
	"github.com/NodeboxHQ/node-dashboard/services/config"
	"github.com/NodeboxHQ/node-dashboard/services/dusk"
	"github.com/NodeboxHQ/node-dashboard/services/linea"
	"github.com/NodeboxHQ/node-dashboard/services/nulink"
	"github.com/NodeboxHQ/node-dashboard/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	uptime2 "github.com/mackerelio/go-osstat/uptime"
	"io/ioutil"
	"math"
	"os"
	"strings"
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

func GetLogo(config *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.Node == "Linea" {
			return c.SendString(fmt.Sprintf(`<img src="/assets/img/logo/linea-logo.png?nodeip=%s" alt="logo-expanded" class="w-52 h-auto object-contain mx-auto block" />`, config.IPv4))
		} else if config.Node == "Dusk" {
			return c.SendString(fmt.Sprintf(`<img src="/assets/img/logo/dusk-logo.png?nodeip=%s" alt="logo-expanded" class="w-52 h-auto object-contain mx-auto block" />`, config.IPv4))
		} else if config.Node == "Nulink" {
			return c.SendString(fmt.Sprintf(`<img src="/assets/img/logo/nulink-logo.png?nodeip=%s" alt="logo-expanded" class="w-52 h-auto object-contain mx-auto block" />`, config.IPv4))
		} else if config.Node == "Babylon" {
			return c.SendString(fmt.Sprintf(`<img src="/assets/img/logo/babylon-logo.png?nodeip=%s" alt="logo-expanded" class="w-52 h-auto object-contain mx-auto block" />`, config.IPv4))
		} else {
			return c.SendString("")
		}
	}
}

func GetCPUUsage(ipv4 string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cpuUsageTemplate := `
                <div class="items-center py-2.5 px-5 border backdrop-blur-md border-cardBackgroundColor rounded-[20px] bg-cardBackgroundColor w-full shadow-md" hx-get="/data/cpu?nodeip=%s" hx-trigger="load" hx-swap="outerHTML transition:true">
                    <h3 class="mb-2.5 text-center text-cardTitleColor text-lg font-semibold">CPU</h3>
                    <div class="flex flex-col">
                        <div class="flex content-between items-center gap-1.5 justify-around flex-col">
                            <div class="w-[45px]">
								<img src="/assets/img/icons/cpu.gif?nodeip=%s" alt="logo-expanded" class="w-52 h-auto object-contain mx-auto block" />
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
			return c.SendString(fmt.Sprintf(cpuUsageTemplate, ipv4, ipv4, 0.0, 0.0))
		}
		time.Sleep(time.Second)

		after, err := cpu.Get()
		if err != nil {
			return c.SendString(fmt.Sprintf(cpuUsageTemplate, ipv4, ipv4, 0.0, 0.0))
		}

		total := float64(after.Total - before.Total)
		idle := float64(after.Idle - before.Idle)
		usage := 100 * (total - idle) / total

		return c.SendString(fmt.Sprintf(cpuUsageTemplate, ipv4, ipv4, math.Round(usage*100)/100, math.Round(usage*100)/100))
	}
}

func GetRAMUsage(ipv4 string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ramUsageTemplate := `
       <div class="items-center py-2.5 px-5 border backdrop-blur-md border-cardBackgroundColor rounded-[20px] bg-cardBackgroundColor w-full shadow-md" hx-get="/data/ram?nodeip=%s" hx-trigger="every 1s" hx-swap="outerHTML transition:true">
           <h3 class="mb-2.5 text-center text-cardTitleColor text-lg font-semibold">RAM</h3>
           <div class="flex flex-col">
                        <div class="flex content-between items-center gap-1.5 justify-around flex-col">
                            <div class="w-[45px]">
							<img src="/assets/img/icons/memory.gif?nodeip=%s" alt="logo-expanded" class="w-52 h-auto object-contain mx-auto block" />
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
			return c.SendString(fmt.Sprintf(ramUsageTemplate, ipv4, ipv4, 0.0, 0.0))
		}

		usage := 100 * (float64(memoryStat.Used) / float64(memoryStat.Total))
		return c.SendString(fmt.Sprintf(ramUsageTemplate, ipv4, ipv4, math.Round(usage*100)/100, math.Round(usage*100)/100))
	}
}

func GetDiskUsage(ipv4 string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		diskUsageTemplate := `
       <div class="items-center py-2.5 px-5 border backdrop-blur-md border-cardBackgroundColor rounded-[20px] bg-cardBackgroundColor w-full shadow-md" hx-get="/data/disk?nodeip=%s" hx-trigger="every 1s" hx-swap="outerHTML transition:true">
           <h3 class="mb-2.5 text-center text-cardTitleColor text-lg font-semibold">Disk Usage</h3>
           <div class="flex flex-col">
                        <div class="flex content-between items-center gap-1.5 justify-around flex-col">
                            <div class="w-[45px]">
							<img src="/assets/img/icons/storage.gif?nodeip=%s" alt="logo-expanded" class="w-52 h-auto object-contain mx-auto block" />
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

		return c.SendString(fmt.Sprintf(diskUsageTemplate, ipv4, ipv4, totalUsage.Used, totalUsage.Used))
	}
}

func GetSystemUptime(ipv4 string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uptimeTemplate := `
       <div class="items-center py-2.5 px-5 border backdrop-blur-md border-cardBackgroundColor rounded-[20px] bg-cardBackgroundColor w-full shadow-md" hx-get="/data/uptime?nodeip=%s" hx-trigger="every 1s" hx-swap="outerHTML transition:true">
           <h3 class="mb-2.5 text-center text-cardTitleColor text-lg font-semibold">Uptime</h3>
           <div class="flex flex-col">
                        <div class="flex content-between items-center gap-1.5 justify-around flex-col">
                            <div class="w-[45px] ml-2">
								<img src="/assets/img/icons/runtime.gif?nodeip=%s" alt="logo-expanded" class="w-52 h-auto object-contain mx-auto block" />
							</div>
                            <div class="pt-[5px] text-center font-normal text-textColor flex flex-col my-auto w-[90%%]"> %s </div>
                        </div>
           </div>
       </div>
	`

		uptime, err := uptime2.Get()

		if err != nil {
			return c.SendString(fmt.Sprintf(uptimeTemplate, ipv4, ipv4, "0s"))
		}

		return c.SendString(fmt.Sprintf(uptimeTemplate, ipv4, ipv4, utils.SecondsToReadable(int(uptime.Seconds()))))
	}
}

func GetActivity(config *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.Node == "Linea" {
			activityTemplate := `
        		<div class="w-64 h-7 mt-5 rounded-full overflow-hidden relative m-0" hx-get="/data/activity?nodeip=NODE_IP" hx-trigger="every 1s" hx-swap="outerHTML" id="activity-bar" ALPINE_TOOLTIP>
						<div class="absolute top-0 left-0 w-full z-0 h-full bg-progressBarBackgroundColor rounded-full"></div>
						<div class="absolute top-0 left-0 h-full rounded-[10px] transition-[width] w-full z-10 %s"></div>
						<div class="items-center text-xs font-bold text-textColor absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-20"> Node %s </div>
        		</div>
			`

			activityTemplate = strings.Replace(activityTemplate, "NODE_IP", config.IPv4, -1)

			status := linea.NodeStatus()
			color := ""
			adjective := ""

			if status.Failure {
				color = "bg-red-500"
				adjective = "Offline"
			} else if status.Syncing {
				color = "bg-yellow-500"
				adjective = "Syncing"
			} else {
				color = "bg-green-500"
				adjective = "Online"
			}

			tippyContent := fmt.Sprintf("<b>Node</b> - %s <br> <b>Owner</b> - %s<br> <b>Private IPv4</b> - %s <br> <b>Public IPv4</b> - %s <br> <b>Public IPv6</b> - %s <br>", config.Node, config.Owner, config.PrivateIPv4, config.IPv4, config.IPv6)

			needBr := ""

			if !status.Failure {
				if status.Syncing {
					tippyContent = tippyContent + fmt.Sprintf("<b>Current Height</b> - %d <br> <b>Max Height</b> - %d", status.CurrentHeight, status.MaxHeight)
					needBr = "<br>"
				} else {
					tippyContent = tippyContent + fmt.Sprintf("<b>Current Height</b> - %d", status.CurrentHeight)
					needBr = "<br>"
				}
			}

			tippyContent = tippyContent + fmt.Sprintf("%s <b>Dashboard Version</b> - %s", needBr, config.NodeboxDashboardVersion)
			activityTemplate = strings.Replace(activityTemplate, "ALPINE_TOOLTIP", fmt.Sprintf(`tooltip-data="%s"`, tippyContent), -1)
			return c.SendString(fmt.Sprintf(activityTemplate, color, adjective))
		} else if config.Node == "Dusk" {
			activityTemplate := `
        		<div class="w-64 h-7 mt-5 rounded-full overflow-hidden relative m-0" hx-get="/data/activity?nodeip=NODE_IP" hx-trigger="every 1s" hx-swap="outerHTML" id="activity-bar" ALPINE_TOOLTIP>
            		<div class="absolute top-0 left-0 w-full z-0 h-full bg-progressBarBackgroundColor rounded-full"></div>
            		<div class="absolute top-0 left-0 h-full rounded-[10px] transition-[width] w-full z-10 %s"></div>
            		<div class="items-center text-sm font-bold text-textColor absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-20"> Node %s </div>
        		</div>
			`

			activityTemplate = strings.Replace(activityTemplate, "NODE_IP", config.IPv4, -1)

			status := dusk.NodeStatus()
			color := ""
			adjective := ""

			if status.Failure {
				color = "bg-red-500"
				adjective = "Offline"
			} else {
				color = "bg-green-500"
				adjective = "Online"
			}

			tippyContent := fmt.Sprintf("<b>Node</b> - %s <br> <b>Version</b> - %s <br> <b>Owner</b> - %s<br> <b>Private IPv4</b> - %s <br> <b>Public IPv4</b> - %s <br> <b>Public IPv6</b> - %s <br>", config.Node, status.Version, config.Owner, config.PrivateIPv4, config.IPv4, config.IPv6)

			needBr := ""

			if !status.Failure {
				tippyContent = tippyContent + fmt.Sprintf("<b>Current Height</b> - %d", status.Height)
				needBr = "<br>"
			}

			tippyContent = tippyContent + fmt.Sprintf("%s <b>Dashboard Version</b> - %s", needBr, config.NodeboxDashboardVersion)
			activityTemplate = strings.Replace(activityTemplate, "ALPINE_TOOLTIP", fmt.Sprintf(`tooltip-data="%s"`, tippyContent), -1)
			return c.SendString(fmt.Sprintf(activityTemplate, color, adjective))
		} else if config.Node == "Nulink" {
			activityTemplate := `
        		<div class="w-64 h-7 mt-5 rounded-full overflow-hidden relative m-0" hx-get="/data/activity?nodeip=NODE_IP" hx-trigger="every 1s" hx-swap="outerHTML" id="activity-bar" ALPINE_TOOLTIP>
            		<div class="absolute top-0 left-0 w-full z-0 h-full bg-progressBarBackgroundColor rounded-full"></div>
            		<div class="absolute top-0 left-0 h-full rounded-[10px] transition-[width] w-full z-10 %s"></div>
            		<div class="items-center text-sm font-bold text-textColor absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-20"> Node %s </div>
        		</div>
			`

			activityTemplate = strings.Replace(activityTemplate, "NODE_IP", config.IPv4, -1)

			status := nulink.NodeStatus()
			color := ""
			adjective := ""

			if status.Online {
				color = "bg-green-500"
				adjective = "Online"
			} else {
				color = "bg-red-500"
				adjective = "Offline"
			}

			tippyContent := fmt.Sprintf("<b>Node</b> - %s <br> <b>Owner</b> - %s<br> <b>Private IPv4</b> - %s <br> <b>Public IPv4</b> - %s <br> <b>Public IPv6</b> - %s", config.Node, config.Owner, config.PrivateIPv4, config.IPv4, config.IPv6)

			tippyContent = tippyContent + fmt.Sprintf("<br> <b>Dashboard Version</b> - %s", config.NodeboxDashboardVersion)
			activityTemplate = strings.Replace(activityTemplate, "ALPINE_TOOLTIP", fmt.Sprintf(`tooltip-data="%s"`, tippyContent), -1)

			return c.SendString(fmt.Sprintf(activityTemplate, color, adjective))
		} else if config.Node == "Babylon" {
			activityTemplate := `
				<div class="w-64 h-7 mt-5 rounded-full overflow-hidden relative m-0" hx-get="/data/activity?nodeip=NODE_IP" hx-trigger="every 1s" hx-swap="outerHTML" id="activity-bar" ALPINE_TOOLTIP>
            		<div class="absolute top-0 left-0 w-full z-0 h-full bg-progressBarBackgroundColor rounded-full"></div>
            		<div class="absolute top-0 left-0 h-full rounded-[10px] transition-[width] w-full z-10 %s"></div>
            		<div class="items-center text-sm font-bold text-textColor absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-20"> Node %s </div>
        		</div>
			`

			activityTemplate = strings.Replace(activityTemplate, "NODE_IP", config.IPv4, -1)

			status := babylon.NodeStatus()
			color := ""
			adjective := ""

			if status.Failure {
				color = "bg-red-500"
				adjective = "Offline"
			} else {
				color = "bg-green-500"
				adjective = "Online"
			}

			tippyContent := fmt.Sprintf("<b>Node</b> - %s <br> <b>Owner</b> - %s<br> <b>Private IPv4</b> - %s <br> <b>Public IPv4</b> - %s <br> <b>Public IPv6</b> - %s <br>", config.Node, config.Owner, config.PrivateIPv4, config.IPv4, config.IPv6)

			needBr := ""

			if !status.Failure {
				tippyContent = tippyContent + fmt.Sprintf("<b>Current Height</b> - %d", status.Height)
				needBr = "<br>"
			}

			tippyContent = tippyContent + fmt.Sprintf("%s <b>Dashboard Version</b> - %s", needBr, config.NodeboxDashboardVersion)
			activityTemplate = strings.Replace(activityTemplate, "ALPINE_TOOLTIP", fmt.Sprintf(`tooltip-data="%s"`, tippyContent), -1)

			return c.SendString(fmt.Sprintf(activityTemplate, color, adjective))
		} else {
			return c.SendString("")
		}
	}
}
