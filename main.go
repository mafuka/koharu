package main

import (
	"fmt"
	"github.com/mafuka/koharu/core"
	"strings"
)

const LOGO = `
/\_/\  /\_/\  /\_/\  /\_/\  /\_/\  /\_/\  /\_/\  /\_/\  /\_/\  /\_/\ 
( o.o )( o.o )( o.o )( o.o )( o.o )( o.o )( o.o )( o.o )( o.o )( o.o )
 > ^ <  > ^ <  > ^ <  > ^ <  > ^ <  > ^ <  > ^ <  > ^ <  > ^ <  > ^ < 
 /\_/\    ██╗  ██╗ ██████╗ ██╗  ██╗ █████╗ ██████╗ ██╗   ██╗    /\_/\ 
( o.o )   ██║ ██╔╝██╔═══██╗██║  ██║██╔══██╗██╔══██╗██║   ██║   ( o.o )
 > ^ <    █████╔╝ ██║   ██║███████║███████║██████╔╝██║   ██║    > ^ < 
 /\_/\    ██╔═██╗ ██║   ██║██╔══██║██╔══██║██╔══██╗██║   ██║    /\_/\ 
( o.o )   ██║  ██╗╚██████╔╝██║  ██║██║  ██║██║  ██║╚██████╔╝   ( o.o )
 > ^ <    ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝     > ^ < 
 /\_/\  /\_/\  /\_/\  /\_/\  /\_/\  /\_/\  /\_/\  /\_/\  /\_/\  /\_/\ 
( o.o )( o.o )( o.o )( o.o )( o.o )( o.o )( o.o )( o.o )( o.o )( o.o )
 > ^ <  > ^ <  > ^ <  > ^ <  > ^ <  > ^ <  > ^ <  > ^ <  > ^ <  > ^ <      
`

const (
	projectName = "Koharu"
	version     = "2.3.3"
	commit      = "114514a"
	buildDate   = "1997"
	repoURL     = "https://github.com/mafuka/koharu"
)

const (
	cfgFile = "config.yml"
)

func printBrand() {
	fmt.Printf("[i] %s %s (%s %s)\n", projectName, version, commit, buildDate)
	fmt.Printf("[i] %s\n", repoURL)
	fmt.Printf("[i] This program is licensed under MIT.\n")
	// keeping line breaks within LOGO to maintain aesthetics
	fmt.Print(strings.TrimPrefix(LOGO, "\n"))
}

func main() {
	printBrand()

	bot := core.New(
		core.WithConfig(core.DefaultConfig()),
		core.WithPProf(),
	)
	//bot.Use()
	err := bot.Run()
	if err != nil {
		panic(err)
	}
}
