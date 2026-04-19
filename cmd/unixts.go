package cmd

import (
	"fmt"
	"orbx/internal/timeutil"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var tz string
var ms bool

var unixtsCmd = &cobra.Command{
	Use:     "unixts [to|from] [value]",
	Short:   "Unix timestamp utilities",
	GroupID: "dev",
	Run: func(cmd *cobra.Command, args []string) {
		// DEFAULT: current unix timestamp
		if len(args) == 0 {
			fmt.Println(time.Now().Unix())
			return
		}

		mode := args[0]

		switch mode {
		// unix → human
		case "to":
			if len(args) < 2 {
				fmt.Println("missing timestamp")
				return
			}

			ts, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid timestamp")
				return
			}

			loc, err := timeutil.GetLocation(tz)
			if err != nil {
				fmt.Println("invalid timezone:", tz)
				return
			}

			t := time.Unix(ts, 0).In(loc)

			fmt.Println(t.Format(timeutil.GetLayout(ms)))
		// human → unix
		case "from":
			if len(args) < 2 {
				fmt.Println("missing datetime")
				return
			}

			input := args[1]

			t, err := time.Parse(time.RFC3339, input)
			if err != nil {
				fmt.Println("invalid format, use RFC3339 like:")
				fmt.Println("2026-04-19T21:57:11+02:00")
				return
			}

			if ms {
				fmt.Println(t.UnixMilli())
			} else {
				fmt.Println(t.Unix())
			}
		default:
			fmt.Println("use: unixts [to|from]")
		}
	},
}

func init() {
	unixtsCmd.Flags().StringVar(&tz, "tz", "", "timezone (e.g. Europe/Belgrade)")
	unixtsCmd.Flags().BoolVar(&ms, "ms", false, "use millisecond precision")
	rootCmd.AddCommand(unixtsCmd)
}
