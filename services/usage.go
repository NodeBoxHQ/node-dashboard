package services

import (
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	uptime2 "github.com/mackerelio/go-osstat/uptime"
	"math"
	"time"
)

func GetCPUUsage(c *fiber.Ctx) error {
	cpuUsageTemplate := `
                <div class="flex items-center space-x-2">
                    <svg width="1em" height="1em" viewBox="0 0 16 16" class="bi bi-cpu-fill" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                        <path fill-rule="evenodd" d="M5.5.5a.5.5 0 0 0-1 0V2A2.5 2.5 0 0 0 2 4.5H.5a.5.5 0 0 0 0 1H2v1H.5a.5.5 0 0 0 0 1H2v1H.5a.5.5 0 0 0 0 1H2v1H.5a.5.5 0 0 0 0 1H2A2.5 2.5 0 0 0 4.5 14v1.5a.5.5 0 0 0 1 0V14h1v1.5a.5.5 0 0 0 1 0V14h1v1.5a.5.5 0 0 0 1 0V14h1v1.5a.5.5 0 0 0 1 0V14a2.5 2.5 0 0 0 2.5-2.5h1.5a.5.5 0 0 0 0-1H14v-1h1.5a.5.5 0 0 0 0-1H14v-1h1.5a.5.5 0 0 0 0-1H14v-1h1.5a.5.5 0 0 0 0-1H14A2.5 2.5 0 0 0 11.5 2V.5a.5.5 0 0 0-1 0V2h-1V.5a.5.5 0 0 0-1 0V2h-1V.5a.5.5 0 0 0-1 0V2h-1V.5zm1 4.5A1.5 1.5 0 0 0 5 6.5v3A1.5 1.5 0 0 0 6.5 11h3A1.5 1.5 0 0 0 11 9.5v-3A1.5 1.5 0 0 0 9.5 5h-3zm0 1a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5h-3z"/>
                    </svg>
                    <div class="text-lg font-semibold">CPU</div>
                </div>
                <div class="mt-2 text-3xl">%v<span class="text-sm ml-1">[%%]</span></div>
                <div class="w-full bg-blue-200 rounded-full h-2.5 dark:bg-blue-700 mt-3">
                    <div class="bg-red-600 h-2.5 rounded-full" style="width: %v%%"></div>
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
	   <div class="flex items-center space-x-2">
           <svg width="1em" height="1em" viewBox="0 0 640 512" fill="currentColor" xmlns="http://www.w3.org/2000/svg"><path d="m640 130.94v-34.94c0-17.67-14.33-32-32-32h-576c-17.67 0-32 14.33-32 32v34.94c18.6 6.61 32 24.19 32 45.06s-13.4 38.45-32 45.06v98.94h640v-98.94c-18.6-6.61-32-24.19-32-45.06s13.4-38.45 32-45.06zm-416 125.06h-64v-128h64zm128 0h-64v-128h64zm128 0h-64v-128h64zm-480 192h64v-26.67c0-8.84 7.16-16 16-16s16 7.16 16 16v26.67h128v-26.67c0-8.84 7.16-16 16-16s16 7.16 16 16v26.67h128v-26.67c0-8.84 7.16-16 16-16s16 7.16 16 16v26.67h128v-26.67c0-8.84 7.16-16 16-16s16 7.16 16 16v26.67h64v-96h-640z"/></svg>
           <div class="text-lg font-semibold">Memory</div>
       </div>
       <div class="mt-2 text-3xl">%v<span class="text-sm ml-1">[%%]</span></div>
       <div class="w-full bg-blue-200 rounded-full h-2.5 dark:bg-blue-700 mt-3">
           <div class="bg-red-600 h-2.5 rounded-full" style="width: %v%%"></div>
       </div>
	`

	memoryStat, err := memory.Get()

	if err != nil {
		return c.SendString(fmt.Sprintf(ramUsageTemplate, 0.0, 0.0))
	}

	usage := 100 * (float64(memoryStat.Used) / float64(memoryStat.Total))
	return c.SendString(fmt.Sprintf(ramUsageTemplate, math.Round(usage*100)/100, math.Round(usage*100)/100))
}

func GetSystemUptime(c *fiber.Ctx) error {
	uptimeTemplate := `
		<div class="flex items-center space-x-2">
                        <svg width="1em" height="1em" fill="currentColor" viewBox="0 0 8 8" xmlns="http://www.w3.org/2000/svg"><path d="m4 0c-2.2 0-4 1.8-4 4s1.8 4 4 4 4-1.8 4-4-1.8-4-4-4zm0 1c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm-.5 1v2.22l.16.13.5.5.34.38.72-.72-.38-.34-.34-.34v-1.81h-1z"/></svg>
                        <div class="text-lg font-semibold">Uptime</div>
                    </div>
                    <div class="mt-2 text-md">%s</div>
	`

	uptime, err := uptime2.Get()

	if err != nil {
		return c.SendString(fmt.Sprintf(uptimeTemplate, "0s"))
	}

	return c.SendString(fmt.Sprintf(uptimeTemplate, utils.SecondsToReadable(int(uptime.Seconds()))))
}
