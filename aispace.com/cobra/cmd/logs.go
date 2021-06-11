/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "logs container",
	Long: `日志长信息`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("执行logs命令")
		//rootCmd 全局的kubeconfig 命令选项
		kubeconfig, err :=cmd.Flags().GetString("kubeconfig")
		fmt.Println(kubeconfig)
		//logsCmd 子命令的namespace 命令选项
		namespace,err := cmd.Flags().GetString("namespace")
		if err != nil{
			fmt.Println("获取命令错误")
			return
		}
		fmt.Printf("命名空间: %s\n",namespace)
		//logsCmd 子命令的其他传入参数这个参数是deployment
		fmt.Println(args)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	logsCmd.Flags().StringP("namespace","n","default","helo message for namespace")
}
