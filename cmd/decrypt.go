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

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt values using key management system",
	Long: `This command will decrypt values that were encrypted using encrypt command. 
	
It requires access to the same key management system (KMS) that was used for encryption.

For example:

  cat encrypted-app-config.hcl | ehcl decrypt > app-config.hcl

`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		message, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal("failed to read:", err)
		}

		result, err := secrets.Decrypt(message)
		if err != nil {
			log.Fatal("failed to Decrypt:", err)
		}

		fmt.Println(string(result))
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)
}
