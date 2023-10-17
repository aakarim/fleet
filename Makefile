dist/starsh:
	# make the binary as small as possible
	go build -ldflags="-s -w" -o ./dist/starsh ./cmd/starsh
clean:
	rm -rf ./dist