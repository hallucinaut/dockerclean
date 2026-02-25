# dockerclean

Docker Cleanup Tool - removes unused Docker containers, volumes, and images.

## Purpose

Clean up Docker resources by removing stopped containers, unused volumes, and dangling images.

## Installation

```bash
go build -o dockerclean ./cmd/dockerclean
```

## Usage

```bash
dockerclean [--dry-run]
```

### Examples

```bash
# Dry run - preview what would be removed
dockerclean --dry-run

# Actual cleanup
dockerclean
```

## Output

### Dry Run

```
=== DOCKER CLEANUP DRY RUN ===

Containers to remove: 3
Volumes to remove: 2
Images to remove: 5

Containers:
  a3f5b8c2 - old-api-service
  b4e6c9d3 - worker-process
  c5f7d0e4 - test-container

Volumes:
  a3f5b8c2 - legacy-data
  b4e6c9d3 - unused-volume

Images:
  a3f5b8c2 - api-service:latest
  b4e6c9d3 - worker:dev
  ...

=== SAFE CLEANUP COMMANDS ===
docker container prune -f
docker volume prune -f
docker image prune -a -f
```

### Actual Cleanup

```
=== DOCKER CLEANUP ===

Removing 3 containers...
Removing 2 volumes...
Removing 5 images...

Cleanup complete!
```

## Dependencies

- Go 1.21+
- github.com/fatih/color
- Docker CLI installed

## Build and Run

```bash
# Build
go build -o dockerclean ./cmd/dockerclean

# Run dry run
go run ./cmd/dockerclean --dry-run

# Run cleanup
go run ./cmd/dockerclean
```

## License

MIT