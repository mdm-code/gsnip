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

const Source = "/usr/local/share/gsnip/snippets"

type flags struct {
	dbMode   *flag.FlagSet
	fMode    *flag.FlagSet
	dbName   *string
	user     *string
	password *string
	port     *int
	host     *string
	fName    *string
	mode     string
}

// TODO: Split to two different commands
// NOTE: https://github.com/calmh/mole/blob/master/cmd/mole/main.go
func parseFlags() (flags, error) {
	a := flags{}
	a.dbMode = flag.NewFlagSet("db", flag.ExitOnError)
	a.fMode = flag.NewFlagSet("file", flag.ExitOnError)

	a.dbName = a.dbMode.String("name", "gsnipdb", "database name")
	a.user = a.dbMode.String("user", "michal", "database user name")
	a.password = a.dbMode.String("pass", "dummy-pass", "database password")
	a.port = a.dbMode.Int("port", 5432, "database postgresql port")
	a.host = a.dbMode.String("host", "localhost", "postgresql host")

	a.fName = a.fMode.String("name", Source, "source file with snippets")

	if len(os.Args) < 2 {
		a.mode = "file"
	} else {
		a.mode = os.Args[1]
	}

	var args []string
	if len(os.Args) >= 3 {
		args = os.Args[2:]
	}

	switch a.mode {
	case "db":
		err := a.dbMode.Parse(args)
		if err != nil {
			return a, err
		}
	case "file":
		err := a.fMode.Parse(args)
		if err != nil {
			return a, err
		}
	default:
		return a, fmt.Errorf("unknown mode")
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Snippet manager written in Go.\n")
		fmt.Fprintf(os.Stderr, "Usage of db:\n")
		a.dbMode.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage of file:\n")
		a.fMode.PrintDefaults()
	}
	return a, nil
}

// Check if there's anything to read on STDIN
func isPiped() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func main() {
	fn := flag.String("f", Source, "Snippets source file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Snippet manager written in Go.\n\nUsage:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var f *os.File
	f, err := os.Open(*fn)
	if err != nil {
		log.Fatal("Could not open " + *fn)
		f, err = os.Open(Source)
		if err != nil {
			log.Fatal("Could not find " + Source)
		}
	}
	defer f.Close()

	parser := parsing.NewParser()
	var snippets snippets.Container
	snippets, err = parser.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	mgr, ok := manager.NewManager(snippets)
	if !ok {
		log.Fatal("failed to initialized snippet manager")
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
			log.Fatal("failed to execute command: ", err)
		}
		os.Stdout.WriteString(output)
	}
}
