package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func init() {
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: elk [args] URLs...")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cmd := exec.Command("cd", "loconotion")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println(cmd.Run())

	for _, a := range args {
		fmt.Println(a)
	}
}
