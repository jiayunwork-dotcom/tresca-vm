# tresca-vm — Go Tresca/von Mises 屈服核算 HTTP 后端与命令行工具

The server starts on `:8080` by default and exposes `POST /api/yield` and
`POST /api/pure-shear`; CLI subcommands `yield`, `tensor`, and `pure-shear`
provide the same calculations in a terminal.

## Build / Run / Test

```text
go build ./...
go run . serve -addr :8080
go run . example -file example/uniaxial.json
go test ./...
```

## Evaluation Image

Evaluation-specific files (do not overwrite project Dockerfile/README):

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md` (this file)

Build and verify in container:

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
