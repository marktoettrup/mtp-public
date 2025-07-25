package main

import (
	"context"
	"fmt"
	"os"

	"bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/cmd"
	ctxhelpers "bitbucket.systematicgroup.local/itmiks/eks-anywhere-tooling/iksctl/pkg/context"

	"github.com/bombsimon/logrusr/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	logrusLogger := logrus.New()
	logger := logrusr.New(logrusLogger)
	ctx := ctxhelpers.WithLogger(context.Background(), &logger)

	err := cmd.Command().ExecuteContext(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
