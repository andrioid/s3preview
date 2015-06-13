# s3preview

Preview service for AWS S3. A small prototype, made at Aalborg Hackathon.

Serves:
- http://{object}
	- Will pass-through the image from S3 (used for debugging purposes)
- http://{object}/{preview-type}
	- If the preview exists, redirect to it
	- If the preview doesn't exist, create it, upload it and then show it

## Features

- No Database of Any Kind
- Only here to act as middle-man so that your other services don't have to create preview files

## Preview Types

Configured in config.toml

Default types are: "large", "small", "thumbnail"