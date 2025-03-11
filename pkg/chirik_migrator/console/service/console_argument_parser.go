package service

import (
	"chickChirick/pkg/chirik_migrator/console/config"
	"chickChirick/pkg/chirik_migrator/file"
	"fmt"
	"os"
)

func ParseInput() error {
	if len(os.Args) <= 1 {
		return fmt.Errorf("not enough input arguments")
	}
	command := os.Args[1]
	if len(command) <= 0 {
		return fmt.Errorf("empty command: %s", command)
	}

	switch command {
	case config.MigrateKey:
		return file.ReadDir(os.Args[2:])
	default:
		return fmt.Errorf("invalid command")
	}
}
