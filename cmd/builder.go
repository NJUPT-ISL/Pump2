/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	ser "github.com/Mr-Linus/Pump2/pkg/builder"
	ya "github.com/Mr-Linus/Pump2/pkg/yaml"
	"github.com/spf13/cobra"
	"os"
)

// builderCmd represents the builder command
var builderCmd = &cobra.Command{
	Use:   "builder",
	Short: "Start Pump2 Builder",
	Long: `This command is used to start Pump2 Builder, which will launch a gRPC 
server using TLS authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		if enableTLS {
			conf, err := ya.ReadConfigYaml(cfgFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			ser.StartWithTLS(conf.Pump2.ServerIP,
				conf.Pump2.ServerPort,
				conf.Pump2.TLS.TLSCrt,
				conf.Pump2.TLS.TLSKey)
		} else {
			ser.StartWithoutTLS(serverIp, serverPort)
		}
	},
}

func init() {
	runCmd.AddCommand(builderCmd)

	// builderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	builderCmd.PersistentFlags().StringVarP(&cfgFile, "configFile", "f", os.Getenv("HOME")+"/pump2.yaml", "Pump2 config file (default is $HOME/pump2.yaml)")
	builderCmd.PersistentFlags().StringVarP(&serverIp, "IP", "i", "0.0.0.0", "IP address used by the Pump2 Builder to run")
	builderCmd.PersistentFlags().StringVarP(&serverPort, "Port", "p", "5020", "Port used by the Pump2 Builder to run")
	builderCmd.PersistentFlags().BoolVarP(&enableTLS, "EnableTLS", "t", false, "Enable the TLS authentication when running Pump2")
}
