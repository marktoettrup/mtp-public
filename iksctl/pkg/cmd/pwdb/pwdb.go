package pwdb

import (
	"github.com/spf13/cobra"
)

const (
	flagPWDBID           = "id"
	flagPWDBBase64Decode = "base64-decode"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "pwdb",
	}

	cmd.AddCommand(passwordCommand())
	cmd.AddCommand(documentCommand())

	return cmd
}
