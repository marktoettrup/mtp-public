package context

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-logr/logr"
)

type ctxKey string

const (
	ctxKeyTmpDir ctxKey = "tmp-dir"
	ctxKeyLogger ctxKey = "logger"
)

func WithLogger(ctx context.Context, logger *logr.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}

func GetLogger(ctx context.Context) (*logr.Logger, error) {
	l, ok := ctx.Value(ctxKeyLogger).(*logr.Logger)
	if !ok {
		return nil, errors.New("could not get logger from context")
	}
	return l, nil
}

func WithTmpDir(ctx context.Context) (context.Context, error) {
	tmpDir, err := os.MkdirTemp("", "airgap-helm-charts")
	if err != nil {
		return nil, fmt.Errorf("could not create tmp dir: %w", err)
	}
	return context.WithValue(ctx, ctxKeyTmpDir, tmpDir), nil
}

func GetTmpDir(ctx context.Context) (string, error) {
	s, ok := ctx.Value(ctxKeyTmpDir).(string)
	if !ok {
		return "", errors.New("could not get tmp dir from context")
	}
	return s, nil
}

func RmTmpDir(ctx context.Context) {
	tmpDir, err := GetTmpDir(ctx)
	if err != nil {
		return
	}
	os.RemoveAll(tmpDir)
}
