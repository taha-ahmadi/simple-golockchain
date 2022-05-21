package main

import (
	"flag"
	"fmt"
	golockchain "github.com/taha-ahmadi/simple-golockchain"
	"os"
)

var allCommand []Command

type Command struct {
	Name        string
	Description string

	Run func(golockchain.Store, ...string) error
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

	fmt.Fprintf(flag.CommandLine.Output(), "Sub commands:\n")

	for i := range allCommand {
		fmt.Fprintf(flag.CommandLine.Output(), "  %s: %s\n", allCommand[i].Name, allCommand[i].Description)
	}

}

func dispatch(store golockchain.Store, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("at least one argument should be assign")
	}

	sub := args[0]
	for i := range allCommand {
		if sub == allCommand[i].Name {
			return allCommand[i].Run(store, args...)
		}
	}
	flag.Usage()
	return fmt.Errorf("invalid comand")
}

func addCommand(name, description string, run func(golockchain.Store, ...string) error) error {
	allCommand = append(allCommand, Command{
		Name:        name,
		Description: description,
		Run:         run,
	})
	return nil
}
