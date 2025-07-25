package pwdb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/flags"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/pwdb"

	"github.com/spf13/cobra"
)

func documentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "document",
	}

	cmd.AddCommand(documentGetCommand())
	cmd.AddCommand(documentCreateCommand())

	return cmd
}

const documentGetFlagOut = "out"

func documentGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get <doc-id>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
				return fmt.Errorf("could not parse pwdb doc id as int: %w", err)
			}

			c := pwdb.New(
				pwdb.WithEndpoint(endpoint),
				pwdb.WithAPIKey(apiKey),
			)

			doc, err := c.GetDocument(cmd.Context(), id)
			if err != nil {
				return err
			}

			out, err := cmd.Flags().GetString(documentGetFlagOut)
			if err != nil {
				return fmt.Errorf("could not get out flag")
			}

			if out == "-" {
				if _, err := io.Copy(os.Stdout, bytes.NewReader(doc)); err != nil {
					return fmt.Errorf("could not write document to stdout: %w", err)
				}
			} else {
				if err := os.WriteFile(out, doc, 0644); err != nil {
					return fmt.Errorf("could not write file: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringP(documentGetFlagOut, "o", "", "Output file. Use '-' to output to stdout")
	if err := cmd.MarkFlagRequired(documentGetFlagOut); err != nil {
		panic(err)
	}

	return cmd
}

const (
	documentCreateFlagFile = "file"
	documentCreateFlagName = "name"
	documentCreateFlagDesc = "description"
)

func documentCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create <password-id>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
				return fmt.Errorf("could not parse pwdb pw id as int: %w", err)
			}

			c := pwdb.New(
				pwdb.WithEndpoint(endpoint),
				pwdb.WithAPIKey(apiKey),
			)

			file, err := cmd.Flags().GetString(documentCreateFlagFile)
			if err != nil {
				return fmt.Errorf("could not read file flag: %w", err)
			}
			name, err := cmd.Flags().GetString(documentCreateFlagName)
			if err != nil {
				return fmt.Errorf("could not read name flag: %w", err)
			}
			desc, err := cmd.Flags().GetString(documentCreateFlagDesc)
			if err != nil {
				return fmt.Errorf("could not read description flag: %w", err)
			}

			if name == "" {
				if file == "-" {
					return errors.New("document name is required when reading from stdin")
				}
				name = filepath.Base(file)
			}
			var bs []byte
			if file == "-" {
				bs, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("could not read input form stdin: %w", err)
				}
			} else {
				bs, err = os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("could not read input file: %w", err)
				}
			}

			_, err = c.CreateDocument(cmd.Context(), id, name, desc, bs)
			return err
		},
	}

	cmd.Flags().StringP(documentCreateFlagFile, "f", "", "Input file. Use '-' to read from stdin")
	if err := cmd.MarkFlagRequired(documentCreateFlagFile); err != nil {
		panic(err)
	}
	cmd.Flags().StringP(documentCreateFlagName, "n", "", "Document name. Will use base file name as default")
	cmd.Flags().StringP(documentCreateFlagDesc, "d", "", "Document description")
	if err := cmd.MarkFlagRequired(documentCreateFlagDesc); err != nil {
		panic(err)
	}

	return cmd

}
