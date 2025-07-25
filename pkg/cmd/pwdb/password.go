package pwdb

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"

	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/flags"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/pwdb"

	"github.com/spf13/cobra"
)

func passwordCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "password",
		Aliases: []string{"pw"},
	}

	cmd.AddCommand(passwordGetCommand())

	return cmd
}

func passwordGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",
	}

	cmd.AddCommand(passwordGetAllCommand())
	cmd.AddCommand(passwordGetUsernameCommand())
	cmd.AddCommand(passwordGetPasswordCommand())
	cmd.AddCommand(passwordGetNotesCommand())

	return cmd
}

func passwordGetAllCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "all <pwdb-id>",
		Args: cobra.ExactArgs(1),
		RunE: createPasswordGetRunE(passwordFieldTypeAll),
	}
}

func passwordGetUsernameCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "username <pwdb-id>",
		Args: cobra.ExactArgs(1),
		RunE: createPasswordGetRunE(passwordFieldTypeUsername),
	}
}

func passwordGetPasswordCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "password <pwdb-id>",
		Args: cobra.ExactArgs(1),
		RunE: createPasswordGetRunE(passwordFieldTypePassword),
	}
}

func passwordGetNotesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "notes <pwdb-id>",
		Args: cobra.ExactArgs(1),
		RunE: createPasswordGetRunE(passwordFieldTypeNotes),
	}

	cmd.Flags().Bool(flagPWDBBase64Decode, false, "enable base64 decode")

	return cmd
}

type passwordFieldType int

const (
	passwordFieldTypeAll = iota
	passwordFieldTypeUsername
	passwordFieldTypePassword
	passwordFieldTypeNotes
)

func createPasswordGetRunE(t passwordFieldType) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		endpoint, err := cmd.Flags().GetString(flags.PWDBEndpoint)
		if err != nil {
			return fmt.Errorf("could not read endpoint flag: %w", err)
		}

		apiKey := os.Getenv("IKS_PWDB_APIKEY")
		if apiKey == "" {
			apiKey, err = cmd.Flags().GetString(flags.PWDBAPIKey)
			if err != nil {
				return fmt.Errorf("could not read api key flag: %w", err)
			}
			if apiKey == "" {
				return errors.New("api-key is required")
			}
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("could not parse pwdb id as int: %w", err)
		}

		c := pwdb.New(
			pwdb.WithEndpoint(endpoint),
			pwdb.WithAPIKey(apiKey),
		)
		res, err := c.GetPassword(cmd.Context(), id)
		if err != nil {
			return err
		}

		switch t {
		case passwordFieldTypeAll:
			fmt.Printf("Username: %s\nPassword: %s\n", res.Username, res.Password)
		case passwordFieldTypeUsername:
			fmt.Println(res.Username)
		case passwordFieldTypePassword:
			fmt.Println(res.Password)
		case passwordFieldTypeNotes:
			decode, err := cmd.Flags().GetBool(flagPWDBBase64Decode)
			if err != nil {
				return fmt.Errorf("could not get base64 decode flag: %w", err)

			}
			if decode {
				decoded, err := base64.StdEncoding.DecodeString(res.Notes)
				if err != nil {
					return fmt.Errorf("could not base64 decode: %w", err)
				}
				fmt.Println(decoded)
			} else {
				fmt.Println(res.Notes)
			}
		}
		return nil
	}
}
