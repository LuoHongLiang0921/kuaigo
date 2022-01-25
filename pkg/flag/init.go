// @Description

package flag

import (
	"flag"
	"os"
)

var defaultFlags = []Flag{
	// HelpFlag prints usage of application.
	&BoolFlag{
		Name:  "help",
		Usage: "--help, show help information",
		Action: func(name string, fs *FlagSet) {
			fs.PrintDefaults()
			os.Exit(0)
		},
	},
}

func init() {
	// procName := filepath.Base(os.Args[0])
	// nfs := flag.NewFlagSet(procName, flag.ExitOnError)
	flagset = &FlagSet{
		FlagSet:  flag.CommandLine,
		flags:    defaultFlags,
		actions:  make(map[string]func(string, *FlagSet)),
		environs: make(map[string]string),
	}
}