// Copyright © 2024 J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package obs is a metapackage for commands
// that dealt with specimen character observations.
package obs

import (
	"github.com/js-arias/command"
	"github.com/js-arias/phydata/cmd/phydata/obs/add"
	"github.com/js-arias/phydata/cmd/phydata/obs/chars"
)

func init() {
	Command.Add(add.Command)
	Command.Add(chars.Command)
}

var Command = &command.Command{
	Usage: "obs <command> [<argument>...]",
	Short: "commands for character observations",
}
