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
	pu "github.com/Mr-Linus/Pump2/pkg/pump2"
	ser "github.com/Mr-Linus/Pump2/pkg/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Pump2 server",
	Long: `Run the gRPC-based Pump2 server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called,"+serverIp+serverPort)
		lis, err := net.Listen("tcp",serverIp+":"+serverPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pu.RegisterPump2Server(s, &ser.P2Server{})
		if err := s.Serve(lis); err != nil{
			log.Fatalf("failed to serve: %v", err)
		}
	},
}


func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.PersistentFlags().StringVarP(&serverIp, "IP", "i", "0.0.0.0", "IP address used by the Pump2 server to run")
	runCmd.PersistentFlags().StringVarP(&serverPort, "Port", "p", "5020", "Port used by the Pump2 server to run")
}
