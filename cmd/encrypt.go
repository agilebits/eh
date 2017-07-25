package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/agilebits/ehcl/secrets"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt protected values using key management system",
	Long: ` 
	
Encrypt command is used to encrypt the protected values in the contents of 
the standard input and write result into the standard output. 

The .hcl file must include the 'ehcl' section.

For example:

  cat app-config.hcl | echl encrypt > encrypted-app-config.hcl
`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		message, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal("failed to read:", err)
		}

		result, err := secrets.Encrypt(message)
		if err != nil {
			log.Fatal("failed to encrypt:", err)
		}

		fmt.Println(string(result))
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringVarP(&env, "env", "e", "dev", "Environment type: 'dev' or 'aws")
	encryptCmd.Flags().StringVarP(&region, "region", "r", "", "AWS Region ('us-east-1')")
	encryptCmd.Flags().StringVarP(&masterKeyID, "master", "m", "", "Master key identifier")
}
