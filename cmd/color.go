package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type RGB struct {
	R, G, B uint8
}

type RGBA struct {
	R, G, B uint8
	A       float64
}

func parseHex(s string) (RGB, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return RGB{}, fmt.Errorf("invalid hex color: expected 6 characters after #")
	}

	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color: %w", err)
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color: %w", err)
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color: %w", err)
	}

	return RGB{uint8(r), uint8(g), uint8(b)}, nil
}

func parseHexA(s string) (RGBA, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 8 {
		return RGBA{}, fmt.Errorf("invalid hex color: expected 8 characters after #")
	}

	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return RGBA{}, fmt.Errorf("invalid hex color: %w", err)
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return RGBA{}, fmt.Errorf("invalid hex color: %w", err)
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return RGBA{}, fmt.Errorf("invalid hex color: %w", err)
	}
	a, err := strconv.ParseUint(s[6:8], 16, 8)
	if err != nil {
		return RGBA{}, fmt.Errorf("invalid hex color: %w", err)
	}

	return RGBA{uint8(r), uint8(g), uint8(b), float64(a) / 255}, nil
}

func rgbToHex(c RGB) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func rgbaToHex(c RGBA) string {
	a := uint8(c.A * 255)
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, a)
}

var colorCmd = &cobra.Command{
	Use:   "color [value]",
	Short: "Convert between color formats (HEX, RGB, RGBA)",
	Long: `Convert between color formats.

Usage:
  orbx color #ff6600
  orbx color #ff6600cc
  orbx color rgb 255 102 0
  orbx color rgba 255 102 0 0.8`,
	GroupID: "dev",
	Args:    cobra.RangeArgs(1, 5),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch strings.ToLower(args[0]) {
		case "rgb":
			if len(args) != 4 {
				return fmt.Errorf("rgb requires 3 values: rgb R G B")
			}
			r, err := strconv.ParseUint(args[1], 10, 8)
			if err != nil {
				return fmt.Errorf("invalid R value: %w", err)
			}
			g, err := strconv.ParseUint(args[2], 10, 8)
			if err != nil {
				return fmt.Errorf("invalid G value: %w", err)
			}
			b, err := strconv.ParseUint(args[3], 10, 8)
			if err != nil {
				return fmt.Errorf("invalid B value: %w", err)
			}
			c := RGB{uint8(r), uint8(g), uint8(b)}
			fmt.Printf("HEX  %s\n", rgbToHex(c))
			fmt.Printf("RGB  rgb(%d, %d, %d)\n", c.R, c.G, c.B)

		case "rgba":
			if len(args) != 5 {
				return fmt.Errorf("rgba requires 4 values: rgba R G B A")
			}
			r, err := strconv.ParseUint(args[1], 10, 8)
			if err != nil {
				return fmt.Errorf("invalid R value: %w", err)
			}
			g, err := strconv.ParseUint(args[2], 10, 8)
			if err != nil {
				return fmt.Errorf("invalid G value: %w", err)
			}
			b, err := strconv.ParseUint(args[3], 10, 8)
			if err != nil {
				return fmt.Errorf("invalid B value: %w", err)
			}
			a, err := strconv.ParseFloat(args[4], 64)
			if err != nil || a < 0 || a > 1 {
				return fmt.Errorf("invalid A value: must be between 0.0 and 1.0")
			}
			c := RGBA{uint8(r), uint8(g), uint8(b), a}
			fmt.Printf("HEX  %s\n", rgbaToHex(c))
			fmt.Printf("RGBA rgba(%d, %d, %d, %.2f)\n", c.R, c.G, c.B, c.A)

		default:
			s := args[0]
			if !strings.HasPrefix(s, "#") {
				return fmt.Errorf("unknown format: use #RRGGBB, #RRGGBBAA, rgb, or rgba")
			}
			hex := strings.TrimPrefix(s, "#")
			switch len(hex) {
			case 3:
				hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
				s = "#" + hex
				fallthrough
			case 6:
				c, err := parseHex(s)
				if err != nil {
					return err
				}
				fmt.Printf("HEX  %s\n", rgbToHex(c))
				fmt.Printf("RGB  rgb(%d, %d, %d)\n", c.R, c.G, c.B)
			case 8:
				c, err := parseHexA(s)
				if err != nil {
					return err
				}
				fmt.Printf("HEX  %s\n", rgbaToHex(c))
				fmt.Printf("RGBA rgba(%d, %d, %d, %.2f)\n", c.R, c.G, c.B, c.A)
			default:
				return fmt.Errorf("invalid hex color: expected #RRGGBB or #RRGGBBAA")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(colorCmd)
}
