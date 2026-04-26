package cmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"
	"orbx/internal/sysutil"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"

	"github.com/spf13/cobra"
)

var faviconOut string

var faviconSizes = []struct {
	filename string
	size     int
}{
	{"favicon-16x16.png", 16},
	{"favicon-32x32.png", 32},
	{"favicon-48x48.png", 48},
	{"apple-touch-icon.png", 180},
	{"apple-touch-icon-152x152.png", 152},
	{"apple-touch-icon-120x120.png", 120},
}

var faviconCmd = &cobra.Command{
	Use:     "favicon [input.png]",
	Short:   "Generate favicon.ico and Apple touch icons from a PNG",
	GroupID: "dev",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer f.Close()

		src, err := png.Decode(f)
		if err != nil {
			return fmt.Errorf("failed to decode PNG: %w", err)
		}

		bounds := src.Bounds()
		if bounds.Dx() != bounds.Dy() {
			return fmt.Errorf("input must be a square PNG (got %dx%d)", bounds.Dx(), bounds.Dy())
		}

		outDir := faviconOut
		if outDir == "" {
			outDir = "."
		}

		// individual PNG sizes
		for _, s := range faviconSizes {
			resized := resizeImage(src, s.size)
			var buf bytes.Buffer
			if err := png.Encode(&buf, resized); err != nil {
				return fmt.Errorf("failed to encode %s: %w", s.filename, err)
			}
			dest := filepath.Join(outDir, s.filename)
			if err := sysutil.WriteFile(dest, buf.Bytes()); err != nil {
				return fmt.Errorf("failed to write %s: %w", s.filename, err)
			}
			fmt.Printf("✓ %s\n", dest)
		}

		// favicon.ico (16, 32, 48)
		icoSizes := []int{16, 32, 48}
		var images []image.Image
		for _, size := range icoSizes {
			images = append(images, resizeImage(src, size))
		}

		var icoBuf bytes.Buffer
		if err := encodeICO(&icoBuf, images); err != nil {
			return fmt.Errorf("failed to encode favicon.ico: %w", err)
		}

		icoPath := filepath.Join(outDir, "favicon.ico")
		if err := sysutil.WriteFile(icoPath, icoBuf.Bytes()); err != nil {
			return fmt.Errorf("failed to write favicon.ico: %w", err)
		}
		fmt.Printf("✓ %s\n", icoPath)

		return nil
	},
}

func resizeImage(src image.Image, size int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

func encodeICO(w io.Writer, images []image.Image) error {
	count := len(images)

	// encode each image as PNG into a buffer
	var bufs [][]byte
	for _, img := range images {
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return err
		}
		bufs = append(bufs, buf.Bytes())
	}

	// ICONDIR header (6 bytes)
	binary.Write(w, binary.LittleEndian, uint16(0))     // reserved
	binary.Write(w, binary.LittleEndian, uint16(1))     // type: 1 = ICO
	binary.Write(w, binary.LittleEndian, uint16(count)) // image count

	// offset starts after ICONDIR (6) + all ICONDIRENTRYs (16 * count)
	offset := uint32(6 + 16*count)

	// ICONDIRENTRY per image (16 bytes each)
	for i, img := range images {
		bounds := img.Bounds()
		w.Write([]byte{byte(bounds.Dx()), byte(bounds.Dy())})      // width, height
		binary.Write(w, binary.LittleEndian, uint8(0))             // color count
		binary.Write(w, binary.LittleEndian, uint8(0))             // reserved
		binary.Write(w, binary.LittleEndian, uint16(1))            // color planes
		binary.Write(w, binary.LittleEndian, uint16(32))           // bits per pixel
		binary.Write(w, binary.LittleEndian, uint32(len(bufs[i]))) // size
		binary.Write(w, binary.LittleEndian, offset)               // offset
		offset += uint32(len(bufs[i]))
	}

	// image data
	for _, buf := range bufs {
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	faviconCmd.Flags().StringVarP(&faviconOut, "out", "o", "", "output directory (default: current directory)")
	rootCmd.AddCommand(faviconCmd)
}
