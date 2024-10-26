/*
Copyright © 2024 Prudhvi Deep prudhvideep1996@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/prudhvideep/fnd/pkg/search"

	"github.com/prudhvideep/fnd/pkg/config"
	"github.com/spf13/cobra"
)

var profile *config.Credentials

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fnd",
	Short: "Find files or patterns locally or in remote/cloud environments",
	Long: `fnd is a versatile CLI tool designed to search for files or patterns both locally and in remote/cloud environments. 

It mimics the behavior of the traditional 'find' command but is optimized for use in modern cloud-based workflows. Whether you need to search through local directories or remote storage systems, fnd provides a simple and intuitive interface to make searching easy.

Usage examples:
  - fnd "*.txt"               // Search for all .txt files locally
  - fnd "error" --remote       // Search for the pattern 'error' in remote environments

With support for various environments, fnd ensures you can quickly find what you're looking for, no matter where it's located.`,
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		configFlag, err := cmd.Flags().GetBool("configure")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if configFlag {
			creds, err := config.GetSshCredentials()
			if err != nil {
				log.Fatal(err)
			}

			config.UpdateConfig(creds)

			return
		}

		remoteFlag, err := cmd.Flags().GetBool("remote")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// fmt.Println("Remote Flag ----> ", remoteFlag)

		typeFlag, err := cmd.Flags().GetString("type")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		dirFlag, err := cmd.Flags().GetString("dir")
		if err != nil {
			os.Exit(1)
		}
		if remoteFlag {
			search.RemoteSearch(args, typeFlag, dirFlag, profile)
		} else {
			search.Find(args, typeFlag, dirFlag)
		}

		fmt.Println(time.Since(t))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.CreateConfigFile()
	profile, _ = config.InitDefaultProfile()

	// fmt.Println("Default Profile ----> ", profile)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fnd.yaml)")
	// rootCmd.AddCommand(configCmd)
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("type", "t", "", "Use this flag to mention the type")
	rootCmd.PersistentFlags().StringP("dir", "d", "", "Use this flag to specify the directory")
	rootCmd.PersistentFlags().BoolP("remote", "r", false, "Use this flag to specify remote execution.")
	rootCmd.PersistentFlags().BoolP("configure", "c", false, "Configure the ssh credentials")
}
