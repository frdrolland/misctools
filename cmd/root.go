// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"strings"
	"sync"
	"time"

	//	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool = false
	wg      sync.WaitGroup
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "misctools",
	Short: "Miscellaneous tools for development",
	Long:  `Miscellaneous tools for development`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	var log = logrus.New()
	log.Out = os.Stdout

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true, // Seems like automatic color detection doesn't work on windows terminals
		FullTimestamp:   true,
		TimestampFormat: time.RFC1123Z,
	})
	logrus.SetLevel(logrus.InfoLevel)

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "misctools.yml", "config file (default is $HOME/misctools.yml)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config/misctools.yml", "config file (default is $HOME/misctools.yml)")
	//	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose mode")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("verbose", "v", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".misctools" (without extension).
		//viper.AddConfigPath(home)
		viper.AddConfigPath(home + "/.misctools")
		viper.AddConfigPath("./config")
		//		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.WithFields(logrus.Fields{
			"File": viper.ConfigFileUsed(),
		}).Debug("Config file read")
	} else {
		logrus.WithFields(logrus.Fields{
			"File": viper.ConfigFileUsed(),
		}).Debug("Could not read config file", err)
	}

	logLevel := viper.GetString("log.level")
	switch strings.ToUpper(logLevel) {
	case logrus.ErrorLevel.String():
		logrus.SetLevel(logrus.ErrorLevel)
	case logrus.WarnLevel.String():
		logrus.SetLevel(logrus.WarnLevel)
	case logrus.InfoLevel.String():
		logrus.SetLevel(logrus.InfoLevel)
	case logrus.DebugLevel.String():
		logrus.SetLevel(logrus.DebugLevel)
	case logrus.TraceLevel.String():
	case "ALL":
		logrus.SetLevel(logrus.TraceLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	//viper.SetDefault("verbose", false)
	//viper.Set("LogFile", LogFile)

	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.WithFields(logrus.Fields{
			"File": e.Name,
		}).Debug("Config file changed")
	})
	viper.WatchConfig()
}
