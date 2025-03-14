date: 2025-03-10
# Docker: Creating and Building Images
![Webserver Illustration](/static/img/docker-logo-blue.png)

This article will show you how to create and build docker images.
## Prerequisites

- Install [Docker](https://docs.docker.com/get-docker/) on your machine.
- Verify installation by running:
  ```sh
  docker --version
  ```

## 1. Create a Dockerfile

A `Dockerfile` is a script that contains instructions to build a Docker image.

Example `Dockerfile` for a simple Go application:

```Dockerfile
# Use the official Golang image as the base
FROM golang:1.20

# Set the working directory
WORKDIR /app

# Copy the Go application source code
COPY . .

# Build the application
RUN go build -o app

# Define the command to run the application
CMD ["./app"]
```

## 2. Build the Docker Image

Run the following command in the directory containing the `Dockerfile`:

```sh
docker build -t my-go-app .
```

Explanation:
- `docker build` - Builds a Docker image.
- `-t my-go-app` - Tags the image as `my-go-app`.
- `.` - Uses the current directory as the build context.

## 3. List Docker Images

To verify that the image was built successfully, run:

```sh
docker images
```

## 4. Run a Docker Container

To create and run a container from the built image:

```sh
docker run --rm -it my-go-app
```

Explanation:
- `docker run` - Runs a container from the specified image.
- `--rm` - Removes the container after it exits.
- `-it` - Runs in interactive mode.
- `my-go-app` - The name of the image.

## 5. Remove Unused Images (Optional)

To free up space, remove unused Docker images:

```sh
docker image prune -a
```

## Conclusion

This guide covered the basics of creating and building Docker images. You can extend the `Dockerfile` to include dependencies, environment variables, and multi-stage builds for optimized production images.
