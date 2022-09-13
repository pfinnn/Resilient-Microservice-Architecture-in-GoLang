TOPTARGETS := all build test build-linux clean
MAKEFILES = $(shell find . -maxdepth 2 -type f -name Makefile)
SUBDIRS   = $(filter-out ./,$(dir $(MAKEFILES)))

SERVICES := forgeservice smelterservice
GCP_PROJECTID = resilient-microservice
GCP_CLUSTERID = cluster-rs

$(TOPTARGETS): $(SUBDIRS)
$(SUBDIRS):
		$(MAKE) -C $@ $(MAKECMDGOALS)

.PHONY: $(TOPTARGETS) $(SUBDIRS)

localBuild:
	cd forgeservice && make build-linux
	cd smelterservice && make build-linux
	docker-compose build

localUp: localBuild
	docker-compose up -d

localDown:
	docker-compose down

deploy:
	kubectl create -f ./deployments/

undeploy:
	kubectl delete -f ./deployments/

clusterUp:
	gcloud container clusters create $(GCP_CLUSTERID) --num-nodes=2 --no-enable-ip-alias

clusterDown:
	gcloud container clusters delete $(GCP_CLUSTERID)
	@for SVC in $(SERVICES); do gcloud container images delete gcr.io/$(GCP_PROJECTID)/$$SVC --force-delete-tags --quiet || exit 1; done;

test-global-resilience:
	cd forgeservice && make build-linux
	cd smelterservice && make build-linux
	docker-compose build
	docker-compose up -d
	cd forgeservice && make test-resilience
	cd smelterservice && make test-resilience
