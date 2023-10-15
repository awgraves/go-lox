BINARY_NAME= lox

build:
	@go build -o ${BINARY_NAME} main.go

run: build
	@./${BINARY_NAME}
