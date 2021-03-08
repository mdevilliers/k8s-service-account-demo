
KIND_INSTANCE=k8s-service-account-demo

# creates a K8s instance
.PHONY: k8s_new
k8s_new:
	kind create cluster --config ./kind/kind.yaml --name $(KIND_INSTANCE)

# deletes a k8s instance
.PHONY: k8s_drop
k8s_drop:
	kind delete cluster --name $(KIND_INSTANCE)

# sets KUBECONFIG for the K8s instance
.PHONY: k8s_connect
k8s_connect:
	kind export kubeconfig --name $(KIND_INSTANCE)

# loads the docker containers into the kind environment
.PHONY: k8s_side_load
k8s_side_load:
	kind load docker-image x/client --name $(KIND_INSTANCE)
	kind load docker-image x/server --name $(KIND_INSTANCE)

# builds all the applications
.PHONY: build_apps
build_apps:
	docker build -t x/client -f ./services/client/Dockerfile ./services/client
	docker build -t x/server -f ./services/server/Dockerfile ./services/server

.PHONY: deploy_apps
deploy_apps: k8s_connect
	kubectl apply -f ./k8s/client-deployment.yaml
	kubectl apply -f ./k8s/server-deployment.yaml
