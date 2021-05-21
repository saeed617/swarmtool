package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/saeed617/swarmtool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string
	config  *swarmtool.Config

	rootCmd = &cobra.Command{
		Use:          "swarmtool",
		Short:        "Utility for backup and restore of swarm cluster state",
		SilenceUsage: true,
	}
)

func init() {
	cobra.OnInitialize(loadConfig)
	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is $HOME/.swarmtool.yaml)")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func loadConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		viper.SetConfigName(".swarmtool")
		viper.AddConfigPath(home)
	}

	viper.SetEnvPrefix("SWARMTOOL")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Ignore not found error to make config file optional
			cobra.CheckErr(err)
		}
	}
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	cobra.CheckErr(err)
}
