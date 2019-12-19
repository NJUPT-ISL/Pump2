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
	"github.com/Mr-Linus/Pump2/pkg/scheduler/server"
	"github.com/spf13/cobra"
	"os"
)

// schedulerCmd represents the scheduler command
var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Start Pump2 Scheduler",
	Long: `This command is used to start Pump2 Scheduler, which will launch a Gin Web 
server to get the image build tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.RunScheduler(cfgFile)
	},
}

func init() {
	runCmd.AddCommand(schedulerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// schedulerCmd.PersistentFlags().String("foo", "", "A help for foo")
	schedulerCmd.PersistentFlags().StringVarP(&cfgFile, "configFile", "f", os.Getenv("HOME")+"/pump2.yaml", "Pump2 config file (default is $HOME/pump2.yaml)")
	schedulerCmd.PersistentFlags().StringVarP(&serverIp, "IP", "i", "0.0.0.0", "IP address used by the Pump2 Builder to run")
	schedulerCmd.PersistentFlags().StringVarP(&serverPort, "Port", "p", "5020", "Port used by the Pump2 Builder to run")
	schedulerCmd.PersistentFlags().BoolVarP(&enableTLS, "EnableTLS", "t", false, "Enable the TLS authentication when running Pump2")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// schedulerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
