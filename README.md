# ISecL Extended Scheduler 
The ISecL Extended Scheduler verifies trust report and asset tag signature in each of the K8s Worker Node annotation against Pod matching expressions in pod yaml file using ISecL Attestation hub public key. The signature verification ensures the integrity of labels created using isecl hostattribute crds on each of the worker nodes. The verification happens at the time of pod scheduling.

## System Requirements
- RHEL 7.5/7.6
- Epel 7 Repo
- Proxy settings if applicable

## Software requirements
- git
- Go 11.4 or newer

# Step By Step Build Instructions


## Install required shell commands

### Install `go 1.11.4` or newer
The `ISecL extended scheduler` requires Go version 11.4 that has support for `go modules`. The build was validated with version 11.4 version of `go`. It is recommended that you use a newer version of `go` - but please keep in mind that the product has been validated with 1.11.4 and newer versions of `go` may introduce compatibility issues. You can use the following to install `go`.
```shell
wget https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz
tar -xzf go1.11.4.linux-amd64.tar.gz
sudo mv go /usr/local
export GOROOT=/usr/local/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

## Build Extended Scheduler

- Git clone the Extended Scheduler
- Run scripts to build the Extended Scheduler

```shell
git clone https://github.com/intel-secl/k8s-extended-scheduler.git
cd k8s-extended-scheduler
make
```

# Links
https://01.org/intel-secl/


# Installation of binary on kubernetes master machine

# Pre-requisites
	cfssl and cfssljson

# Execute below steps on K8s master node for installing isecl-k8s-scheduler
1. Generate Server cert for extended scheduler by executing below steps.
	```
	K8S_CERT_CN="`hostname` Kubernetes Extended Scheduler"
	K8S_SANS="hostname,IP,kubernetesmaster2"
	K8S_SERVER_CA_KEY=/etc/kubernetes/pki/ca.key
	K8S_SERVER_CA_CERT=/etc/kubernetes/pki/ca.crt

	cat > k8sextscheduler.json <<EOF
	{
	    "hosts": [`echo $K8S_SANS | tr -s "[:space:]"| sed 's/,/\",\"/g' | sed 's/^/\"/' | sed 's/$/\"/' | sed 's/,/,\n/g'`],
            "CN": "$K8S_CERT_CN",
	    "key": {
               "algo": "rsa",
               "size": 2048
            }
	}
	EOF

	cfssl genkey k8sextscheduler.json | cfssljson -bare k8sextscheduler
	cfssl sign -csr=k8sextscheduler.csr -ca-key=${K8S_SERVER_CA_KEY} -ca=${K8S_SERVER_CA_CERT} | cfssljson -bare k8sextscheduler
	exp_date=`cfssl certinfo -cert k8sextscheduler.pem | grep not_after | cut -d T -f 1 | tr -d '[:space:]-"' | cut -d : -f 2`
	chmod 700 k8sextscheduler.pem
	chmod 700 k8sextscheduler-key.pem
	exp_date=`cfssl certinfo -cert k8sextscheduler.pem | grep not_after | cut -d T -f 1 | tr -d '[:space:]-"' | cut -d : -f 2`
	mv  k8sextscheduler-key.pem server.key
	mv k8sextscheduler.pem server.crt
	```

2. Copy the complete compiled source code to K8s master and run the below command to install extended scheduler service

    ```
    make install
    ```
3. Create tag_prefix.conf file in k8s master machine in the path /opt/isecl-k8s-extensions/ with below json.
    
    ```
       {
            "trusted":<<prefix for trust tag>>
       }
    ```


4. In order to bring up the Kubernetes cluster along with the extended scheduler, we need to:

	(a) Create scheduler_policy.json file as per the template below. This contains the configration required by base scheduler to reach the extended scheduler. This file should be placed with the extended scheduler binary (inside /opt/isecl-k8s-extensions/bin). Rewrite URL as per IP address and port.
	```
	{
		"kind" : "Policy",
		"apiVersion" : "v1",
		"predicates" : [
			{"name" : "PodFitsHostPorts"},
			{"name" : "PodFitsResources"},
			{"name" : "NoDiskConflict"},
			{"name" : "MatchNodeSelector"},
			{"name" : "HostName"}
			],
		"priorities" : [
			{"name" : "LeastRequestedPriority", "weight" : 1},
			{"name" : "BalancedResourceAllocation", "weight" : 1},
			{"name" : "ServiceSpreadingPriority", "weight" : 1},
			{"name" : "EqualPriority", "weight" : 1}
			],
		"extenders" : [
			{"urlPrefix": "https://<<k8s master ip address>>:8888/",
	    		 "apiVersion": "v1beta1",
	    		 "filterVerb": "filter",
            		 "weight": 5,
	                 "enableHttps": true
			}
    		]
	}
	```
	

6.  Transfer the public key from the ISecl Attestation Hub and place it in /etc/kubernetes/pki/ folder 
	(a) login to the machine where ISecl Attestation Reporting Hub is installed 
	(b) scp the AH's public key (or use other ways to transfer the key to k8s)
		scp /opt/attestation-hub/configuration/hub_public_key.pem <<user>>@<<k8smasternode>>:/etc/kubernetes/pki/.


7. Edit /opt/isecl-k8s-extenstions/isecl-extended-scheduler-config.json file with appropriate URL ip address of k8smaster, port, paths to certificates(generated in pre-requisites 1) and keys(from step 2)
 


8. Configue the manifest of K8s Base scheduler
  
        * Add scheduler-policy.json under kube-scheduler section /etc/kubernetes/manifests/kube-scheduler.yaml as mentioned below
	```
	spec:
          containers:
	  - command:
            - kube-scheduler
              --policy-config-file : "<path_to_scheduler_policy.json>scheduler_policy.json"
	```

	* Add mount path for isecl extended scheduler under container section /etc/kubernetes/manifests/kube-scheduler.yaml as mentioned below
	```
	containers:
		volumeMounts:
		- mountPath: /etc/kubernetes/scheduler.conf
		name: kubeconfig
		readOnly: true
		- mountPath: /opt/isecl-k8s-extensions/bin/
		name: extendedsched
		readOnly: true
	```

	* Add volume path for isecl extended scheduler under volumes section /etc/kubernetes/manifests/kube-scheduler.yaml as mentioned below
	```
	spec:
	volumes:
	- hostPath:
		path: /etc/kubernetes/scheduler.conf
		type: FileOrCreate
		name: kubeconfig
	- hostPath:
		path: /opt/isecl-k8s-extensions/bin/
		type: ""
		name: extendedsched
	```


9. Run below commands to enable service daemon (to activate newly added service)
       
	```
	systemctl daemon-reload   
	```

10. Run the extended scheduler using below command	
	```
	systemctl start isecl-k8s-scheduler.service
	```

11. Restart Kubelet which restart all the k8s services including kube base schedular
	```
	systemctl restart kubelet
	```

12. To check status of this service run below command
	```
	systemctl status isecl-k8s-scheduler.service
	```

13. To stop this service run below command
	```
	systemctl stop isecl-k8s-scheduler.service
	```

14. To check the port where extended scheduler is listening, the default the port for extended schedular service is 8888
	```
	netstat -tlpn | grep 8888
	```
