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
	"github.com/bkyoung/maxwell/internal/systemd"
	"os"

	"github.com/spf13/cobra"
)

// uninstallServiceCmd represents the uninstallService command
var uninstallServiceCmd = &cobra.Command{
	Use:   "uninstall",
	Short: fmt.Sprintf("Uninstall %s as a systemd unit", executableName),
	Long: fmt.Sprintf("Uninstall %s as a systemd unit", executableName),
	Run: func(cmd *cobra.Command, args []string) {
		unitPath := fmt.Sprintf("/etc/systemd/system/%s.service", executableName)
		agent := systemd.New(executableName,
			systemd.WithUnitPath(unitPath),
		)
		err := systemd.Uninstall(agent);if err != nil {
			fmt.Printf("Failed to uninstall systemd unit: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully uninstalled %s.service\n", executableName)
	},
}

func init() {
	rootCmd.AddCommand(uninstallServiceCmd)
}
