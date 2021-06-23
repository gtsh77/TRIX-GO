package main

import (
	"flag"
	"fmt"
	"github.com/gtsh77/TRIX-GO/q3parser"
	_ "github.com/gtsh77/TRIX-GO/mlib"
	"os"
)

func main() {

	argc := len(os.Args)

	if argc < 2 {
		fmt.Printf("usage: trix2 [-c </path/to/file>] || [-m <mode> [opt]]\n")
	} else {
		cc := flag.Bool("c", false, "compile flag")
		sm := flag.Bool("m", false, "mode flag")
		// fs := flag.Bool("f", false, "fullscreen flag")
		// fr := flag.Bool("b", false, "framerate flag")
		flag.Parse()

		if *cc {
			if argc < 3 {
				fmt.Printf("specify map path\n")
			} else {
				q3parser.ParseMap(os.Args[2])
			}
		} else if *sm {
			if argc < 3 {
				fmt.Printf("specify mode name\n")
			} else if os.Args[2] == "load" {
				if argc < 4 {
					fmt.Printf("specify level name\n")
				} else {
					//load level
				}
			} else {
				fmt.Printf("unknown mode\n")
			}
		}
	}

	//fmt.Printf(q3parser.ParseMap("test"))
	//fmt.Printf("%d",len(os.Args))
}
