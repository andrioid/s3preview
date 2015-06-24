package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/rlmcpherson/s3gof3r"
	"github.com/willf/bloom"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
)

var previewBloom *bloom.BloomFilter

func PopulatePreviewCache() {
	svc := s3.New(
		aws.Auth{
			AccessKey: configuration.AWS_Key,
			SecretKey: configuration.AWS_Secret},
		aws.USEast,
		nil)

	bucket := svc.Bucket(configuration.Preview_Bucket)

	resp, err := bucket.List(configuration.Preview_Prefix, "", "", 0)

	if err != nil {
		fmt.Println("prebucket list", err)
		return
	}

	bloomSize := uint(1)

	if len(resp.Contents) > 0 {
		bloomSize = uint(len(resp.Contents))
	}

	previewBloom = bloom.NewWithEstimates(bloomSize, 0.001)

	for _, elem := range resp.Contents {
		previewBloom.AddString(elem.Key)
	}

	fmt.Printf("%d existing previews added to cache.\n", len(resp.Contents))

}

type PreviewPart struct {
	Path  string
	Types []string
}

func GenerateMissing() {
	svc := s3.New(
		aws.Auth{
			AccessKey: configuration.AWS_Key,
			SecretKey: configuration.AWS_Secret},
		aws.USEast,
		nil)

	pbu := svc.Bucket(configuration.Asset_Bucket)

	resp, err := pbu.List(configuration.Asset_Prefix, "", "", 0)

	if err != nil {
		fmt.Printf("Error listing from Asset Bucket: '%s'\n", err)
	}

	missingPreviews := make(chan PreviewPart, 20)

	s3g := s3gof3r.New("", s3gof3r.Keys{configuration.AWS_Key, configuration.AWS_Secret, ""}) // I know, right
	assBucket := s3g.Bucket(configuration.Asset_Bucket)
	preBucket := s3g.Bucket(configuration.Preview_Bucket)

	go func() {
		for {
			missing := <-missingPreviews

			// Fetch Original
			assRead, _, err := assBucket.GetReader(missing.Path, nil)
			if err != nil {
				fmt.Printf("asset error (%s): %s\n", missing.Path, err.Error())
				continue
			}

			//TODO: Use content-type to handle other types
			//fmt.Println("http header", h)

			orgImg, err := imaging.Decode(assRead)

			if err != nil {
				fmt.Println(err)
				continue
			}

			log.Printf("GET %s\n", missing.Path)

			for _, j := range missing.Types {
				Opt := configuration.Previews[j]
				img, err := Preview(&orgImg, Opt)

				if err != nil {
					fmt.Printf("img error: %s\n", err.Error())
					continue
				}

				path := PreviewName(missing.Path, j)
				//fmt.Printf("Preview Path: %s\n", path)

				hdr := make(http.Header)
				hdr.Add("Content-Type", "image/jpg")
				w, err := preBucket.PutWriter(path, hdr, nil)
				if err != nil {
					fmt.Printf("preview writer: %s\n", err.Error())
					continue
				}
				imaging.Encode(w, img, imaging.JPEG)
				w.Close()
				log.Printf("PUT %s", path)

			}

		}

	}()

	for _, elem := range resp.Contents {
		// Ignore md5 sums, they don't need thumbnails
		if match, _ := regexp.MatchString("^.md5/*", elem.Key); match == true {
			continue
		}

		if match, _ := regexp.MatchString("(.jpg|.png|.jpeg)$", elem.Key); match == false {
			continue
		}

		// If we're using the same bucket, then we need to exclude the preview prefix from the asset list
		if configuration.Preview_Bucket == configuration.Asset_Bucket && strings.HasPrefix(elem.Key, configuration.Preview_Prefix) == true {
			continue
		}

		fmt.Printf("%s does not have prefix %s\n", elem.Key, configuration.Preview_Prefix)

		var types []string

		for key, _ := range configuration.Previews {
			name := PreviewName(elem.Key, key)

			exists := previewBloom.TestString(name)
			if exists == false {
				types = append(types, key)
			}
		}

		if len(types) > 0 {
			missingPreviews <- PreviewPart{
				elem.Key,
				types,
			}

		}

	}

	//fmt.Println(resp2)
}

func PreviewName(obj string, previewtype string) (name string) {
	return path.Join(configuration.Preview_Prefix, previewtype, obj)
}
