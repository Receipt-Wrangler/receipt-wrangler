package services

import (
	"fmt"
	"math"

	"gopkg.in/gographics/imagick.v3/imagick"
)

// The copies of receipt images embedded in the claim packet's PDF. They are for
// reading, not evidence: the evidence is the untouched original in images/. So
// each is decoded small, oriented, flattened onto white and re-encoded as a
// modest JPEG, which keeps the HTML Chromium has to swallow small on a 1 GB host.
const (
	// packetImageMaxSide is the longest side, in pixels, of an embedded image (or
	// of one slice of a tall one). About 200 DPI across a Letter page's width.
	packetImageMaxSide = 1600
	// packetImageQuality is the JPEG quality of an embedded copy.
	packetImageQuality = 80
	// packetSliceAspect is the height/width of one embedded slice, close to the
	// printable area of a Letter page, so a slice fills a page without shrinking.
	packetSliceAspect = 1.3
	// packetMaxSlicesPerImage caps how many page-sized slices one image becomes.
	// A taller image gets taller slices, each shrunk to fit its page.
	packetMaxSlicesPerImage = 12
)

// packetImageLayout is how one image is embedded: resized to Width x Height and,
// when it is much taller than a page, cut into slices SliceHeight tall.
type packetImageLayout struct {
	Width       uint
	Height      uint
	SliceHeight uint
}

func (layout packetImageLayout) sliceCount() int {
	return int(math.Ceil(float64(layout.Height) / float64(layout.SliceHeight)))
}

// packetLayoutFor decides the embedded size of a width x height image. Pure, so
// the arithmetic is testable without ImageMagick.
//
// An image no taller than a page (with some slack) is shrunk to fit
// packetImageMaxSide on its longest side and embedded whole. A taller one — a
// rasterized multi-page PDF, which BuildDisplayImageString's pipeline stacks into
// one tall JPEG, or a long till receipt — would be unreadable shrunk onto one
// page, so it keeps a readable width and is cut into page-sized slices instead.
// Images are never enlarged.
func packetLayoutFor(width, height uint) packetImageLayout {
	if width == 0 || height == 0 {
		return packetImageLayout{Width: width, Height: height, SliceHeight: max(height, 1)}
	}

	aspect := float64(height) / float64(width)
	if aspect <= packetSliceAspect*1.15 {
		scale := math.Min(1, float64(packetImageMaxSide)/float64(max(width, height)))
		w := max(uint(math.Round(float64(width)*scale)), 1)
		h := max(uint(math.Round(float64(height)*scale)), 1)
		return packetImageLayout{Width: w, Height: h, SliceHeight: h}
	}

	sliceWidthLimit := uint(math.Floor(float64(packetImageMaxSide) / packetSliceAspect))
	w := min(width, sliceWidthLimit)
	h := max(uint(math.Round(float64(height)*float64(w)/float64(width))), 1)
	sliceHeight := max(uint(math.Round(float64(w)*packetSliceAspect)), 1)
	if minimum := uint(math.Ceil(float64(h) / packetMaxSlicesPerImage)); sliceHeight < minimum {
		sliceHeight = minimum
	}
	return packetImageLayout{Width: w, Height: h, SliceHeight: sliceHeight}
}

