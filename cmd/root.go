/*
Copyright © 2024 Jeff May <jeff.n.may@gmail.com>

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
	"os"
	"starq/internal/starq"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "starq [-r <jq-rule> ...] [-c <config-file>.yaml] [-a <jq-rule> ...]",
	Short: "starq transforms a document from stdin using jq-style rules.",
	Long:  `starq (pronounced "star-q") transforms a JSON or YAML document from stdin using jq-style rules and outputs to stdout.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.ParseFlags(args)
		exitIf(err)

		var opts starq.Opts
		opts.ConfigFile, err = cmd.Flags().GetString(flagConfigFile)
		exitIf(err)
		opts.PrependRules, err = cmd.Flags().GetStringArray(flagPrependRule)
		exitIf(err)
		opts.AppendRules, err = cmd.Flags().GetStringArray(flagAppendRule)
		exitIf(err)

		opts.Input = cmd.InOrStdin()
		opts.Output = cmd.OutOrStdout()

		err = starq.Run(opts)
		exitIf(err)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func exitIf(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.Flags().StringP(flagConfigFile, "c", "", "uses a YAML config file to define the rules")
	rootCmd.Flags().StringArrayP(flagPrependRule, "r", []string{}, "prepends this rule to apply to the input before applying the config rules")
	rootCmd.Flags().StringArrayP(flagAppendRule, "a", []string{}, "appends this rule to apply after all the custom rules and config rules")
}

// Avoid typos by defining the flag names as constants
const (
	flagConfigFile  = "config"
	flagPrependRule = "rule"
	flagAppendRule  = "append-rule"
)
