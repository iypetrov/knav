build:
	@go build -o ./bin/main

local:
	@sudo cp bin/main /usr/local/bin/knav
	@sudo codesign --force --sign - /usr/local/bin/knav
	@codesign -dv /usr/local/bin/knav 2>&1 | head -3
