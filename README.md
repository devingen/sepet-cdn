# Sepet CDN
CDN server for Sepet.

## Building and running the server

Build and run `./cmd/sepet-cdn/sepet-cdn.go`.

Check the environment variables defined in `config/app.go`.
They must all have the prefix `SEPET_` like `SEPET_S3_ACCESS_KEY_ID`.

Provide `SEPET_S3_ENDPOINT` while using MinIO as the file server.

## Running the Docker image

```
// pull the latest image
docker pull devingen/sepet-cdn:VERSION_HERE

// stop and remove any existing container
docker stop sepet-cdn && docker rm sepet-cdn

// run the container
docker run \
  --restart always \
  --name sepet-cdn \
  -e SEPET_CDN_PORT=80 \
  -e SEPET_CDN_LOG_LEVEL=debug \
  -e SEPET_CDN_DAL_UPDATE_INTERVAL=5s \
  -e SEPET_CDN_CACHE_RESET_INTERVAL=1m \
  -e SEPET_CDN_API_URL=http://localhost:1005 \
  -e SEPET_CDN_S3_ENDPOINT=http://localhost:9000 \
  -e SEPET_CDN_S3_ACCESS_KEY_ID=ACCESSKEYIDFORTHEFILESERVER \
  -e SEPET_CDN_S3_SECRET_ACCESS_KEY=ACCESSKEYFORTHEFILESERVER \
  -e SEPET_CDN_S3_REGION=region-of-the-cdn \
  -e SEPET_CDN_S3_BUCKET=the-root-bucket-name-in-s3 \
  -p 80:80 \
  devingen/sepet-cdn:VERSION_HERE
```

## Development

### Releasing new Docker image
```
docker build --platform linux/amd64 -t devingen/sepet-cdn:0.0.6 .
docker push devingen/sepet-cdn:0.0.6
```
