# tresca-vm

tresca-vm is a Go yield-criteria calculator for Tresca and von Mises
conditions. It takes a 3D stress tensor or the three principal stresses plus a
yield strength, returns both equivalent stresses, safety factors, and which
criterion reaches yield first. The pure-shear endpoint gives the comparison
values for `2*tau = sigma_y` and `sqrt(3)*tau = sigma_y` boundaries. The
service is exposed through HTTP JSON endpoints and CLI subcommands with no
web page.

## Usage

Run the HTTP server:

```bash
go run . serve -addr :8080
```

Evaluate from the command line:

```bash
go run . yield -principal 100,0,-50 -sigmaY 200
go run . tensor -sxx 100 -syy -50 -sigmaY 200
go run . pure-shear -tau 100 -sigmaY 200
```

Run the uniaxial scenario:

```bash
go run . example -file example/uniaxial.json
```

The example uses principal stresses `[120, 0, 0]` with `sigma_y = 240`; both
criteria return a safety factor of 2.

## HTTP API

```text
POST /api/yield        {"principal":[100,0,-50],"sigma_y":200}
POST /api/yield        {"sxx":100,"syy":-50,"sigma_y":200}
POST /api/pure-shear   {"tau":100,"sigma_y":200}
POST /api/plane-stress {"sxx":120,"syy":-40,"txy":30,"sigma_y":300}
GET  /health
```

Invalid yield strengths, asymmetric tensors, and non-finite components return
an error body with HTTP 400.

## Stress Conventions

Principal stresses are sorted so `s1 >= s2 >= s3`. Tresca uses
`s1 - s3`; von Mises uses `sqrt(0.5 * sum((si-sj)^2))`. Hydrostatic stress
does not enter the deviatoric equivalent stress. Plane stress keeps the zero
third principal value in the sorted list.

## Code Layout

```text
internal/tensor   tensor representation, parsing, eigenvalues, invariants
internal/yield    Tresca/Mises equivalents, safety, yield classification
internal/server   HTTP handlers and JSON responses
internal/cli      subcommand parsing and terminal output
example/          offline scenario JSON files
```

## Build and Test

```bash
export GOTOOLCHAIN=local CGO_ENABLED=0
go build ./...
go test ./...
```

The Dockerfile builds the server binary and starts it on port 8080.
