package encodingutil

import (
	"fmt"
	"io"
	"os"

	"github.com/mmilanovic4/orbx/internal/sysutil"
)

func GetInputData(text string, file string) ([]byte, error) {
	// 1. File
	if file != "" {
		data, err := sysutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		return data, nil
	}

	// 2. Text
	if text != "" {
		return []byte(text), nil
	}

	// 3. STDIN
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin stat: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("failed to read stdin: %w", err)
		}
		return data, nil
	}

	// Nothing provided
	return nil, fmt.Errorf("no input provided")
}
