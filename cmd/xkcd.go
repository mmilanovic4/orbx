package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/spf13/cobra"
)

type XKCD struct {
	Title string `json:"title"`
	Img   string `json:"img"`
	Alt   string `json:"alt"`
	Num   int    `json:"num"`
}

var xkcdCmd = &cobra.Command{
	Use:     "xkcd",
	Short:   "Fetch latest XKCD comic",
	GroupID: "misc",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("https://xkcd.com/info.0.json")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		var comic XKCD
		if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
			fmt.Println("Decode error:", err)
			return
		}

		fmt.Printf("\n#%d - %s\n", comic.Num, comic.Title)
		fmt.Println("Image:", comic.Img)
		fmt.Println("\nAlt text:")
		fmt.Println(comic.Alt)

		exec.Command("open", comic.Img).Start()
	},
}

func init() {
	rootCmd.AddCommand(xkcdCmd)
}