// packetDisplaySlices turns displayable image bytes (already through
// GetBytesFromImageBytes, so a PDF is a JPEG by now) into the JPEG slices the
// packet embeds, top to bottom.
//
// Memory: the image is pinged for its size first, and a JPEG is then decoded at
// a reduced scale (the jpeg:size hint lets libjpeg skip most of the pixels), so a
// 12-megapixel phone photo or a stacked PDF never sits in memory at full size.
func packetDisplaySlices(displayable []byte) ([][]byte, error) {
	width, height, isJpeg, err := pingImage(displayable)
	if err != nil {
		return nil, err
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if isJpeg {
		// The EXIF orientation is applied after decoding, so the hint has to hold
		// for either orientation: ask for the larger of the two scales needed.
		upright := packetLayoutFor(width, height)
		turned := packetLayoutFor(height, width)
		scale := math.Max(float64(upright.Width)/float64(width), float64(turned.Width)/float64(height))
		hintWidth := uint(math.Ceil(float64(width) * scale))
		hintHeight := uint(math.Ceil(float64(height) * scale))
		if err := mw.SetOption("jpeg:size", fmt.Sprintf("%dx%d", hintWidth, hintHeight)); err != nil {
			return nil, err
		}
	}

	if err := mw.ReadImageBlob(displayable); err != nil {
		return nil, err
	}
	// Only the first frame of a multi-frame image (GIF, multi-page TIFF).
	mw.SetFirstIterator()
	frame := mw.GetImage()
	if frame == nil || !frame.IsVerified() {
		return nil, fmt.Errorf("the image has no frame to embed")
	}
	defer frame.Destroy()

	if err := prepareFrameForJpeg(frame); err != nil {
		return nil, err
	}

	layout := packetLayoutFor(frame.GetImageWidth(), frame.GetImageHeight())
	if layout.Width != frame.GetImageWidth() || layout.Height != frame.GetImageHeight() {
		if err := frame.ResizeImage(layout.Width, layout.Height, imagick.FILTER_LANCZOS); err != nil {
			return nil, err
		}
	}

	slices := make([][]byte, 0, layout.sliceCount())
	for top := uint(0); top < layout.Height; top += layout.SliceHeight {
		height := min(layout.SliceHeight, layout.Height-top)
		slice, err := encodePacketSlice(frame, layout.Width, height, top)
		if err != nil {
			return nil, err
		}
		slices = append(slices, slice)
	}
	return slices, nil
}

func pingImage(data []byte) (width uint, height uint, isJpeg bool, err error) {
	ping := imagick.NewMagickWand()
	defer ping.Destroy()

	if err := ping.PingImageBlob(data); err != nil {
		return 0, 0, false, err
	}
	ping.SetFirstIterator()
	return ping.GetImageWidth(), ping.GetImageHeight(), ping.GetImageFormat() == "JPEG", nil
}

// prepareFrameForJpeg applies the EXIF orientation (a phone photo is often stored
// sideways and relies on it, and stripping metadata would lose it), converts to
// sRGB (a CMYK JPEG would otherwise print as garbage) and flattens any alpha onto
// white (JPEG has none, and transparency would come out black).
func prepareFrameForJpeg(frame *imagick.MagickWand) error {
	if err := frame.AutoOrientImage(); err != nil {
		return err
	}
	if err := frame.TransformImageColorspace(imagick.COLORSPACE_SRGB); err != nil {
		return err
	}

	white := imagick.NewPixelWand()
	defer white.Destroy()
	if !white.SetColor("white") {
		return fmt.Errorf("could not set the background colour")
	}
	if err := frame.SetImageBackgroundColor(white); err != nil {
		return err
	}
	if frame.GetImageAlphaChannel() {
		if err := frame.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE); err != nil {
			return err
		}
	}
	return frame.SetImageDepth(8)
}

func encodePacketSlice(frame *imagick.MagickWand, width uint, height uint, top uint) ([]byte, error) {
	slice := frame.GetImageRegion(width, height, 0, int(top))
	if slice == nil || !slice.IsVerified() {
		return nil, fmt.Errorf("could not cut the image at %d", top)
	}
	defer slice.Destroy()

	if err := slice.StripImage(); err != nil {
		return nil, err
	}
	if err := slice.SetImageFormat("jpeg"); err != nil {
		return nil, err
	}
	if err := slice.SetImageCompressionQuality(packetImageQuality); err != nil {
		return nil, err
	}
	// GetImageBlob reports warnings as errors; a blob that came with one is usable.
	blob, err := slice.GetImageBlob()
	if err != nil && len(blob) == 0 {
		return nil, err
	}
	return blob, nil
}
