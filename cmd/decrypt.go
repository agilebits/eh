package cmd

import (
	"fmt"
	"log"

	"github.com/agilebits/eh/secrets"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt values in .hcl file",
	Long: `This command will decrypt values that were protected using encrypt command. 
	
For example:

  eh decrypt app-config.hcl
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
