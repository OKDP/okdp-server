/*
 *    Copyright 2024 okdp.io
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package cmd

import (
	"os"

	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/constants"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "okdp-server",
	Short: "OKDP Server",
	Run:   runOkdpServer,
}

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("server.port", 8090)
	viper.SetConfigName("okdp-server")
	viper.SetConfigType("yaml")

	RootCmd.PersistentFlags().String("config", "config.yaml", "Path to configuration file")
	if err := viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config")); err != nil {
		panic("Unable to read server configuration: " + err.Error())
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}

func runOkdpServer(_ *cobra.Command, _ []string) {
	config := config.GetAppConfig()
	log.SetupGlobalLogger(config.Logging)

	server := server.NewOKDPServer(config)
	log.Info("ListenAddress %s: ", config.Server.ListenAddress)
	log.Info("Port %d: ", config.Server.Port)
	log.Info("okdp server started on port %d, requests api on %s", config.Server.Port, constants.OkdpServerBaseURL)
	log.Fatal(server.ListenAndServe())
}
