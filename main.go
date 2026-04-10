package main

import (
	"github.com/backsofangels/grimoire/cmd"
	_ "github.com/backsofangels/grimoire/internal/providers/android"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
