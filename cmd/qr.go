package cmd

import (
	"fmt"
	"orbx/internal/encodingutil"
	"path/filepath"
	"strings"

	go_qr "github.com/piglig/go-qr"
	"github.com/spf13/cobra"
)

var (
	qrFile  string
	qrLevel string
	qrOut   string
)

var qrCmd = &cobra.Command{
	Use:     "qr [text]",
	Short:   "Generate a QR code from text or URL",
	GroupID: "util",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		data, err := encodingutil.GetInputData(input, qrFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		text := strings.TrimSpace(string(data))
		if text == "" {
			return fmt.Errorf("input is empty")
		}

		ecc, err := parseQRLevel(qrLevel)
		if err != nil {
			return err
		}

		qr, err := go_qr.EncodeText(text, ecc)
		if err != nil {
			return fmt.Errorf("failed to encode QR code: %w", err)
		}

		if qrOut != "" {
			return saveQRToFile(qr, qrOut)
		}

		renderQRTerminal(qr)
		return nil
	},
}

func parseQRLevel(level string) (go_qr.Ecc, error) {
	switch level {
	case "1":
		return go_qr.Low, nil
	case "2":
		return go_qr.Medium, nil
	case "3":
		return go_qr.Quartile, nil
	case "4":
		return go_qr.High, nil
	default:
		return go_qr.Medium, fmt.Errorf("unknown level %q: use 1, 2, 3 or 4", level)
	}
}

func renderQRTerminal(qr *go_qr.QrCode) {
	size := qr.Size()
	for r := 0; r <= size+1; r += 2 {
		for c := 0; c <= size+1; c++ {
			x := c - 1
			topY := r - 1
			botY := r
			topDark := x >= 0 && x < size && topY >= 0 && topY < size && qr.Module(x, topY)
			botDark := x >= 0 && x < size && botY >= 0 && botY < size && qr.Module(x, botY)
			switch {
			case topDark && botDark:
				fmt.Print("█")
			case topDark:
				fmt.Print("▀")
			case botDark:
				fmt.Print("▄")
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func saveQRToFile(qr *go_qr.QrCode, path string) error {
	cfg := go_qr.NewQrCodeImgConfig(10, 4)
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		if err := qr.PNG(cfg, path); err != nil {
			return fmt.Errorf("failed to save PNG: %w", err)
		}
	case ".svg":
		if err := qr.SVG(cfg, path); err != nil {
			return fmt.Errorf("failed to save SVG: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file extension %q: use .png or .svg", ext)
	}
	return nil
}

func init() {
	qrCmd.Flags().StringVarP(&qrFile, "file", "f", "", "read input from file")
	qrCmd.Flags().StringVarP(&qrLevel, "level", "l", "medium", "error correction level (1, 2, 3, 4)")
	qrCmd.Flags().StringVarP(&qrOut, "out", "o", "", "save to file (.png or .svg)")
	rootCmd.AddCommand(qrCmd)
}
