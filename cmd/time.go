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
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	iso bool = false
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Tool to work with time/date/timestamps",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tog, _ := cmd.Flags().GetBool("datetime")
		if true == tog {
			transformFromIsoDateTime(args)
		} else {
			transformFromEpoch(args)
		}
	},
}

//
// Transform milliseconds since epoch to readable date/time
//
func transformFromEpoch(args []string) {
	if nil == args || len(args) == 0 {
		fmt.Printf("error: No timestamp given in parameter\n")
		return
	}
	for _, ts := range args {
		result, err := msToTime(ts)
		if nil != err {
			fmt.Printf("error: on %s: %s\n", ts, err)
		} else {
			fmt.Printf("%s ==> %s\n", ts, result)
		}
	}
}

//
// Transform milliseconds since epoch to readable date/time
//
func transformFromIsoDateTime(args []string) {
	if nil == args || len(args) == 0 {
		fmt.Printf("error: No RFC3339-compatible date/time given in parameter\n")
		return
	}
	for _, ts := range args {

		result, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			fmt.Printf("error: on %s: %s\n", ts, err)
		} else {
			fmt.Printf("%s ==> %d\n", ts, (result.UnixNano() / 1000000))
		}
	}
}

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func init() {
	rootCmd.AddCommand(timeCmd)

	timeCmd.Flags().BoolVarP(&iso, "datetime", "d", false, "transform ISO date/time to timestamps (milliseconds from epoch)")

	// Here you will define your flags and configuration settings.
	//rootCmd.PersistentFlags().BoolVar(&cfgFile, "config", "config/misctools.yml", "config file (default is $HOME/misctools.yml)")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
