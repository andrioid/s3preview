package main

import (
	"github.com/disintegration/imaging"
	"image"
	"io"
)

// Does Preview Exist? (objectname, type)
// Read Original (objectname)
// Create PreviewBase, original for image files (objectname) (image)
// Generate Preview (io.Reader, previewType) (*image, type)
// Upload Preview (path)
// Redirect to Preview (objectname, type)

// When do I need to buffer Image
// - Not when creating on the fly preview
// - When creating multiple preview files
// - When preview base has to be created from content

// Translates previewTypes to Imaging functions. Returns image interface
func Preview(r io.Reader, opt PreviewOptions) (img image.Image, err error) {

	oimg, err := imaging.Decode(r)
	if err != nil {
		return img, err
	}

	if opt.Method == "resize" {
		img = imaging.Resize(oimg, opt.Width, opt.Height, imaging.Box)
	} else {
		img = imaging.Thumbnail(oimg, opt.Width, opt.Height, imaging.Linear)
	}

	return // img and err are implicit
}

// Path of Preview file
// - Used to check if exists, uploading, redirecting
func PreviewPath(objname string, ptype string) (path string) {
	return ""
}

func Generate() {

}
