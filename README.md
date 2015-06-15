# s3preview

Note: This project is VERY early development and I will probably change (or break) stuff.

Service that sits in between Amazon's s3 service and your thumbnail (or preview) needs. It will create preview files on the fly, upload them to s3 and redirect directly to s3 if the preview flie exists.

## Motivation

I have a lot of pet projects and a reoccuring theme is that I need some sort of thumbnails for various purposes and importing graphic libraries to all projects seems silly. So to be all hip and stuff, I decided to create a microservice that has the single purpose of creating preview files (and pointing you to them).

I also needed a project with a very narrow scope for Aalborg Hackathon.

## API

- GET /{object}?t={preview-type}

Default preview types are "small", "large", "thumbnail"

That's all folks!

## Configuration

Previews and various tweaks are stored in config.toml.

TODO: Cleanup defaults and allow env overrides to config file.

### Thumbnail

Resize into desired width and height. Then crop the image.

### Resize

Resize to desired width and height. If you want to maintain aspect ratio, set one of the size values to 0.
