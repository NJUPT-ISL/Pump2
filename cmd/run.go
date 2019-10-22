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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Pump2 server",
	Long: `Run the gRPC-based Pump2 server`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}


func init() {
	rootCmd.AddCommand(runCmd)
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/etc/pump2/pump2.yaml", "Pump2 config file (default is $HOME/pump2.yaml)")
	runCmd.PersistentFlags().StringVarP(&serverIp, "IP", "i", "0.0.0.0", "IP address used by the Pump2 server to run")
	runCmd.PersistentFlags().IntVarP(&serverPort, "Port", "p", 5020, "Port used by the Pump2 server to run")
	runCmd.PersistentFlags().BoolVarP(&enableTLS, "EnableTLS", "-t", false, "Enable the TLS authentication when running Pump2")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".Pump2" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".Pump2")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}