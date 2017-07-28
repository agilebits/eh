package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ehcl",
	Short: "Encrypt and decrypt protected values in .hcl files",
	Long: `	
This utility is used to protect secrets in the server apps. 
	
It relies on the key management system (KMS) provided by the server environment.
For example, Amazon Web Services KMS is used for servers running on EC2 virtual 
machines. A fake KMS implementation can be used when running on developer machines 
to avoid storing unencrypted secrets in version control systems.

For example:

	cat plaintext.config.hcl | ehcl encrypt > encrypted.config.hcl
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var inplace bool

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ehcl/config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")      // name of config file (without extension)
	viper.AddConfigPath("$HOME/.ehcl") // adding home directory as first search path
	viper.AutomaticEnv()               // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
