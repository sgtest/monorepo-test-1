/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package federation

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-5/pkg/api"
)

// ServerAddressByClientCIDR helps the client to determine the server address that they should use, depending on the clientCIDR that they match.
type ServerAddressByClientCIDR struct {
	// The CIDR with which clients can match their IP to figure out the server address that they should use.
	ClientCIDR string
	// Address of this server, suitable for a client that matches the above CIDR.
	// This can be a hostname, hostname:port, IP or IP:port.
	ServerAddress string
}

// ClusterSpec describes the attributes of a kubernetes cluster.
type ClusterSpec struct {
	// A map of client CIDR to server address.
	// This is to help clients reach servers in the most network-efficient way possible.
	// Clients can use the appropriate server address as per the CIDR that they match.
	// In case of multiple matches, clients should use the longest matching CIDR.
	ServerAddressByClientCIDRs []ServerAddressByClientCIDR
	// Name of the secret containing kubeconfig to access this cluster.
	// The secret is read from the kubernetes cluster that is hosting federation control plane.
	// Admin needs to ensure that the required secret exists. Secret should be in the same namespace where federation control plane is hosted and it should have kubeconfig in its data with key "kubeconfig".
	// This will later be changed to a reference to secret in federation control plane when the federation control plane supports secrets.
	// This can be left empty if the cluster allows insecure access.
	// +optional
	SecretRef *api.LocalObjectReference
}

type ClusterConditionType string

// These are valid conditions of a cluster.
const (
	// ClusterReady means the cluster is ready to accept workloads.
	ClusterReady ClusterConditionType = "Ready"
	// ClusterOffline means the cluster is temporarily down or not reachable
	ClusterOffline ClusterConditionType = "Offline"
)

// ClusterCondition describes current state of a cluster.
type ClusterCondition struct {
	// Type of cluster condition, Complete or Failed.
	Type ClusterConditionType
	// Status of the condition, one of True, False, Unknown.
	Status api.ConditionStatus
	// Last time the condition was checked.
	// +optional
	LastProbeTime metav1.Time
	// Last time the condition transit from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// (brief) reason for the condition's last transition.
	// +optional
	Reason string
	// Human readable message indicating details about last transition.
	// +optional
	Message string
}

// ClusterStatus is information about the current status of a cluster updated by cluster controller periodically.
type ClusterStatus struct {
	// Conditions is an array of current cluster conditions.
	// +optional
	Conditions []ClusterCondition
	// Zones is the list of availability zones in which the nodes of the cluster exist, e.g. 'us-east1-a'.
	// These will always be in the same region.
	// +optional
	Zones []string
	// Region is the name of the region in which all of the nodes in the cluster exist.  e.g. 'us-east1'.
	// +optional
	Region string
}

// +genclient=true
// +nonNamespaced=true

// Information about a registered cluster in a federated kubernetes setup. Clusters are not namespaced and have unique names in the federation.
type Cluster struct {
	metav1.TypeMeta
	// Standard object's metadata.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta

	// Spec defines the behavior of the Cluster.
	// +optional
	Spec ClusterSpec
	// Status describes the current status of a Cluster
	// +optional
	Status ClusterStatus
}

// A list of all the kubernetes clusters registered to the federation
type ClusterList struct {
	metav1.TypeMeta
	// Standard list metadata.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta

	// List of Cluster objects.
	Items []Cluster
}

// Temporary/alpha structures to support custom replica assignments within FederatedReplicaSet.

// A set of preferences that can be added to federated version of ReplicaSet as a json-serialized annotation.
// The preferences allow the user to express in which clusters he wants to put his replicas within the
// mentioned FederatedReplicaSet.
type FederatedReplicaSetPreferences struct {
	// If set to true then already scheduled and running replicas may be moved to other clusters to
	// in order to bring cluster replicasets towards a desired state. Otherwise, if set to false,
	// up and running replicas will not be moved.
	// +optional
	Rebalance bool

	// A mapping between cluster names and preferences regarding local ReplicaSet in these clusters.
	// "*" (if provided) applies to all clusters if an explicit mapping is not provided. If there is no
	// "*" that clusters without explicit preferences should not have any replicas scheduled.
	// +optional
	Clusters map[string]ClusterReplicaSetPreferences
}

// Preferences regarding number of replicas assigned to a cluster replicaset within a federated replicaset.
type ClusterReplicaSetPreferences struct {
	// Minimum number of replicas that should be assigned to this Local ReplicaSet. 0 by default.
	// +optional
	MinReplicas int64

	// Maximum number of replicas that should be assigned to this Local ReplicaSet. Unbounded if no value provided (default).
	// +optional
	MaxReplicas *int64

	// A number expressing the preference to put an additional replica to this LocalReplicaSet. 0 by default.
	Weight int64
}
