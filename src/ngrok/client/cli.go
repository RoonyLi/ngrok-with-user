package client

import (
	"flag"
	"fmt"
	"ngrok/version"
	"os"
)

const usage1 string = `Usage: %s [OPTIONS] <local port or address>
Options:
`

const usage2 string = `
Examples:
	ngrok 80
	ngrok -subdomain=example 8080
	ngrok -proto=tcp 22
	ngrok -hostname="example.com" -httpauth="user:password" 10.0.0.1


Advanced usage: ngrok [OPTIONS] <command> [command args] [...]
Commands:
	ngrok start [tunnel] [...]    Start tunnels by name from config file
	ngork start-all               Start all tunnels defined in config file
	ngrok list                    List tunnel names from config file
	ngrok help                    Print help
	ngrok version                 Print ngrok version

Examples:
	ngrok start www api blog pubsub
	ngrok -log=stdout -config=ngrok.yml start ssh
	ngrok start-all
	ngrok version

`

type Options struct {
	config    string
	logto     string
	loglevel  string
	user string
	password string
	command   string
	args      []string
}

func ParseArgs() (opts *Options, err error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage1, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, usage2)
	}

	config := flag.String(
		"config",
		"",
		"Path to ngrok configuration file. (default: $HOME/.ngrok)")

	logto := flag.String(
		"log",
		"none",
		"Write log messages to this file. 'stdout' and 'none' have special meanings")

	loglevel := flag.String(
		"log-level",
		"DEBUG",
		"The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")

	user := flag.String(
		"user",
		"",
		"username for identifying an bestngrok.top account")

	password := flag.String(
		"password",
		"",
		"password for identifying an bestngrok.top account")
	
	flag.Parse()

	opts = &Options{
		config:    *config,
		logto:     *logto,
		loglevel:  *loglevel,
		user:      *user,
		password:  *password,
		command:   flag.Arg(0),
	}

	switch opts.command {
	case "list":
		opts.args = flag.Args()[1:]
	case "start":
		opts.args = flag.Args()[1:]
	case "start-all":
		opts.args = flag.Args()[1:]
	case "version":
		fmt.Println(version.MajorMinor())
		os.Exit(0)
	case "help":
		flag.Usage()
		os.Exit(0)
	case "":
		err = fmt.Errorf("Error: Specify a local port to tunnel to, or " +
			"an ngrok command.\n\nExample: To expose port 80, run " +
			"'ngrok 80'")
		return

	default:
		if len(flag.Args()) > 1 {
			err = fmt.Errorf("You may only specify one port to tunnel to on the command line, got %d: %v",
				len(flag.Args()),
				flag.Args())
			return
		}

		opts.command = "default"
		opts.args = flag.Args()
	}

	return
}
