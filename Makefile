build:
	@go build -o ./bin/main

local:
	@sudo cp bin/main /usr/local/bin/knav
