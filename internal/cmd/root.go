package cmd

import (
	"flag"
	"fmt"
	"os"
)

const Version = "2.0.2"

func AsciiArt() {
	fmt.Println(" _   _           _     ______      __   __")
	fmt.Println("| \\ | |         | |    | ___ \\     \\ \\ / /")
	fmt.Println("|  \\| | ___   __| | ___| |_/ / ___  \\ V / ")
	fmt.Println("| . ` |/ _ \\ / _` |/ _ \\ ___ \\/ _ \\ /   \\ ")
	fmt.Println("| |\\  | (_) | (_| |  __/ |_/ / (_) / /^\\ \\")
	fmt.Println("\\_| \\_/\\___/ \\__,_|\\___\\____/ \\___/\\/   \\/")

	fmt.Printf("\t\t\t\t         v%s\n", Version)
}

func ParseFlags() string {
	configPath := flag.String("config", "./config.json", "path to config file")
	help := flag.Bool("help", false, "print help and exit")
	version := flag.Bool("version", false, "print version and exit")

	flag.Parse()

	if *version {
		println(Version)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	return *configPath
}
