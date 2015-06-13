# s3preview

Preview service for AWS S3. A small prototype, made at Aalborg Hackathon.

Serves:
- http://{bucketname}/{object}
	- Will pass-through the image from S3
- http://{bucketname}/{object}/thumbnail
	- If the thumbnail exists, redirect to it
	- If the thumbnail doesn't exist, create it, upload it and then show it

## Features

- No Database of Any Kind
- Only here to act as middle-man so that your other services don't have to create preview files