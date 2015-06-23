package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/rlmcpherson/s3gof3r"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	//"net/url"
	//"image/color"
	"path"
)

func registerHandlers(r *mux.Router) {
	//r.HandleFunc("/passthrough/{object:[0-9a-z/_.-]+}", PassthroughHandler)
	r.HandleFunc("/{object:[0-9A-Za-z/_.-]+}", ThumbnailHandler)
	//r.HandleFunc("/{object:[0-9a-z/-_.]+}/{previewType:[a-zA-Z0-9_-]+}", ThumbnailHandler)
	r.HandleFunc("/{object}/debug", DebugHandler)

}

func PassthroughHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("PassthroughHandler called")
	object := mux.Vars(r)["object"]
	fmt.Printf("GET /%s\n", object)

	k, err := s3gof3r.EnvKeys() // get S3 keys from environment
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	// Open bucket to put file into
	s3 := s3gof3r.New("", k)
	b := s3.Bucket(configuration.Asset_Bucket)

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
	fmt.Println(rw, "Hello", configuration.Asset_Bucket, object)
}

func DebugHandler(rw http.ResponseWriter, r *http.Request) {
	object := mux.Vars(r)["object"]

	basepath := "http://localhost:8097"
	newPath := path.Join(basepath, configuration.Preview_Bucket, object, "thumbnail")
	fmt.Fprintf(rw, "%s", newPath)
}

func ThumbnailHandler(rw http.ResponseWriter, r *http.Request) {
	object := mux.Vars(r)["object"]
	previewType := r.FormValue("t")

	if previewType == "" {
		http.Error(rw, "previewType empty. Add \"/passthrough/\" in front of your URL to see the original", 400)
		//fmt.Fprintf(rw, "previewType empty. Add \"/passthrough/\" in front of your URL to see the original")
		return
	}

	typeOptions, ok := configuration.Previews[previewType]

	if ok != true {
		http.Error(rw, "previewType not configured", 400)
		return
	}

	k, err := s3gof3r.EnvKeys() // get S3 keys from environment
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	// Create s3 path for thumbnail
	s3path := path.Join(configuration.Preview_Prefix, previewType, object)
	s3url := fmt.Sprintf("http://%s.%s/%s", configuration.Preview_Bucket, configuration.StorageDomain, path.Join(configuration.Preview_Prefix, previewType, object))

	// Ask Mr. Bloom
	exists := previewBloom.TestString(s3path)

	if exists == true {
		// Thumbnail exists, redirect and return
		//fmt.Fprintf(rw, "Redirecting to: %s", s3url)

		http.Redirect(rw, r, s3url, 301)
		return
	}

	// Fetch image and generate stuff
	// - TODO: Check if we're too busy to create the thumbnail now. Return a temporary error 502 if we are.
	// - TODO: Add the object into preview queue for later processing

	// Open bucket to put file into
	s3 := s3gof3r.New("", k)
	b := s3.Bucket(configuration.Asset_Bucket)

	rb, _, err := b.GetReader(path.Join(configuration.Asset_Prefix, object), nil)

	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	orgImg, err := imaging.Decode(rb)
	log.Printf("GET %s", path.Join(configuration.Asset_Prefix, object))

	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	dstImg, err := Preview(&orgImg, typeOptions)

	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	// Put into Preview Bucket
	pb := s3.Bucket(configuration.Preview_Bucket)

	hdr := make(http.Header)
	hdr.Add("Content-Type", "image/jpg")

	prw, err := pb.PutWriter(s3path, hdr, nil)

	if err = imaging.Encode(prw, dstImg, imaging.JPEG); err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	if err = prw.Close(); err != nil {
		fmt.Fprintf(rw, err.Error())
		return
	}

	// Output to browser
	err = imaging.Encode(rw, dstImg, imaging.JPEG)
	if err != nil {
		fmt.Fprintf(rw, err.Error())
		return

	}
}
