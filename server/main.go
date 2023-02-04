package main

import (
	"fmt"

	"parrotDB/server/parrotserver"
)

var banner = "\n__________                                    __ ________ __________ \n\\______   \\_____ ______" +
	"________  ____________/  |\\______ \\\\______   \\\n |     ___/\\__  \\\\_  __ \\_  __ \\/  _ \\_  __ \\   __\\" +
	"    |  \\|    |  _/\n |    |     / __ \\|  | \\/|  | \\(  <_> )  | \\/|  | |    `   \\    |   \\\n |____|    (" +
	"____  /__|   |__|   \\____/|__|   |__|/_______  /______  /\n                \\/                                " +
	"       \\/       \\/ \n"

func main() {
	fmt.Printf(banner)
	s := parrotserver.NewServer()
	s.ListenAndServe()
}
