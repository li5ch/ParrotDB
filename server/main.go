package main

import "parrotDB/server/parrotserver"

func main() {
	s := parrotserver.NewServer()
	s.ListenAndServe()
}
