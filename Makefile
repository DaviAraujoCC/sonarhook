all: build


build: 
	@echo "Building..."
	CGO_ENABLED=0 go build -o sonarhook .