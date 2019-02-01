# ISecL K8S Extensions Extended Scheduler
# Works in tandem with K8s scheduler to return a filtered list of nodes as per predicates on CRDs
# Author:  <manux.ullas@intel.com>
DESCRIPTION="ISecL K8S Extended Scheduler"

SERVICE=isecl-k8s-scheduler
SYSTEMINSTALLDIR=/opt/isecl-k8s-extensions/bin/
CONFIGDIR=/opt/cit-k8s-extensions/config
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
	@mkdir -p ${SYSTEMINSTALLDIR}
	@mkdir -p ${CONFIGDIR}
	@cp isecl-extended-scheduler-config.json ${CONFIGDIR}
	@cp -f ${SERVICE}-${VERSION} ${SYSTEMINSTALLDIR}
	@cp -f ${SERVICECONFIG} ${SERVICEINSTALLDIR}
        

# Uninstalls the service binary and the service config files
.PHONY: uninstall
uninstall:
	@service ${SERVICE} stop && rm -rf ${SERVICEINSTALLDIR}/${SERVICE} ${SERVICEINSTALLDIR}/${SERVICECONFIG} ${CONFIGDIR} ${SYSTEMINSTALLDIR}${SERVICE}-${VERSION}

# Removes the generated service config and binary files
.PHONY: clean
clean:
	@rm -rf vendor/
	@rm -f ${SERVICE}-${VERSION}
