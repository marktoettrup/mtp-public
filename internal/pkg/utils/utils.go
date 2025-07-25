package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func GetDefaultFilter(name string) (int, int, string) {
	page := 0
	limit := 50
	filter := fmt.Sprintf("tolower(name) eq '%s'", name)

	return page, limit, filter
}

func MarkRequiredFlags(cmd *cobra.Command, flags ...string) error {
	for _, flag := range flags {
		if err := cmd.MarkFlagRequired(flag); err != nil {
			return fmt.Errorf("failed to mark flag '%s' as required: %w", flag, err)
		}
	}
	return nil
}

func ParseMemoryToGB(memoryStr string) (int64, error) {
	memoryStr = strings.TrimSpace(memoryStr)
	if len(memoryStr) < 3 {
		return 0, fmt.Errorf("invalid memory format: %s", memoryStr)
	}

	unit := memoryStr[len(memoryStr)-2:]
	valueStr := memoryStr[:len(memoryStr)-2]

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory value: %s", valueStr)
	}

	switch unit {
	case "Gi":
		return value, nil
	case "Mi":
		return value / 1024, nil
	case "Ki":
		return value / (1024 * 1024), nil
	default:
		if valueWithoutSuffix, err := strconv.ParseInt(memoryStr, 10, 64); err == nil {
			return valueWithoutSuffix, nil
		}
		return 0, fmt.Errorf("unsupported memory unit: %s", unit)
	}
}

func ParseStorageToGB(storageStr string) (int64, error) {
	storageStr = strings.TrimSpace(storageStr)
	if len(storageStr) < 3 {
		return 0, fmt.Errorf("invalid storage format: %s", storageStr)
	}

	unit := storageStr[len(storageStr)-2:]
	valueStr := storageStr[:len(storageStr)-2]

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid storage value: %s", valueStr)
	}

	switch unit {
	case "Gi":
		return value, nil
	case "Mi":
		return value / 1024, nil
	case "Ki":
		return value / (1024 * 1024), nil
	default:
		if valueWithoutSuffix, err := strconv.ParseInt(storageStr, 10, 64); err == nil {
			return valueWithoutSuffix, nil
		}
		return 0, fmt.Errorf("unsupported storage unit: %s", unit)
	}
}
