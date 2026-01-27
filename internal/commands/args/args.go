package args

import (
	"flag"
	"log"
	"os"
	"strings"
)

type Positional struct {
	Url        string
	MethodHttp string
}

type Optional struct {
	Help       bool
	Version    bool
	Goroutines int
}

type Args struct {
	Positional
	Optional
}

func Parse() *Args {
	var positional Positional
	var optional Optional

	flag.IntVar(&optional.Goroutines, "goroutines", 5, "...")
	flag.BoolVar(&optional.Help, "help", false, "...")
	flag.BoolVar(&optional.Version, "version", false, "...")

	flag.Parse()

	positional.MethodHttp = strings.ToUpper(flag.Arg(0))
	positional.Url = flag.Arg(1)

	switch {
	case optional.Help:
		flag.PrintDefaults()
		os.Exit(0)
	case optional.Version:
		log.Println("v1.0.0")
		os.Exit(0)
	case len(positional.Url) == 0:
		log.Println("invalid url")
		os.Exit(1)
	case len(positional.MethodHttp) <= 0:
		log.Println("invalid method http")
		os.Exit(1)
	}

	return &Args{
		Positional: positional,
		Optional:   optional,
	}
}
