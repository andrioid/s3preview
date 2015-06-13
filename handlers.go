package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/rlmcpherson/s3gof3r"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	//"net/url"
	"image/color"
	"path"
)

func registerHandlers(r *mux.Router) {
	r.HandleFunc("/{object}", PassthroughHandler)
	r.HandleFunc("/{object}/debug", DebugHandler)
	r.HandleFunc("/{object}/{previewType}", ThumbnailHandler)
}

func PassthroughHandler(rw http.ResponseWriter, r *http.Request) {
	object := mux.Vars(r)["object"]

	k, err := s3gof3r.EnvKeys() // get S3 keys from environment
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	// Open bucket to put file into
	s3 := s3gof3r.New("", k)
	b := s3.Bucket(configuraton.Asset_Bucket)

	rb, h, err := b.GetReader(object, nil)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}
	// stream to standard output
	if _, err = io.Copy(rw, rb); err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}
	err = rb.Close()
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	fmt.Println(h) // print key header data
	fmt.Println(rw, "Hello", configuraton.Asset_Bucket, object)
}

func DebugHandler(rw http.ResponseWriter, r *http.Request) {
	object := mux.Vars(r)["object"]

	basepath := "http://localhost:8097"
	newPath := path.Join(basepath, configuraton.Preview_Bucket, object, "thumbnail")
	fmt.Fprintf(rw, "%s", newPath)
}

func ThumbnailHandler(rw http.ResponseWriter, r *http.Request) {
	object := mux.Vars(r)["object"]
	previewType := mux.Vars(r)["previewType"]

	typeOptions, ok := configuraton.Previews[previewType]

	if ok != true {
		fmt.Fprintf(rw, "previewType not found")
		return

	}

	k, err := s3gof3r.EnvKeys() // get S3 keys from environment
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	// Create s3 path for thumbnai
	s3path := path.Join(configuraton.Preview_Prefix, previewType, object)
	s3url := "http://andridk-assets.s3.amazonaws.com/" + path.Join(configuraton.Preview_Prefix, previewType, object)

	resp, err := http.Head(s3url)
	if err != nil {
		fmt.Fprintf(rw, err.Error())
		return
	}

	if resp.StatusCode == 200 {
		// Thumbnail exists, redirect and return
		//fmt.Fprintf(rw, "Redirecting to: %s", s3url)

		http.Redirect(rw, r, s3url, 301)
		return
	}

	// Fetch image and generate stuff

	// Open bucket to put file into
	s3 := s3gof3r.New("", k)
	b := s3.Bucket(configuraton.Asset_Bucket)

	rb, _, err := b.GetReader(object, nil)

	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	img, err := imaging.Decode(rb)
	if err != nil {
	}

	pb := s3.Bucket(configuraton.Preview_Bucket)

	hdr := make(http.Header)

	hdr.Set("Content-Type", "image/jpg")
	prw, err := pb.PutWriter(s3path, hdr, nil)
	if err != nil {
		fmt.Fprintf(rw, err.Error())
		return
	}

	dstImg := imaging.New(typeOptions.Width, typeOptions.Height, color.NRGBA{255, 0, 0, 255})
	if typeOptions.Method == "thumbnail" {
		dstImg = imaging.Thumbnail(img, typeOptions.Width, typeOptions.Height, imaging.Linear)
	} else if typeOptions.Method == "resize" {
		dstImg = imaging.Resize(img, typeOptions.Width, typeOptions.Height, imaging.Box)
	} else {
		fmt.Fprintf(rw, "Preview method '%s', not implemented.", typeOptions.Method)
		return
	}

	err = imaging.Encode(prw, dstImg, imaging.JPEG)
	if err != nil {
	}

	if err = prw.Close(); err != nil {
		fmt.Fprintf(rw, err.Error())
		return
	}

	err = imaging.Encode(rw, dstImg, imaging.JPEG)
	if err != nil {
		fmt.Fprintf(rw, err.Error())
		return

	}
}
