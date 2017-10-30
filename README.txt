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
1. Install the AH public Key and Server cert and server public key
2. Re-Start the base scheduler with the policy.json configuration
3. CIT extended scheduler should be installed 
4. Copy the binary from bin folder and config file from src folder to /opt folder 
5. Run the exetended scheduler as below
	service cit-extended-scheduler start
