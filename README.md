# hermes-go-yaml

A YAML file format validator written in Go.

## Usage

```bash
# Validate specific files
yaml-validator file1.yaml file2.yaml

# Validate from stdin
cat config.yaml | yaml-validator
```

## Features

- Validates YAML syntax and structure
- Reports errors with line and column numbers
- Supports multiple files and stdin input
- Exit code 0 for valid, 1 for invalid

## Build

```bash
go build -o yaml-validator .
```

## Docker

```bash
docker run --rm -v $(pwd):/data ghcr.io/fys05/hermes-go-yaml:latest /data/config.yaml
```

## K8s Deployment

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```
