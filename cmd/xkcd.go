package cmd

import (
	"encoding/json"
	"fmt"
	"orbx/internal/netutil"
	"os/exec"

	"github.com/spf13/cobra"
)

type XKCD struct {
	Title string `json:"title"`
	Img   string `json:"img"`
	Alt   string `json:"alt"`
	Num   int    `json:"num"`
}

const (
	XKCD_BASE_URL = "https://xkcd.com"
)

var (
	xkcdOpen bool
)

var xkcdCmd = &cobra.Command{
	Use:     "xkcd",
	Short:   "Fetch latest XKCD comic",
	GroupID: "misc",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := netutil.Get(HN_BASE_URL + "/info.0.json")
		if err != nil {
			return fmt.Errorf("failed to fetch comic: %w", err)
		}

		var comic XKCD
		if err := json.Unmarshal(resp.Body, &comic); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		fmt.Printf("#%d - %s\n", comic.Num, comic.Title)
		fmt.Println("Image:", comic.Img)
		fmt.Println("Alt text:")
		fmt.Println(comic.Alt)

		if xkcdOpen {
			exec.Command("open", comic.Img).Start()
		}

		return nil
	},
}

func init() {
	xkcdCmd.Flags().BoolVar(&xkcdOpen, "open", false, "open in browser")
	rootCmd.AddCommand(xkcdCmd)
}
