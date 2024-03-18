package main

import (
	"fmt"
	"strings"

	"github.com/mafuka/koharu/bot"
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
	version     = "dev"
	commit      = "none"
	buildDate   = "unknown"
	repoURL     = "https://github.com/mafuka/koharu"
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

	cfg := bot.DefaultConfig()
	cfg.Admin = []int{2982280379}
	cfg.Log.File = "console"
	cfg.Log.Level = bot.DebugLevel
	cfg.Server.Port = 8081

	b := bot.New(cfg)

	err := b.Run()
	if err != nil {
		bot.Log().Fatalf("Bot exited abnormally: %+v", err)
	}
}
