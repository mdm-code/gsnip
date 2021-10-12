package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mdm-code/gsnip/manager"
	"github.com/mdm-code/gsnip/parsing"
	"github.com/mdm-code/gsnip/snippets"
)

const Src = "/usr/local/share/gsnip/snippets"

var (
	fCmd  *flag.FlagSet
	dbCmd *flag.FlagSet
	src   string
	dbnm  string
	mode  string
	user  string
	pswd  string
	port  string
	host  string
)

// Parse command-line subcommands
func parseFlags(args []string) error {
	fCmd = flag.NewFlagSet("file", flag.ExitOnError)
	dbCmd = flag.NewFlagSet("db", flag.ExitOnError)

	fCmd.StringVar(&src, "src", Src, "source name")
	dbCmd.StringVar(&dbnm, "dbnm", dbnm, "db name")
	dbCmd.StringVar(&user, "user", user, "db user name")
	dbCmd.StringVar(&pswd, "pswd", pswd, "db password")
	dbCmd.StringVar(&port, "port", port, "db port")
	dbCmd.StringVar(&host, "host", host, "db host")

	if len(os.Args) < 2 {
		return fmt.Errorf("file or db subcommand is required")
	}

	var err error
	switch os.Args[1] {
	case "file":
		err = fCmd.Parse(os.Args[2:])
	case "db":
		err = dbCmd.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
	if err != nil {
		return err
	}
	return nil
}

// Check if there's anything to read on STDIN
func isPiped() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func main() {
	var snpts snippets.Container

	err := parseFlags(os.Args)
	if err != nil {
		fmt.Fprint(os.Stderr, "failed to parse command line arguments")
		os.Exit(1)
	}

	if dbCmd.Parsed() {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host,
			port,
			user,
			pswd,
			dbnm,
		)
		snpts, err = snippets.NewSnippetsDB("postgres", dsn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to set up snippet container due to %s", err)
			os.Exit(1)
		}
	} else if fCmd.Parsed() {
		var f *os.File
		f, err := os.Open(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open %s\n", src)
			os.Exit(1)
			f, err = os.Open(Src)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not find %s\n", Src)
				os.Exit(1)
			}
		}
		defer f.Close()
		parser := parsing.NewParser()
		snpts, err = parser.Parse(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		}
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}

	mgr, ok := manager.NewManager(snpts)
	if !ok {
		fmt.Fprint(os.Stderr, "failed to initialized snippet manager")
		os.Exit(1)
	}

	if isPiped() {
		s := bufio.NewScanner(os.Stdin)
		s.Split(bufio.ScanWords)
		var params []string
		for s.Scan() {
			params = append(params, s.Text())
		}

		output, err := mgr.Execute(params...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute: %s", err)
			os.Exit(1)
		}
		os.Stdout.WriteString(output)
	}
}
