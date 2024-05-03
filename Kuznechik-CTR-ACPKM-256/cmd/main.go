package main

import (
	"Kuznechik-CTR-ACPKM-256/pkg/app"
	"log"
	"os"
)

func main() {
	// Password = "Pa$$w0rd"
	if len(os.Args) == 2 {
		if os.Args[1] != "97c94ebe5d767a353b77f3c0ce2d429741f2e8c99473c3c150e2faa3d14c9da6" {
			log.Fatalln("Password isn't corrected!")
		} else {
			log.Println("Password is corrected!")
			a := app.NewApp()
			a.Run()
		}
	} else {
		if len(os.Args) == 3 {
			if os.Args[1] == "-t" && os.Args[2] != "97c94ebe5d767a353b77f3c0ce2d429741f2e8c99473c3c150e2faa3d14c9da6" {
				log.Fatalln("Password isn't corrected!")
			} else {
				log.Println("Password is corrected!")
				a := app.NewApp()
				a.Test()
			}
		}
		log.Fatalln("Password wasn't send to program!")
	}
}
