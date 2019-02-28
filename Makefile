# CIT K8S Extensions Extended Scheduler
# Works in tandem with K8s scheduler to return a filtered list of nodes as per predicates on CRDs
# Author:  <manux.ullas@intel.com>
DESCRIPTION="CIT K8S Extended Scheduler"

SERVICE=citk8sscheduler
SYSTEMINSTALLDIR=/opt/cit_k8s_extensions/bin/
CONFIGDIR=/opt/cit_k8s_extensions/config
SERVICEINSTALLDIR=/etc/systemd/system/
SERVICECONFIG=${SERVICE}.service
VERSION := 1.0-SNAPSHOT
BUILD := `date +%FT%T%z`

# LDFLAGS
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"


# Generate the service binary and executable
.DEFAULT_GOAL: $(SERVICE)
$(SERVICE):
	@rm -f ${GOPATH}/pkg/mod/k8s.io/kubernetes@v0.0.0-20170928221357-0b9efaeb34a2/plugin/pkg/scheduler/api/zz_generated.deepcopy.go
	go build ${LDFLAGS} -o ${SERVICE}-${VERSION} ${SOURCES}

# Install the service binary and the service config files
.PHONY: install
install:
	@mkdir -p ${SYSTEMINSTALLDIR}
	@mkdir -p ${CONFIGDIR}
	@cp cit-extended-scheduler-config.json ${CONFIGDIR}
	@cp -f ${SERVICE}-${VERSION} ${SYSTEMINSTALLDIR}
	@cp -f ${SERVICECONFIG} ${SERVICEINSTALLDIR}
        

# Uninstalls the service binary and the service config files
.PHONY: uninstall
uninstall:
	@service ${SERVICE} stop && rm -rf ${SERVICEINSTALLDIR}/${SERVICE} ${SERVICEINSTALLDIR}/${SERVICECONFIG} ${CONFIGDIR} ${SYSTEMINSTALLDIR}${SERVICE}-${VERSION}

# Removes the generated service config and binary files
.PHONY: clean
clean:
	@rm -f ${SERVICE}-${VERSION}
