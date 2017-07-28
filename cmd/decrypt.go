package cmd

import (
	"fmt"
	"log"

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
  ehcl decrypt encrypted-app-config.hcl > app-config.hcl
  ehcl decrypt -i app-config.hcl
`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := getURL(args)
		if err != nil {
			log.Fatal("failed to get url: ", err)
		}

		message, err := read(url)
		if err != nil {
			log.Fatal("failed to read:", err)
		}

		result, err := secrets.Decrypt(message)
		if err != nil {
			log.Fatal("failed to decrypt: ", err)
		}

		if isFileURL(url) && inplace {
			if err := write(url, result); err != nil {
				log.Fatal("failed to write:", err)
			}
		} else {
			fmt.Println(string(result))
		}
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().BoolVarP(&inplace, "inplace", "i", false, "Decrypt file in-place")
}
