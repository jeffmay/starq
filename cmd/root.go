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
	"os"
	"starq/internal/starq"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new instance of the base command without any subcommands
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "starq [<config-file>.yaml...] [-r <jq-rule> ...] [-a <jq-rule> ...]",
		Short: "starq transforms a document from stdin using jq-style rules.",
		Long:  `starq (pronounced "star-q") transforms a JSON or YAML document from stdin using jq-style rules and outputs to stdout.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.ParseFlags(args)
			if err != nil {
				return err
			}

			var opts starq.Opts
			opts.ConfigFiles = cmd.Flags().Args()
			opts.PrependRules, err = cmd.Flags().GetStringArray(flagPrependRule)
			if err != nil {
				return err
			}
			opts.AppendRules, err = cmd.Flags().GetStringArray(flagAppendRule)
			if err != nil {
				return err
			}

			if opts.IsEmpty() {
				return cmd.Help()
			}

			runner := starq.NewRunner(cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr())
			return runner.RunAllTransformers(&opts)
		},
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = NewRootCmd()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	InitRootCmd(rootCmd)
}

func InitRootCmd(rootCmd *cobra.Command) {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.Flags().StringArrayP(flagPrependRule, "r", []string{}, "prepends this rule to apply to the input before applying the config rules")
	rootCmd.Flags().StringArrayP(flagAppendRule, "a", []string{}, "appends this rule to apply after all the custom rules and config rules")
}

// Avoid typos by defining the flag names as constants
const (
	flagPrependRule = "rule"
	flagAppendRule  = "append-rule"
)
