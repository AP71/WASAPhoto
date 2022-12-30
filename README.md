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
$ docker run -it --rm -v /service/database/file:/app/service/database/file:rw -p 3000:3000 wasa-photos-backend:latest
```

### Frontend

```
$ docker run -it --rm -p 8081:80 wasa-photos-frontend:latest
```

## License

See [LICENSE](LICENSE).