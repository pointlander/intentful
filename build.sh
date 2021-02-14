CGO_ENABLED=0 go build
docker build -t pointlander/intentful:latest .
docker push pointlander/intentful:latest
