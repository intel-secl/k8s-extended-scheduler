# CIT K8S Extensions Extended Scheduler
# Works in tandem with K8s scheduler to return a filtered list of nodes as per predicates on CRDs
# Author:  <manux.ullas@intel.com>
DESCRIPTION="CIT K8S Extended Scheduler"

SERVICE=citk8sscheduler
SYSTEMINSTALLDIR=/opt/cit_k8s_extensions/bin/
SERVICEINSTALLDIR=/etc/systemd/system/
SERVICECONFIG=${SERVICE}.service
VERSION := 1.0-SNAPSHOT
BUILD := `date +%FT%T%z`

# LDFLAGS
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"


# Generate the service binary and executable
.DEFAULT_GOAL: $(SERVICE)
$(SERVICE):
	glide install -v
	# Remove this file from the deps since it causes build failures
	@rm -f vendor/k8s.io/kubernetes/plugin/pkg/scheduler/api/zz_generated.deepcopy.go
	go build ${LDFLAGS} -o ${SERVICE}-${VERSION} ${SOURCES}

# Install the service binary and the service config files
.PHONY: install
install:
	@service ${SERVICE} stop
	@mkdir -p ${SYSTEMINSTALLDIR}
	@cp -f ${SERVICE}-${VERSION} ${SYSTEMINSTALLDIR}
	@cp -f ${SERVICECONFIG} ${SERVICEINSTALLDIR}
        

# Uninstalls the service binary and the service config files
.PHONY: uninstall
uninstall:
	@service ${SERVICE} stop && rm -f ${SERVICEINSTALLDIR}/${SERVICE} ${SERVICEINSTALLDIR}/${SERVICECONFIG}

# Removes the generated service config and binary files
.PHONY: clean
clean:
	@rm -rf vendor/
	@rm -f ${SERVICE}-${VERSION}
