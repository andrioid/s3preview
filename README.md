# s3preview

## What is it?

A web-service that attempts to remove the complexity of previews (or thumbnails) from applications. It utilitizes Amazon's Object Storage (s3) for both original files and previews.

### Example Use

If you have an image on s3, say "s3preview-demo/1.jpg" in the bucket "andridk-assets" and for some reason, you would like to serve a 100 pixel wide preview of that image.

Then you can set s3preview up to read from "andridk-assets" and link to it from HTML like

```
<img src="http://s3preview-demo.herokuapp.com/s3preview-demo/1.jpg?t=small">
```

If the preview exists on your "preview-bucket", then s3preview will redirect the browser to the actual file on s3. If the preview doesn't exist, it will be generated, served to the browser and uploaded to s3.

![s3preview demo image](http://s3preview-demo.herokuapp.com/s3preview-demo/1.jpg?t=small)

## Features

- Resizing while maintaining aspect-ratio
- Thumbnail. Resize and crop
- Transparent. It will serve the same paths as your s3 bucket.
- Should be fast (I'll measure that later)
- On-start generation of any missing preview images.
- Caching of preview-list from s3. *
	- I use a Bloom filter so there is a possibility of false positives. The result is that s3preview might redirect to a preview file that doesn't exist yet. Don't worry, I have a plan.
- Stateless. You can scale out by starting more instances of the service.


## That sounds awesome. How can I try it?

### Heroku

If you see a picture below the example, then you've already been served an image from s3preview. If not, Heroku is probably forcing my Dyno to sleep.

You may also try your own: 

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/andrioid/s3preview/tree/master)

### Docker

### Build your own

```bash
git clone https://github.com/andrioid/s3preview
docker build -t yours3preview .
docker run -ti \
	-e AWS_ACCESS_KEY_ID=<yourkey> -e AWS_SECRET_ACCESS_KEY=<yoursecret> -e ASSET_BUCKET=<yourassetbucket> -e PREVIEW_BUCKET=<yourpreviewbucket> \
	-e ASSET_PREFIX="" -e PREVIEW_PREFIX="s3preview/"
	yours3preview 
```

### Use my [automated build](https://registry.hub.docker.com/u/andrioid/s3preview/)

```bash
docker run -ti \
	-e AWS_ACCESS_KEY_ID=<yourkey> -e AWS_SECRET_ACCESS_KEY=<yoursecret> -e ASSET_BUCKET=<yourassetbucket> -e PREVIEW_BUCKET=<yourpreviewbucket> \
	-e ASSET_PREFIX="" -e PREVIEW_PREFIX="s3preview/"
	andrioid/s3preview
```

## Status

It's not pretty, but it works.


## Configuration

### config.toml
These are the standard thumbnails (that may change). 

There are two supported methods.

- "resize": Scales the image. If either height, or width is 0 it will respect aspect-ratio
- "thumbnail": Scales and crop the image.

```toml
# Squared 200 pixel thumbnail (crop)
[previews.sq200]
width = 200
height = 200
method = "thumbnail"

# Squared 100 pixel thumbnail
[previews.thumbnail]
width = 100
height = 100
method = "thumbnail"

# 100 pixel high scaled
[previews.small]
width = 0
height = 100
method = "resize"

# Squared 500 pixel thumb 
[previews.large]
width = 500
height = 500
method = "thumbnail"
```

### Command Line Interface (and env)

Can be produced by running `s3preview --help`

```
NAME:
   s3preview - Previews (or thumbnails) for AWS S3 objects

USAGE:
   s3preview [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR(S): 
   Andri Óskarsson <andri80@gmail.com> 

COMMANDS:
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --aws-key 		AWS Access Key. Needs to have S3 Access. [$AWS_ACCESS_KEY_ID]
   --aws-secret 	AWS Secret Key [$AWS_SECRET_ACCESS_KEY]
   --port "80"		HTTP listen port [$PORT]
   --asset-bucket 	Bucket to retrieve originals from [$ASSET_BUCKET]
   --asset-prefix 	Prefix for Asset Bucket. E.g. /originals [$ASSET_PREFIX]
   --preview-bucket 	Bucket to store previews on. Needs to be public for redirects to work. [$PREVIEW_BUCKET]
   --preview-prefix 	Prefix for Preview Bucket. E.g. /s3preview [$PREVIEW_PREFIX]
   --generate		Will generate missing previews during startup. [$GENERATE]
   --help, -h		show help
```

## TODO

In no particular order

- Clean up the code (isn't it always there?)
- Add support for mp4 video files.
- Better handle too much load. Return an error when too busy with preview-generation.

## Thanks
- [Aalborg Hackathon](https://www.linkedin.com/groups/Aalborg-Hackathon-7453429/about): For getting me out of the house so I could work on this.

## License

MIT License ([tldr](https://tldrlegal.com/license/mit-license))