package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"

	"github.com/evan-buss/openbooks/cli"
	"github.com/evan-buss/openbooks/irc"
	"github.com/evan-buss/openbooks/server"
)

var logIRC bool
var cliMode bool
var userName string
var port string

// Retrieve command line arguments and set appropriate variables
func init() {
	user, _ := user.Current()
	userName = strings.Split(user.Name, " ")[0]
	flag.BoolVar(&logIRC, "log", false, "Save IRC logs to irc_log.txt")
	flag.BoolVar(&cliMode, "cli", false, "Launch OpenBooks in the terminal instead of the web UI")
	flag.StringVar(&userName, "name", userName, "Use a name that differs from your account name. One word only")
	flag.StringVar(&port, "port", "5228", "Set the local network port for browser mode")
}

// Determine what mode to run the application in (CLI or Web Server)
func main() {
	flag.Parse()

	// Username can be supplied via ARGS or found from the user's system name
	if strings.Contains(userName, " ") {
		// If there is a space split it and take the first word
		fmt.Println("Using first word of entered username")
		userName = strings.Split(userName, " ")[0]
	}

	conn := irc.New(userName, userName)
	conn.Logging = logIRC

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		conn.Disconnect()
		os.Exit(1)
	}()

	if cliMode {
		cli.Start(conn)
	} else {
		server.Start(conn, port)
	}
}
