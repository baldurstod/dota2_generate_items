.PHONY: build clean

BINARY_NAME=dota2_generate_items

build:
	go build ${BINARY_NAME}

run: build
	./${BINARY_NAME} -i "./var" -r "./var" -o "./var"

clean:
	go clean
