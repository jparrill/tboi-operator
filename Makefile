USER := padajuan
OP_NAME := tboi-operator
DOCK_TAG := latest

.PHONY: gen
gen:
	operator-sdk generate k8s

.PHONY: build
build:
	operator-sdk build docker.io/${USER}/${OP_NAME}:${DOCK_TAG}
	docker push docker.io/$USER/${OP_NAME}:${DOCK_TAG}

.PHONY: deploy
deploy:
	oc new-project ${OP_NAME}
	oc create -f deploy/rbac.yaml
	oc create -f deploy/crd.yaml
	oc create -f deploy/operator.yaml
	oc create -f deploy/cr.yaml
	oc get all

.PHONY: clean
clean:
	oc delete -f deploy/cr.yaml
	oc delete -f deploy/operator.yaml
	oc delete -f deploy/rbac.yaml
	oc delete -f deploy/crd.yaml
	oc delete project ${OP_NAME}
	oc project default
