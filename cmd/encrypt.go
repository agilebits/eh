package cmd

import (
	"fmt"
	"log"

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
  echl encrypt app-config.hcl > encrypted-app-config.hcl
  echl encrypt -i app-config.hcl
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

		result, err := secrets.Encrypt(message)
		if err != nil {
			log.Fatal("failed to encrypt:", err)
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
	RootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().BoolVarP(&inplace, "inplace", "i", false, "Encrypt file in-place")
}
