/*
Copyright © 2021 Brandon Young <bkyoung@gmail.com>

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
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var executableName string
var defaultCfgFile = fmt.Sprintf("/etc/%s/config.yaml", executableName)
var version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   executableName,
	Short: fmt.Sprintf("%s secrets agent", executableName),
	Long: fmt.Sprintf(`%s secrets agent for updating vital system configuration elements securely`, executableName),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Bind flags for the command -- this runs for every subcommand, too
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			fmt.Printf("error reading flags: %s\n", err)
			os.Exit(1)
		}

		viper.SetEnvPrefix(executableName)
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.AutomaticEnv()
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cfgFileUsage := fmt.Sprintf("config file (default is %s)", defaultCfgFile)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", cfgFileUsage)
	rootCmd.SetVersionTemplate("{{printf \"maxwell %s\\n\" .Version}}")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(path.Dir(defaultCfgFile))
		viper.SetConfigName(strings.Split(path.Base(defaultCfgFile), ".")[0])
		viper.SetConfigType(path.Ext(path.Base(defaultCfgFile)))
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}
