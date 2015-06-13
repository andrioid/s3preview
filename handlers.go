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
	"path"
)

func registerHandlers(r *mux.Router) {
	r.HandleFunc("/{bucket:[a-z0-9-]+}/{object}", PassthroughHandler)
	r.HandleFunc("/{bucket:[a-z0-9-]+}/{object}/debug", DebugHandler)
	r.HandleFunc("/{bucket:[a-z0-9-]+}/{object}/thumbnail", ThumbnailHandler)

}

func PassthroughHandler(rw http.ResponseWriter, r *http.Request) {
	bucket := mux.Vars(r)["bucket"]
	object := mux.Vars(r)["object"]

	k, err := s3gof3r.EnvKeys() // get S3 keys from environment
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	// Open bucket to put file into
	s3 := s3gof3r.New("", k)
	b := s3.Bucket(bucket)

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
	fmt.Println(rw, "Hello", bucket, object)
}

func DebugHandler(rw http.ResponseWriter, r *http.Request) {
	bucket := mux.Vars(r)["bucket"]
	object := mux.Vars(r)["object"]

	basepath := "http://localhost:8097"
	newPath := path.Join(basepath, bucket, object, "thumbnail")
	fmt.Fprintf(rw, "%s", newPath)
}

func ThumbnailHandler(rw http.ResponseWriter, r *http.Request) {
	bucket := mux.Vars(r)["bucket"]
	object := mux.Vars(r)["object"]

	k, err := s3gof3r.EnvKeys() // get S3 keys from environment
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	// Create s3 path for thumbnai
	s3path := path.Join(preview_prefix, bucket, "thumbnail", object)
	s3url := "http://andridk-assets.s3.amazonaws.com/" + path.Join(preview_prefix, bucket, "thumbnail", object)

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
	b := s3.Bucket(bucket)

	rb, h, err := b.GetReader(object, nil)
	fmt.Println(h)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	img, err := imaging.Decode(rb)
	if err != nil {
	}

	pb := s3.Bucket(preview_bucket)

	hdr := make(http.Header)

	hdr.Set("Content-Type", "image/jpg")
	prw, err := pb.PutWriter(s3path, hdr, nil)
	if err != nil {
		fmt.Fprintf(rw, err.Error())
		return
	}

	dstImg := imaging.Thumbnail(img, 100, 100, imaging.Linear)

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

	//fmt.Fprintln(rw, "bucket", preview_bucket, "preview bucket path", s3path)
	//fmt.Fprintln(rw, "preview bucket url", s3url)
	// Create url for thumbnail

	// Check if the thumbnail exists
	// Create Thumbnail
	// Upload Thumbnail

	//imaging.Encode(os., dstImg, imaging.PNG)
	//imaging.Save(dstImg, "thumbnail.jpg")

	//fmt.Fprintln(rw, "insert thumbnail here")
}
