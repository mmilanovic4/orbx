package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mmilanovic4/orbx/internal/encodingutil"

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

// qrQuietZone is the empty margin around the code, in modules. The QR spec
// asks for 4, with less some scanners do not find the code among the text
// around it.
const qrQuietZone = 4

// renderQRTerminal draws two module rows per line with half blocks.
func renderQRTerminal(qr *go_qr.QrCode) {
	size := qr.Size()
	dark := func(x, y int) bool {
		return x >= 0 && x < size && y >= 0 && y < size && qr.Module(x, y)
	}

	for y := -qrQuietZone; y < size+qrQuietZone; y += 2 {
		var line strings.Builder
		for x := -qrQuietZone; x < size+qrQuietZone; x++ {
			top, bottom := dark(x, y), dark(x, y+1)
			switch {
			case top && bottom:
				line.WriteString("█")
			case top:
				line.WriteString("▀")
			case bottom:
				line.WriteString("▄")
			default:
				line.WriteString(" ")
			}
		}
		fmt.Println(line.String())
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
	qrCmd.Flags().StringVarP(&qrLevel, "level", "l", "2", "error correction level (1, 2, 3, 4)")
	qrCmd.Flags().StringVarP(&qrOut, "out", "o", "", "save to file (.png or .svg)")
	rootCmd.AddCommand(qrCmd)
}
