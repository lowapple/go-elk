package main

import (
	"flag"
	"fmt"
	"github.com/lowapple/go-elk/src/common/config"
	"github.com/lowapple/go-elk/src/notion"
	files "github.com/lowapple/go-elk/src/utils"
	"log"
	"os"
)

func init() {
	flag.StringVar(&config.OutputPath, "o", "./dist", "notion output static page dir path")
}

func download(URL string) {
	err := notion.New().Extract(URL)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: elk [args] URLs...")
		flag.PrintDefaults()
		os.Exit(1)
	}
	for _, a := range args {
		err := files.IsNotExistMkDir(config.OutputPath)
		if err != nil {
			log.Fatal(err)
			return
		}
		download(a)
	}
}
