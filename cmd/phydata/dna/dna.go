// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package dna is a metapackage for commands
// that dealt with DNA sequences.
package dna

import (
	"github.com/js-arias/command"
	"github.com/js-arias/phydata/cmd/phydata/dna/add"
)

func init() {
	Command.Add(add.Command)
}

var Command = &command.Command{
	Usage: "dna <command> [<argument>...]",
	Short: "commands for DNA sequences",
}
