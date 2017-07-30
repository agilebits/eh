package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "eh",
	Short: "Encrypt and decrypt protected values in .hcl files",
	Long: `	
Protect secrets in .hcl files. 
	
This utility relies on the key management system (KMS) provided by the server environment.
For example, Amazon Web Services KMS is used for servers running on EC2 virtual 
machines. A special "local" implementation can be used when running on developer machines 
to avoid storing unencrypted secrets in version control systems.

The utility uses a special 'eh' element in the .hcl file that defines encryptin parameters.

For example:

	eh encrypt -i config.hcl
	eh read config.hcl
	eh decrypt -i config.hcl
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
