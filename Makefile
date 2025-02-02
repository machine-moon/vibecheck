all: build


build: 
	@go build -o server

run: build
	@./server

clean:
	@rm -f server

.PHONY: all build run clean
