package cmd

import (
	"fmt"
	"orbx/internal/dateutil"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	tz string
	ms bool
)

var unixtsCmd = &cobra.Command{
	Use:     "unixts [to|from] [value]",
	Short:   "Unix timestamp utilities",
	GroupID: "dev",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println(time.Now().Unix())
			return nil
		}

		mode := args[0]

		switch mode {
		case "to":
			if len(args) < 2 {
				return fmt.Errorf("missing timestamp")
			}

			ts, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid timestamp: %w", err)
			}

			loc, err := dateutil.GetLocation(tz)
			if err != nil {
				return fmt.Errorf("invalid timezone %q: %w", tz, err)
			}

			t := time.Unix(ts, 0).In(loc)
			fmt.Println(t.Format(dateutil.GetLayout(ms)))

		case "from":
			if len(args) < 2 {
				return fmt.Errorf("missing datetime")
			}

			t, err := time.Parse(time.RFC3339, args[1])
			if err != nil {
				return fmt.Errorf("invalid format, use RFC3339 (e.g. 2026-04-19T21:57:11+02:00)")
			}

			if ms {
				fmt.Println(t.UnixMilli())
			} else {
				fmt.Println(t.Unix())
			}

		default:
			return fmt.Errorf("unknown mode %q — use: to, from", mode)
		}

		return nil
	},
}

func init() {
	unixtsCmd.Flags().StringVar(&tz, "tz", "", "timezone (e.g. Europe/Belgrade)")
	unixtsCmd.Flags().BoolVar(&ms, "ms", false, "use millisecond precision")
	rootCmd.AddCommand(unixtsCmd)
}
