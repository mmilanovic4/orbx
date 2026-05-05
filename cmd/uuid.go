package cmd

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	uuidCompact   bool
	uuidUpper     bool
	uuidVersion   int
	uuidName      string
	uuidNamespace string
)

var namespaces = map[string]uuid.UUID{
	"dns":  uuid.NameSpaceDNS,
	"url":  uuid.NameSpaceURL,
	"oid":  uuid.NameSpaceOID,
	"x500": uuid.NameSpaceX500,
}

func resolveNamespacedArgs(namespace, name string) (uuid.UUID, string, error) {
	if name == "" {
		return uuid.UUID{}, "", fmt.Errorf("--name is required for v3 and v5")
	}
	if namespace == "" {
		return uuid.UUID{}, "", fmt.Errorf("--namespace is required for v3 and v5 (dns, url, oid, x500)")
	}
	ns, ok := namespaces[strings.ToLower(namespace)]
	if !ok {
		return uuid.UUID{}, "", fmt.Errorf("unknown namespace %q, use: dns, url, oid, x500", namespace)
	}
	return ns, name, nil
}

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate a UUID (v1, v3, v4, v5, v6, v7)",
	Long: `Generate a UUID.

Supported versions:
  v1  MAC address + timestamp (privacy risk)
  v3  MD5 hash of namespace + name (deterministic)
  v4  Random (default)
  v5  SHA-1 hash of namespace + name (deterministic)
  v6  Sortable v1 reordering
  v7  Timestamp + random, sortable (good for database keys)

Namespaces for v3/v5: dns, url, oid, x500`,
	GroupID: "dev",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			id  string
			err error
			u   uuid.UUID
		)

		switch uuidVersion {
		case 1:
			u, err = uuid.NewUUID()
		case 3:
			ns, name, e := resolveNamespacedArgs(uuidNamespace, uuidName)
			if e != nil {
				return e
			}
			u = uuid.NewMD5(ns, []byte(name))
		case 4:
			u, err = uuid.NewRandom()
		case 5:
			ns, name, e := resolveNamespacedArgs(uuidNamespace, uuidName)
			if e != nil {
				return e
			}
			u = uuid.NewSHA1(ns, []byte(name))
		case 6:
			u, err = uuid.NewV6()
		case 7:
			u, err = uuid.NewV7()
		default:
			return fmt.Errorf("unsupported UUID version: %d (supported: 1, 3, 4, 5, 6, 7)", uuidVersion)
		}

		if err != nil {
			return fmt.Errorf("failed to generate UUID: %w", err)
		}

		id = u.String()

		if uuidCompact {
			id = strings.ReplaceAll(id, "-", "")
		}
		if uuidUpper {
			id = strings.ToUpper(id)
		}

		fmt.Println(id)
		return nil
	},
}

func init() {
	uuidCmd.Flags().BoolVar(&uuidCompact, "compact", false, "remove dashes from UUID")
	uuidCmd.Flags().BoolVar(&uuidUpper, "upper", false, "uppercase UUID")
	uuidCmd.Flags().IntVar(&uuidVersion, "version", 4, "UUID version (1, 3, 4, 5, 6, 7)")
	uuidCmd.Flags().StringVar(&uuidName, "name", "", "name for v3/v5")
	uuidCmd.Flags().StringVar(&uuidNamespace, "namespace", "", "namespace for v3/v5 (dns, url, oid, x500)")
	rootCmd.AddCommand(uuidCmd)
}
