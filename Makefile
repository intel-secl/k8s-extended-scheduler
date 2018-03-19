# CIT K8S Extended Scheduler
# Builds, Installs and Uninstalls the K8S Extended Scheduler
DESCRIPTION="CIT K8S Extended Scheduler"

SERVICE=citk8s-extended-scheduler
SYSTEMINSTALLDIR=/opt/citk8s/$(SERVICE)

VERSION := 1.0.0
BUILD := `date +%FT%T%z`

# LDFLAGS
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"


# Generate the service binary and executable
.DEFAULT_GOAL: $(SERVICE)
$(SERVICE):
	@glide update -v
# Remove this file from the deps since it causes build failures
	@rm -f vendor/k8s.io/kubernetes/plugin/pkg/scheduler/api/zz_generated.deepcopy.go
	go build ${LDFLAGS} -o ${SERVICE} ${SOURCES}

# Install the service binary and the service config files
.PHONY: install
install:
	@cp -f ${SERVICE} ${SERVICEINSTALLDIR}

# Uninstalls the service binary and the service config files
.PHONY: uninstall
uninstall:
	@rm -f ${SERVICEINSTALLDIR}/${SERVICE}

# Removes the generated service config and binary files
.PHONY: clean
clean:
	@rm -rf vendor/
	@rm -f ${SERVICE}
