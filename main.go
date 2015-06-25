package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type PreviewOptions struct {
	Width  int
	Height int
	Method string
}

type Config struct {
	Previews       map[string]PreviewOptions
	Asset_Bucket   string
	Asset_Prefix   string
	Preview_Bucket string
	Preview_Prefix string
	StorageDomain  string
	ListenPort     int
	ListenPortSSL  int
	Seed           bool
	AWS_Key        string
	AWS_Secret     string
}

var configuration Config

func init() {
}

func main() {

	app := cli.NewApp()
	app.Name = "s3preview"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Andri Óskarsson",
			Email: "andri80@gmail.com",
		},
	}
	app.Usage = "Previews (or thumbnails) for AWS S3 objects"
	app.Version = "0.1.0"
	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "aws-key",
			Usage:  "AWS Access Key. Needs to have S3 Access.",
			EnvVar: "AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "aws-secret",
			Usage:  "AWS Secret Key",
			EnvVar: "AWS_SECRET_ACCESS_KEY",
		},
		cli.IntFlag{
			Name:   "port",
			Usage:  "HTTP listen port",
			Value:  80,
			EnvVar: "PORT",
		},
		cli.StringFlag{
			Name:   "asset-bucket",
			Usage:  "Bucket to retrieve originals from",
			EnvVar: "ASSET_BUCKET",
		},
		cli.StringFlag{
			Name:   "asset-prefix",
			Usage:  "Prefix for Asset Bucket. E.g. /originals",
			EnvVar: "ASSET_PREFIX",
		},
		cli.StringFlag{
			Name:   "preview-bucket",
			Usage:  "Bucket to store previews on. Needs to be public for redirects to work.",
			EnvVar: "PREVIEW_BUCKET",
		},
		cli.StringFlag{
			Name:   "preview-prefix",
			Usage:  "Prefix for Preview Bucket. E.g. /s3preview",
			EnvVar: "PREVIEW_PREFIX",
		},
		cli.BoolFlag{
			Name:   "generate",
			Usage:  "Will generate missing previews during startup.",
			EnvVar: "GENERATE",
		},
	}

	// TODO: Finish CLI'fying the program

	app.Action = func(c *cli.Context) {
		Configure(c)
		PopulatePreviewCache()
		if c.Bool("generate") == true {
			GenerateMissing()
		}

		r := mux.NewRouter()

		http.Handle("/", r)
		registerHandlers(r)
		//fmt.Println(configuration.Previews)
		//	fmt.Println(configuration.Previews["small"])
		fmt.Printf("I'm listening (port %d)...\n", configuration.ListenPort)
		http.ListenAndServe(fmt.Sprintf(":%d", configuration.ListenPort), nil)
	}

	app.Run(os.Args)
}

func Configure(c *cli.Context) {
	if _, err := toml.DecodeFile("config.toml", &configuration); err != nil {
		panic(err)
	}
	if c.String("aws-key") != "" {
		configuration.AWS_Key = c.String("aws-key")
	}
	if c.String("aws-secret") != "" {
		configuration.AWS_Secret = c.String("aws-secret")
	}
	if c.String("preview-bucket") != "" {
		configuration.Preview_Bucket = c.String("preview-bucket")
	}
	if c.String("preview-prefix") != "" {
		configuration.Preview_Prefix = c.String("preview-prefix")
	}
	if c.String("asset-bucket") != "" {
		configuration.Asset_Bucket = c.String("asset-bucket")
	}
	if c.String("asset-prefix") != "" {
		configuration.Asset_Prefix = c.String("asset-prefix")
	}

	if c.Int("port") != 0 {
		configuration.ListenPort = c.Int("port")
	}

	if configuration.AWS_Key == "" || configuration.AWS_Secret == "" {
		fmt.Println("AWS Configuration Required (key, secret, asset bucket, preview bucket).")
		os.Exit(1)
	}

	if configuration.StorageDomain == "" {
		configuration.StorageDomain = "s3.amazonaws.com"
	}

	if configuration.Asset_Bucket == "" || configuration.Preview_Bucket == "" {
		fmt.Println("Asset-Bucket and Preview-Bucket must be configured")
		os.Exit(1)
	}

	// Default value for preview_prefix, but only if sharing buckets
	if configuration.Asset_Bucket == configuration.Preview_Bucket {
		if configuration.Preview_Prefix == "" {
			configuration.Preview_Prefix = "s3preview/"
		}
	}
}
