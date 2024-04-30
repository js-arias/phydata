// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// PhyData is a tool for management of character data
// for phylogenetic analysis.
package main

import (
	"github.com/js-arias/command"
	"github.com/js-arias/phydata/cmd/phydata/obs"
)

var app = &command.Command{
	Usage: "phydata <command> [<argument>...]",
	Short: "a tool for phylogenetic data management",
}

func init() {
	app.Add(obs.Command)
}

func main() {
	app.Main()
}
