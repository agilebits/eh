package cmd

import (
	"fmt"
	"log"

	"github.com/agilebits/eh/secrets"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read and decrypt file and all included fragments",
	Long: `In addition to decrypting the protected values this command will also read and merge all included files specified in the eh include. 
	
For example:

  eh read app-config.hcl
`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := getURL(args)
		if err != nil {
			log.Fatal("failed to get url: ", err)
		}

		result, err := secrets.Read(url)
		if err != nil {
			log.Fatal("failed to read:", err)
		}

		fmt.Println(string(result))
	},
}

func init() {
	RootCmd.AddCommand(readCmd)
}
