package cmd

import (
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/airgap"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/flags"
	ntnx "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/ntnx"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/parse"
	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd/pwdb"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "iksctl",
	}

	cmd.PersistentFlags().String(flags.PWDBEndpoint, "https://pwdb.systematicgroup.local/api", "PWDB endpoint")
	cmd.PersistentFlags().String(flags.PWDBAPIKey, "", "PWDB api key (Envvar: IKS_PWDB_APIKEY)")

	cmd.AddCommand(airgap.Command())
	cmd.AddCommand(parse.Command())
	cmd.AddCommand(pwdb.Command())
	cmd.AddCommand(ntnx.NtnxCommand())

	return cmd
}
