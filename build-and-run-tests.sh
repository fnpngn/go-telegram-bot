docker build -t test-server -f ./testing/Dockerfile ./testing
docker run --env-file=token.env -p 8080:8080 --name test-server-container test-server
