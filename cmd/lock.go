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
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofrs/flock"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	running    bool = false
	autoCreate bool = false
	dirToLock  []string
	locks      []*flock.Flock // an empty list
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock a file to avoid operations on it",
	Long:  `This commands locks a file on a file system until CTRL-C or process is killed. (DEBUG)`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("This command requires at least one file argument")
		}
		dirToLock = args
		log.WithFields(log.Fields{
			"dir to lock": dirToLock,
		}).Info("Files to lock")
		// TODO Mettre un mode verbose
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Catch signals to cleanup before exiting
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM|syscall.SIGKILL)
		/*
			go func() {
				<-c
				Shutdown()
			}()
		*/
		for i := 0; i < len(dirToLock); i++ {

			// First check that file/dir exists
			if _, err := os.Stat(dirToLock[i]); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					log.WithFields(log.Fields{
						"target": dirToLock[i],
					}).Fatal("File/directory does not exist")
					return
				} else {
					// Schrodinger: file may or may not exist. See err for details.
					// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
					log.Fatal(err)
				}
			}

			lockFilePath := dirToLock[i] + string(os.PathSeparator) + "misctools.lock"
			log.WithFields(log.Fields{
				"dir":      dirToLock[i],
				"lockFile": lockFilePath,
			}).Info("trying to acquire lock")
			fileLock := flock.NewFlock(lockFilePath)
			locked, err := fileLock.TryLock()
			if err != nil {
				// handle locking error
				log.WithFields(log.Fields{
					"file/dir": dirToLock[i],
				}).Fatal("Cannot lock file/directory : ", err)
			}

			if locked {
				locks = append(locks, fileLock)
				// do work
				//fileLock.Unlock()
			}
		}

		running = true
		done := make(chan bool, 1)
		go func() {
			log.Info("Press CTRL+C to stop the process...")
			sig := <-sigs
			log.WithFields(log.Fields{
				"signal": sig,
			}).Info("Signal caught")
			done <- true
		}()
		<-done
		log.Info("exiting...")
		for i := 0; i < len(locks); i++ {
			locks[i].Unlock()
		}
		/*
			fileLock := flock.NewFlock("/var/lock/go-lock.lock")

			locked, err := fileLock.TryLock()

			if err != nil {
				// handle locking error
				log.Fatal(err)
			}

			if locked {
				// do work
				fileLock.Unlock()
			}
		*/
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
	// viper.SetDefault("lock.file", "/var/lock/lockdev/misctools.lock")
	lockCmd.Flags().BoolP("debug", "d", false, "Help message for toggle")
}

func Shutdown() {
	if running {
		log.Println("filelock server is shutting down...")
		wg.Done()
		running = false
		os.Exit(0)
	} else {
		log.Println("filelock server was already stopped...")
	}
}
