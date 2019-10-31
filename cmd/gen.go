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
	"os"

	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a profile template",
	Long: `This command is used to generate a configuration file template. 
You can execute the launcher gen command to generate the pump2.yaml file by default in the /etc/pump2 location.
If you need to modify the location of the generated file, you can use the -f option to modify the location where
the file was created, for example:
Pump2 gen -f /etc/pump2.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Create(cfgFile)
		if err != nil {
			fmt.Println("Generate a profile error:")
			fmt.Println(err)
			_ = f.Close()
			return
		}
		var exampleYaml = `pump2:
  serverip: 0.0.0.0
  serverport: 5020
  tls:
	tlskey: /etc/pump2/tls.key
	tlscrt: /etc/pump2/tls.crt
`
		_, err = f.WriteString(exampleYaml)
		if err != nil {
			fmt.Println(err)
			_ = f.Close()
			return
		}
		fmt.Println("Generate a profile:" + cfgFile)
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genCmd.PersistentFlags().StringVarP(&cfgFile, "file", "f", "/etc/pump2/pump2.yaml", "Path to generate the configuration file.")
}
