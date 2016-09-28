run-example: examples/node-docker-example.tar.gz
	go run main.go $^ --config examples/cb-demo.json

examples/node-docker-example.tar.gz:
	cd examples && curl -O https://storage.googleapis.com/container-builder-examples/node-docker-example.tar.gz
