/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"strconv"

	"github.com/andersryanc/playing-with-go/gen"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a random string",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := cmd.Flags().GetString("type")
		if err != nil {
			return fmt.Errorf("unable to get type flag: %v", err)
		}

		var chars string
		switch t {
		case "alpha":
			chars = gen.AlphaChars
		case "alphanumeric":
			chars = gen.AlphaNumericChars
		case "hex":
			chars = gen.HexChars
		case "numeric":
			chars = gen.NumericChars
		case "default":
			chars = gen.DefaultChars
		}

		if len(args) == 0 {
			args = append(args, "8")
		}

		for i := 0; i < len(args); i++ {
			num, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Printf("unable to convert arg (%q) to int: %v\n", args[i], err)
				continue
			}
			fmt.Println(gen.RandomString(chars, num))
		}

		return nil
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
	genCmd.Flags().IntP("length", "l", 8, "the length of the random string")
	genCmd.Flags().StringP("type", "t", "default", "the type of string to be generated (accepted: alpha, alphanumeric, hex, numeric)")
}
