# cb - a Container Builder

It builds container images locally using [Google Cloud Container Builder config file](https://cloud.google.com/container-builder/docs/api/build-requests#build_steps).

Why just do `docker build`? It is useful to provide an easy way to manage multiple steps builds.
(e.g.
- Don't want to contain go command and depencendies just for a single go binary docker image.
- Want to transpile frontent js files into a bundle then add it to a docker image.

## Install

`go get -u github.com/hiroshi/cb`

Make sure you have `$GOPATH/bin` in your `$PATH`.

## Usage

`cb SOURCE.tar.gz --config CONFIG.(json|yml)`

### Notes
- The [`source`](https://cloud.google.com/container-builder/docs/api/build-requests#source_location) field in config will be ignored as well as `gcloud alpha container builds create` do. Specify SOURCE as 1st argument.
- The ['images'](https://cloud.google.com/container-builder/docs/api/build-requests#resulting_images) field in config will be ignored. The `cb` command is intended for local builds so always pushing images are not supposed to be welcome.

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

- Support `env`, `wait_for`, `id`, `dir` of [`steps`](https://cloud.google.com/container-builder/docs/api/build-requests#build_steps)

## References
- [Build request - Google Cloud Container Builder](https://cloud.google.com/container-builder/docs/api/build-requests)
- [gcloud alpha container builds create](https://cloud.google.com/sdk/gcloud/reference/alpha/container/builds/create)
