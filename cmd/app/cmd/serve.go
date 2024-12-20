package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/config"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server and listen for oncoming requests",
	Run: func(cmd *cobra.Command, args []string) {
		configEnv, err := cmd.Flags().GetString("env")
		if err != nil {
			panic(err.Error())
		}

		configFileName := fmt.Sprintf("%s.%s", config.DefaultConfigName, configEnv)
		cfg := config.NewConfig(configFileName, config.DefaultConfigLocation)
		s := server.NewWithConfig(cfg)

		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer s.Shutdown()
		err = s.Start()
		if err != nil {
			panic(err.Error())
		}

		<-exit
		fmt.Println("Stopping server...")
	},
}
