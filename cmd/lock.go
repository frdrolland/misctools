// Copyright © 2018 frdrolland@yahoo.fr
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/theckman/go-flock"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock a file to avoid operations on it",
	Long:  `This commands locks a file on a file system until CTRL-C or process is killed.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO Mettre un mode verbose
		fmt.Println("lock called")
		fileLock := flock.NewFlock("/var/lock/go-lock.lock")

		locked, err := fileLock.TryLock()

		if err != nil {
			// handle locking error
		}

		if locked {
			// do work
			fileLock.Unlock()
		}
	},
}

func init() {
	rootCmd.AddCommand(lockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
