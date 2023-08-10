APP_NAME := app
ARGS_START := start server
BIN_PATH := ./bin
CMD_PATH := ./cmd
DOCS_OUTPUT_PATH := ./internal/docs
HELM_CHARTS_PATH := ./deploy/charts
IMAGE_PATH := ${APP_NAME}/${APP_NAME}
IMAGE_REGISTRY := docker.io
IMAGE_TAG := latest
K8S_NAMESPACE := default
KAFKA_ALIAS := localhost
KAFKA_CERTS_PATH := ./.data/kafka/config/certs
SWAGGO_URL := github.com/swaggo/swag

# everything before this next line is configurable in the Makefile.properties
-include Makefile.properties
# everything after this line is a derived value from the above variables

HELM_CHART_PATH := ${HELM_CHARTS_PATH}/${APP_NAME}
HELM_RELEASE_NAME := $(APP_NAME)
IMAGE_URL := ${IMAGE_REGISTRY}/${IMAGE_PATH}
KAFKA_CA_KEY_PATH := ${KAFKA_CERTS_PATH}/ca-key
KAFKA_CA_CERT_PATH := ${KAFKA_CERTS_PATH}/ca-cert
KAFKA_CLIENT_CERT_PATH := ${KAFKA_CERTS_PATH}/client-cert
KAFKA_CLIENT_KEY_PATH := ${KAFKA_CERTS_PATH}/client-key
KAFKA_CLIENT_P12_PATH := ${KAFKA_CERTS_PATH}/client.p12
KAFKA_JKS_KEYSTORE_PATH := ${KAFKA_CERTS_PATH}/kafka.keystore.jks
KAFKA_JKS_TRUSTSTORE_PATH := ${KAFKA_CERTS_PATH}/kafka.truststore.jks
KIND_CLUSTER_NAME := ${APP_NAME}

ifeq ("${GOOS}", "windows")
BINARY_EXT := ".exe"
endif

binary: docs-swaggo
	@echo "building binary for ${APP_NAME} for os/arch $$(go env GOOS)/$$(go env GOARCH)..."
	@mkdir -p "${BIN_PATH}"
	@go build \
		-ldflags "\
			-extldflags 'static' -s -w \
			-X ${APP_NAME}/internal/constants.AppName=${APP_NAME} \
			-X ${APP_NAME}/internal/constants.Version=$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse HEAD | head -c 6) \
		" \
		-o "${BIN_PATH}/${APP_NAME}${BINARY_EXT}" \
		"${CMD_PATH}/${APP_NAME}"
	@cp "${BIN_PATH}/${APP_NAME}${BINARY_EXT}" "${BIN_PATH}/${APP_NAME}_$$(go env GOOS)_$$(go env GOARCH)${BINARY_EXT}"

deps:
	@echo "updating dependencies..."
	@go mod tidy
	@go mod vendor

docs-swaggo: 
	@echo "generating documentation package..."
	swag init \
		--generalInfo "./internal/api/api.go" \
		--output "${DOCS_OUTPUT_PATH}" \
		--parseInternal \
		--parseDependency \
		--parseDepth 8

image:
	@echo "building image ${IMAGE_URL}:${IMAGE_TAG}..."addr
	docker build -t ${IMAGE_URL}:${IMAGE_TAG} .

deploy-k8s:
	@echo "deploying to kubernetes cluster..."
	@cd ${HELM_CHART_PATH} \
		&& helm upgrade --install \
			--namespace ${K8S_NAMESPACE} \
			--create-namespace \
			--values ./values.yaml \
			"${HELM_RELEASE_NAME}" .

deploy-kind:
	@echo "deploying to local kubernetes cluster..."
	@kubectl config use-context kind-${KIND_CLUSTER_NAME}
	@$(MAKE) deploy-k8s

install-swaggo:
	@echo "installing swag from ${SWAGGO_URL}..."
	go get -u ${SWAGGO_URL}/cmd/swag@latest

install-swaggo-ci:
	@echo "installing swag from ${SWAGGO_URL}..."
	go install ${SWAGGO_URL}/cmd/swag@latest

kafka-jks: # ref https://www.ibm.com/docs/en/cloud-paks/cp-biz-automation/20.0.x?topic=emitter-preparing-ssl-certificates-kafka
	rm -rf ${KAFKA_CERTS_PATH}/*
	mkdir -p ${KAFKA_CERTS_PATH}
	echo '*' > ${KAFKA_CERTS_PATH}/.gitignore
	echo '!.gitignore' >> ${KAFKA_CERTS_PATH}/.gitignore

	# create certificate authority
	openssl req -new -x509 -keyout ${KAFKA_CA_KEY_PATH} -out ${KAFKA_CA_CERT_PATH} -days 365

	# create client certificate
	openssl req -new -newkey rsa:2048 -nodes -keyout ${KAFKA_CLIENT_KEY_PATH} -out ${KAFKA_CLIENT_CERT_PATH} -days 365
	openssl x509 -req -days 365 -in ${KAFKA_CLIENT_CERT_PATH} -CA ${KAFKA_CA_CERT_PATH} -CAkey ${KAFKA_CA_KEY_PATH} -out ${KAFKA_CLIENT_CERT_PATH} -set_serial 01 -sha256

	# package client data into client keystore
	openssl pkcs12 -export -in ${KAFKA_CLIENT_CERT_PATH} -inkey ${KAFKA_CLIENT_KEY_PATH} -name user > ${KAFKA_CLIENT_P12_PATH}
	keytool -importkeystore -srckeystore ${KAFKA_CLIENT_P12_PATH} -destkeystore ${KAFKA_JKS_KEYSTORE_PATH} -srcstoretype pkcs12 -alias user

	# package certificate authority into server truststore
	keytool -keystore ${KAFKA_JKS_TRUSTSTORE_PATH} -alias CARoot -import -file ${KAFKA_CA_CERT_PATH}

	chmod 644 ${KAFKA_CERTS_PATH}/*

kind-load: image
	@echo "loading image into local kubernetes cluster..."
	@kind load docker-image ${IMAGE_URL}:${IMAGE_TAG} --name ${KIND_CLUSTER_NAME}

nats-nkey:
	nk -gen user -pubout

start: docs
	@go run "${CMD_PATH}/${APP_NAME}" ${ARGS_START}

start-kind:
	@echo "initialising local kubernetes cluster..."
	@kind create cluster --name ${KIND_CLUSTER_NAME}

start-kafka:
	@docker-compose run kafka

start-mongo:
	@docker-compose run mongo

start-mysql:
	@docker-compose run mysql

start-nats:
	@docker-compose run nats

start-postgres:
	@docker-compose run postgres

start-redis:
	@docker-compose run redis

test:
	@go test -v -coverpkg=./... -coverprofile=./tests/cover.out ./...
	@go tool cover -func ./tests/cover.out
	@go tool cover -html ./tests/cover.out -o ./tests/cover.html
