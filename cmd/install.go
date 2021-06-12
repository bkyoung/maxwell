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
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bkyoung/maxwell/internal/systemd"
)

//go:embed assets/unit.service.tpl
var unit []byte

// installServiceCmd installs the systemd service unit to disk
var installServiceCmd = &cobra.Command{
	Use:   "install",
	Short: fmt.Sprintf("Install %s as a systemd unit", executableName),
	Long: fmt.Sprintf(`Install %s as a systemd unit`, executableName),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		executablePath := viper.GetString("executable-path")
		if executablePath == "" {
			executablePath, err = os.Executable()
			if err != nil {
				fmt.Printf("Error determining path to currently running executable: %s\n", err)
				os.Exit(1)
			}
		}
		data := struct{
			Config 		string
			Executable  string
			Path		string
		}{
			Config: cfgFile,
			Executable: executableName,
			Path: executablePath,
		}

		unitPath := fmt.Sprintf("/etc/systemd/system/%s.service", executableName)

		unitContent := &bytes.Buffer{}
		t := template.Must(template.New("unit-tmpl").Parse(string(unit)))
		err = t.Execute(unitContent, data); if err != nil {
			fmt.Printf("Error rendering systemd unit file: %s\n", err)
			os.Exit(1)
		}

		agent := systemd.New(executableName,
			systemd.WithConfigPath(cfgFile),
			systemd.WithUnitPath(unitPath),
			systemd.WithUnitContent(unitContent.Bytes()),
			systemd.WithUnitDisabled(viper.GetBool("disable")),
			systemd.WithExecutablePath(executablePath),
		)

		if viper.GetBool("write-config") {
			err = viper.WriteConfig();if err != nil {
				fmt.Printf("Failed to write config file: %s\n", err)
				os.Exit(1)
			}
		}

		err = systemd.Install(agent);if err != nil {
			fmt.Printf("Failed to install systemd unit: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully installed %s.service\n", executableName)
	},
}

func init() {
	rootCmd.AddCommand(installServiceCmd)
	installServiceCmd.Flags().Bool("disable", false, "Disable systemd unit after installing")
	installServiceCmd.Flags().String("executable-path", "", "Path to ExecStart (default is current path of binary)")
	installServiceCmd.Flags().Bool("write-config", false, "Write configuration file to disk (default is false)")
}
