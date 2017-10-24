package algorithm

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/plugin/pkg/scheduler/api"
)

//FilteredHost is used for getting the nodes and pod details and verify and return if pod key matches with annotations
func FilteredHost(args *schedulerapi.ExtenderArgs) ( *schedulerapi.ExtenderFilterResult, error) {
	result := []v1.Node{}
	failedNodesMap := schedulerapi.FailedNodesMap{}

	//Get the list of nodes and pods from base scheduler
	nodes := args.Nodes
	pod := args.Pod

	//Check for presence of Affinity tag in pod specification
	if pod.Spec.Affinity != nil && pod.Spec.Affinity.NodeAffinity != nil {
		if pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil && pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms != nil {

			//get the nodeselector data
			nodeSelectorData := pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
			fmt.Println("Node Affinity tag found in pod specification")
			fmt.Println(nodeSelectorData)

			for _, node := range nodes.Items {
				//allways check for the trust tag signed report
				if cipherVal, ok := node.Annotations["TrustTagSignedReport"]; ok {
					for _, nodeSelector := range nodeSelectorData {
						//match the data from the pod node selector tag to the node annotation 
						if CheckAnnotationAttrib(cipherVal, nodeSelector.MatchExpressions) {
							result = append(result, node)
						} else {
							failedNodesMap[node.Name] = fmt.Sprintf("Annotation validation failed in extended-scheduler")
						}
					}
				}
			}
		} else {
			for _, node := range nodes.Items {
				fmt.Println("No Node Selector terms tag found in pod specification")
				result = append(result, node)
			}
		}
	} else {
		for _, node := range nodes.Items {
			fmt.Println("No Node Affinity tag found in pod specification")
			result = append(result, node)
		}
	}

	glog.V(4).Infof("Returning following nodelist from extended scheduler: %v", result)
	fmt.Println("Returning following nodelist from extended scheduler: %v", result)
	return &schedulerapi.ExtenderFilterResult{
		Nodes:       &v1.NodeList{Items: result},
		NodeNames:   nil,
		FailedNodes: failedNodesMap,
	}, nil
}

