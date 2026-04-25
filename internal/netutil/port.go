package netutil

import (
	"fmt"
	"strconv"
)

func ParsePort(input string) (int, error) {
	port, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q: must be a number", input)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port must be between 1 and 65535")
	}

	return port, nil
}
