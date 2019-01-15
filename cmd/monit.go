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
	"log"
	"sync"

	"github.com/frdrolland/misctools/elk"
	"github.com/frdrolland/misctools/influx"
	"github.com/spf13/cobra"
)

// monitCmd represents the monit command
var wg sync.WaitGroup
var monitCmd = &cobra.Command{
	Use:   "monit",
	Short: "Start a monitoring server to simulate Elasticsearch and/or InfluxDB",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting elasticsearch mock...")
		influx.Startup()
		wg.Add(1)
		elk.Startup()
		wg.Add(1)
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(monitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//monitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
