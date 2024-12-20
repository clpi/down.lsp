package config

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	Config = cobra.Command{
		Use:     "config <command>",
		Aliases: []string{"cfg", "conf"},
		Long:    "config",
		Short:   "c",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("config")
		},
	}
)
