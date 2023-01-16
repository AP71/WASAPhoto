# Wasa Photo!

## How to build container images

### Backend

```sh
$ docker build -t wasa-photos-backend:latest -f Dockerfile.backend .
```

### Frontend

```sh
$ docker build -t wasa-photos-frontend:latest -f Dockerfile.frontend .
```

## How to run container images

### Backend

```sh
$ docker run -it --rm -u "$(id -u):$(id -g)" -p 3000:3000 wasa-photos-backend:latest
```

### Frontend

```
$ docker run -it --rm -p 8080:80 wasa-photos-frontend:latest
```

## License

See [LICENSE](LICENSE).
