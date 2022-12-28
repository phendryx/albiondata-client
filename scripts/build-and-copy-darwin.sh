docker build -f ./Dockerfile.build.darwin -t albiondataclient-darwin .
docker run --name builder albiondataclient-darwin
docker cp builder:/usr/src/app/update-darwin-amd64.gz ./update-darwin-amd64.gz
docker cp builder:/usr/src/app/albiondata-client-amd64-mac.zip ./albiondata-client-amd64-mac.zip
