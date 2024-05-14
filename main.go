package main

import (
	"log"

	"github.com/raefon/td-stream/commands"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	commands.Execute()
}
