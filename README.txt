##############################################################################
Pre-requisites for building the code
##############################################################################
1. Install GO
2. Install MAVEN
3. Install maven plugin for GO - from https://github.com/raydac/mvn-golang
4. Set GOROOT to path of go folder. For example : /usr/local/go
5. Install JAVA SDK and set JAVA_HOME

Run below command from terminal inside main directory
-- mvn package
On build success, you can find binary in "bin" folder
Binary name - citk8sscheduler-1.0-SNAPSHOT

##############################################################################
Installation of binary
##############################################################################
Pre-requisites
1. Kubernetes cluster should be up and running
2. CIT custom controller should be installed and running
3. Edit config file in src folder
4. Copy the binary from bin folder and config file from src folder to root folder 
Install this binary on kubernetes master
