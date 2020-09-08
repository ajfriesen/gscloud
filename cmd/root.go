package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gridscale/gscloud/render"
	"github.com/gridscale/gscloud/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	account    string
	rt         *runtime.Runtime
	jsonFlag   bool
	quietFlag  bool
	renderOpts render.Options
)

const (
	defaultAPIURL = "https://api.gridscale.io"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "gscloud",
	Short:             "the CLI for the gridscale cloud",
	Long:              `gscloud lets you manage objects on gridscale.io via command line. It provides a Docker-CLI comparable command line that allows you to create, manipulate, and remove objects on gridscale.io.`,
	DisableAutoGenTag: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Register following initializers only when we are not running tests.
	if !runtime.UnderTest() {
		cobra.OnInitialize(initConfig, initRuntime)
	}

	rootCmd.PersistentFlags().StringVar(&configFile, "config", runtime.ConfigPath(), fmt.Sprintf("Specify a configuration file"))
	rootCmd.PersistentFlags().StringVarP(&account, "account", "", "default", "Specify the account used")
	rootCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "j", false, "Print JSON to stdout instead of a table")
	rootCmd.PersistentFlags().BoolVarP(&renderOpts.NoHeader, "noheading", "", false, "Do not print column headings")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "Print only IDs of objects")
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Use default paths.
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(runtime.ConfigPath())
		viper.AddConfigPath(".")
	}
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Not found. Disregard
		} else if _, ok := err.(*os.PathError); ok && commandWithoutConfig(os.Args) {
			// --config given along with make-config → we're about to create that file. Disregard
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

// initRuntime initializes the client for a given account.
func initRuntime() {
	theRuntime, err := runtime.NewRuntime(account)
	if err != nil {
		log.Fatal(err)
	}
	rt = theRuntime
}

// commandWithoutConfig return true if current command does not need a config file.
// Called from within a cobra initializer function. Unfortunately there is no
// way of getting the current command from an cobra initializer so we scan the
// command line again.
func commandWithoutConfig(cmdLine []string) bool {
	var noConfigNeeded = []string{
		"make-config", "version", "manpage", "completion",
	}
	for _, cmd := range noConfigNeeded {
		if contains(cmdLine, cmd) {
			return true
		}
	}
	return false
}

// contains tests whether string e is in slice s.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
