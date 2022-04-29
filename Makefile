all: build


build: 
	@echo "Building..."
	CGO_ENABLED=0 go build -o sonarhook .

docker-build:
	@echo "Building Docker image..."
	docker build -t ${IMG} .

docker-push:
	@echo "Pushing Docker image..."
	docker push ${IMG}