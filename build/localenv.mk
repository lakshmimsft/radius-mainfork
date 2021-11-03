# ------------------------------------------------------------
# Copyright (c) Microsoft Corporation.
# Licensed under the MIT License.
# ------------------------------------------------------------

# Local development experience using Radius is facilitated by a local K8s API server.
# We use kcp (https://github.com/kcp-dev/kcp)
#
# Assumptions:
# 1. kcp binary is located under ~/.rad/bin
# 2. kcp has been started (cd ~/.rad/bin && ./kcp start)
# 3. The resulting kcp configuration files are under ~/.rad/bin/.kcp

##@ Support for local development using Radius
radiusd-run: check-kcp-running ## Run Radius daemon (controller host/API server)
	KUBECONFIG=~/.rad/bin/.kcp/data/admin.kubeconfig go run ./cmd/radiusd/main.go -zap-devel

radiusd-crd-install: check-kcp-running ## Install CRDs related to local development (set up the local K8s API server)
	KUBECONFIG=~/.rad/bin/.kcp/data/admin.kubeconfig kubectl apply -f ./deploy/Chart/crds/radius.dev_executables.yaml

radiusd-crd-uninstall: check-kcp-running ## Uninstall CRDs relatd to local development
	KUBECONFIG=~/.rad/bin/.kcp/data/admin.kubeconfig kubectl delete -f ./deploy/Chart/crds/radius.dev_executables.yaml

check-kcp-running:
	./build/check-kcp.sh