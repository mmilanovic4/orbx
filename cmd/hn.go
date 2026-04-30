package cmd

import (
	"encoding/json"
	"fmt"
	"orbx/internal/netutil"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

type HNStory struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	By    string `json:"by"`
	Score int    `json:"score"`
	Time  int64  `json:"time"`
}

const (
	HN_BASE_URL = "https://hacker-news.firebaseio.com"
	HN_ITEM_URL = "https://news.ycombinator.com/item?id="
)

var (
	hnCount int
)

var hnCmd = &cobra.Command{
	Use:     "hn",
	Short:   "Fetch top stories from Hacker News",
	GroupID: "misc",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := netutil.Get(HN_BASE_URL + "/v0/topstories.json")
		if err != nil {
			return fmt.Errorf("failed to fetch top stories: %w", err)
		}

		var ids []int
		if err := json.Unmarshal(resp.Body, &ids); err != nil {
			return fmt.Errorf("failed to decode story list: %w", err)
		}

		if hnCount > len(ids) {
			hnCount = len(ids)
		}

		fmt.Printf("Hacker News — Top %d Stories\n\n", hnCount)

		for i, id := range ids[:hnCount] {
			storyResp, err := netutil.Get(HN_BASE_URL + "/v0/item/" + strconv.Itoa(id) + ".json")
			if err != nil {
				return fmt.Errorf("failed to fetch story %d: %w", id, err)
			}

			var story HNStory
			if err := json.Unmarshal(storyResp.Body, &story); err != nil {
				return fmt.Errorf("failed to decode story %d: %w", id, err)
			}

			fmt.Printf("%d. %s\n", i+1, story.Title)
			fmt.Printf("   %d points by %s\n", story.Score, story.By)
			fmt.Printf("   ID: %d\n", story.ID)

			if story.URL != "" {
				fmt.Printf("   %s\n", story.URL)
			}

			fmt.Println()
		}

		return nil
	},
}

var hnOpenCmd = &cobra.Command{
	Use:   "open [id]",
	Short: "Open a Hacker News story by ID in the browser",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			return fmt.Errorf("invalid story ID: %s", args[0])
		}

		storyResp, err := netutil.Get(HN_BASE_URL + "/v0/item/" + strconv.Itoa(id) + ".json")
		if err != nil {
			return fmt.Errorf("failed to fetch story: %w", err)
		}

		var story HNStory
		if err := json.Unmarshal(storyResp.Body, &story); err != nil {
			return fmt.Errorf("failed to decode story: %w", err)
		}

		url := story.URL
		if url == "" {
			url = HN_ITEM_URL + strconv.Itoa(story.ID)
		}

		fmt.Printf("Opening: %s\n", story.Title)
		fmt.Printf("URL: %s\n", url)

		if err := exec.Command("open", url).Start(); err != nil {
			return fmt.Errorf("failed to open browser: %w", err)
		}

		return nil
	},
}

func init() {
	hnCmd.Flags().IntVar(&hnCount, "count", 10, "number of stories to fetch")
	hnCmd.AddCommand(hnOpenCmd)
	rootCmd.AddCommand(hnCmd)
}
