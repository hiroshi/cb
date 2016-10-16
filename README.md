# cb - a Container Builder

It builds container images locally using [Google Cloud Container Builder config file](https://cloud.google.com/container-builder/docs/api/build-requests#build_steps).

Why not just do `docker build`? It will be useful to provide an easy way to manage multiple steps builds.
- We love small docker images:
  - Don't want to contain golang environment. We love a single go binary docker image.
  - Don't Want to contain frontend js build environment for a web app.

### What it looks like?

```
steps:
- name: gcr.io/cloud-builders/docker
  args: ["build", "-t", "cb-build", "-f", "Dockerfile.build", "."]
- name: cb-build
  args: ["cp", "/go/src/cb/cb", "/workspace"]
- name: gcr.io/cloud-builders/docker
  args: ["build", "-t", "cb", "."]
```
This is an example config file. It will build a golang single binary image of `cb` command itself (not useful though).
- 1st step - Build a temporary image. It builds a go binary using `golang` base image as usual.
- 2nd step - Run the resulted image of 1st step. It copies the golang binary in the image to workspace volume.
- 3rd step - Build a final image from scratch. Just add the `cb` command from workplace volume.

```
$ docker images
REPOSITORY      TAG          IMAGE ID            CREATED             SIZE
cb              latest       05994f135ea4        2 days ago          3.208 MB
cb-build        latest       61f9b946f604        2 days ago          680.9 MB
...
```

## Install

`go get -u github.com/hiroshi/cb`

Make sure you have `$GOPATH/bin` in your `$PATH`.

## Usage

`cb SOURCE.tar.gz --config CONFIG.(json|yml)`

### Notes
- The [`source`](https://cloud.google.com/container-builder/docs/api/build-requests#source_location) field in config will be ignored as well as `gcloud alpha container builds create` do. Specify SOURCE as 1st argument.
- The [`images`](https://cloud.google.com/container-builder/docs/api/build-requests#resulting_images) field in config will be ignored. The `cb` command is intended for local builds so always pushing images are not supposed to be welcome.

## How it works
- 1) Create a volume for `workspace` with `docker volume create`.
- 2) Expand `SOURCE` into the `workspace` volume with `docker copy`.
- 3) `docker run` an image with volumes `/var/run/docker.sock//var/run/docker.sock/`, `WORKSPACE_VOLUME:/workspace`.
  - If the image have docker command like [this](https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/docker), you can do `docker build` or anything in container with SOURCE at hand.
- 4) Repeat 3) with different image and args as you specifed in `steps` field of CONFIG.

Do you get it? No? See and run examples, I hope it may help you understand.

## Examples
`make run-example`

## TODO

- Support `wait_for` and `id` fields of [`steps`](https://cloud.google.com/container-builder/docs/api/build-requests#build_steps)

## References
- [Build request - Google Cloud Container Builder](https://cloud.google.com/container-builder/docs/api/build-requests)
- [gcloud alpha container builds create](https://cloud.google.com/sdk/gcloud/reference/alpha/container/builds/create)
