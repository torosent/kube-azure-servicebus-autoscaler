.PHONY: test clean compile build push

IMAGE=torosent/kube-azure-servicebus-autoscaler
VERSION=1.0.0

clean:
	rm -f kube-azure-servicebus-autoscaler

compile: clean
	GOOS=linux go build .

build: compile
	docker build -t $(IMAGE):$(VERSION) .

push: build
	docker push $(IMAGE):$(VERSION)
