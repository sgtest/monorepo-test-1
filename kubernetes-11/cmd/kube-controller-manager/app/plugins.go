/*
Copyright 2014 The Kubernetes Authors.

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

package app

import (
	// This file exists to force the desired plugin implementations to be linked.
	// This should probably be part of some configuration fed into the build for a
	// given binary target.

	"fmt"

	// Cloud providers
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/apis/componentconfig"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/cloudprovider/providers"

	// Volume plugins
	"github.com/golang/glog"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/cloudprovider"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/cloudprovider/providers/aws"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/cloudprovider/providers/azure"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/cloudprovider/providers/gce"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/cloudprovider/providers/openstack"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/cloudprovider/providers/photon"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/cloudprovider/providers/vsphere"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/util/io"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/aws_ebs"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/azure_dd"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/azure_file"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/cinder"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/flexvolume"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/flocker"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/gce_pd"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/glusterfs"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/host_path"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/nfs"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/photon_pd"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/portworx"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/quobyte"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/rbd"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/scaleio"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-11/pkg/volume/vsphere_volume"
)

// ProbeAttachableVolumePlugins collects all volume plugins for the attach/
// detach controller. VolumeConfiguration is used ot get FlexVolumePluginDir
// which specifies the directory to search for additional third party volume
// plugins.
// The list of plugins is manually compiled. This code and the plugin
// initialization code for kubelet really, really need a through refactor.
func ProbeAttachableVolumePlugins(config componentconfig.VolumeConfiguration) []volume.VolumePlugin {
	allPlugins := []volume.VolumePlugin{}

	allPlugins = append(allPlugins, aws_ebs.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, gce_pd.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, cinder.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, flexvolume.ProbeVolumePlugins(config.FlexVolumePluginDir)...)
	allPlugins = append(allPlugins, portworx.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, vsphere_volume.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, azure_dd.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, photon_pd.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, scaleio.ProbeVolumePlugins()...)
	return allPlugins
}

// ProbeControllerVolumePlugins collects all persistent volume plugins into an
// easy to use list. Only volume plugins that implement any of
// provisioner/recycler/deleter interface should be returned.
func ProbeControllerVolumePlugins(cloud cloudprovider.Interface, config componentconfig.VolumeConfiguration) []volume.VolumePlugin {
	allPlugins := []volume.VolumePlugin{}

	// The list of plugins to probe is decided by this binary, not
	// by dynamic linking or other "magic".  Plugins will be analyzed and
	// initialized later.

	// Each plugin can make use of VolumeConfig.  The single arg to this func contains *all* enumerated
	// options meant to configure volume plugins.  From that single config, create an instance of volume.VolumeConfig
	// for a specific plugin and pass that instance to the plugin's ProbeVolumePlugins(config) func.

	// HostPath recycling is for testing and development purposes only!
	hostPathConfig := volume.VolumeConfig{
		RecyclerMinimumTimeout:   int(config.PersistentVolumeRecyclerConfiguration.MinimumTimeoutHostPath),
		RecyclerTimeoutIncrement: int(config.PersistentVolumeRecyclerConfiguration.IncrementTimeoutHostPath),
		RecyclerPodTemplate:      volume.NewPersistentVolumeRecyclerPodTemplate(),
		ProvisioningEnabled:      config.EnableHostPathProvisioning,
	}
	if err := AttemptToLoadRecycler(config.PersistentVolumeRecyclerConfiguration.PodTemplateFilePathHostPath, &hostPathConfig); err != nil {
		glog.Fatalf("Could not create hostpath recycler pod from file %s: %+v", config.PersistentVolumeRecyclerConfiguration.PodTemplateFilePathHostPath, err)
	}
	allPlugins = append(allPlugins, host_path.ProbeVolumePlugins(hostPathConfig)...)

	nfsConfig := volume.VolumeConfig{
		RecyclerMinimumTimeout:   int(config.PersistentVolumeRecyclerConfiguration.MinimumTimeoutNFS),
		RecyclerTimeoutIncrement: int(config.PersistentVolumeRecyclerConfiguration.IncrementTimeoutNFS),
		RecyclerPodTemplate:      volume.NewPersistentVolumeRecyclerPodTemplate(),
	}
	if err := AttemptToLoadRecycler(config.PersistentVolumeRecyclerConfiguration.PodTemplateFilePathNFS, &nfsConfig); err != nil {
		glog.Fatalf("Could not create NFS recycler pod from file %s: %+v", config.PersistentVolumeRecyclerConfiguration.PodTemplateFilePathNFS, err)
	}
	allPlugins = append(allPlugins, nfs.ProbeVolumePlugins(nfsConfig)...)
	allPlugins = append(allPlugins, glusterfs.ProbeVolumePlugins()...)
	// add rbd provisioner
	allPlugins = append(allPlugins, rbd.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, quobyte.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, azure_file.ProbeVolumePlugins()...)

	allPlugins = append(allPlugins, flocker.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, portworx.ProbeVolumePlugins()...)
	allPlugins = append(allPlugins, scaleio.ProbeVolumePlugins()...)

	if cloud != nil {
		switch {
		case aws.ProviderName == cloud.ProviderName():
			allPlugins = append(allPlugins, aws_ebs.ProbeVolumePlugins()...)
		case gce.ProviderName == cloud.ProviderName():
			allPlugins = append(allPlugins, gce_pd.ProbeVolumePlugins()...)
		case openstack.ProviderName == cloud.ProviderName():
			allPlugins = append(allPlugins, cinder.ProbeVolumePlugins()...)
		case vsphere.ProviderName == cloud.ProviderName():
			allPlugins = append(allPlugins, vsphere_volume.ProbeVolumePlugins()...)
		case azure.CloudProviderName == cloud.ProviderName():
			allPlugins = append(allPlugins, azure_dd.ProbeVolumePlugins()...)
		case photon.ProviderName == cloud.ProviderName():
			allPlugins = append(allPlugins, photon_pd.ProbeVolumePlugins()...)
		}
	}

	return allPlugins
}

// AttemptToLoadRecycler tries decoding a pod from a filepath for use as a recycler for a volume.
// If successful, this method will set the recycler on the config.
// If unsuccessful, an error is returned. Function is exported for reuse downstream.
func AttemptToLoadRecycler(path string, config *volume.VolumeConfig) error {
	if path != "" {
		recyclerPod, err := io.LoadPodFromFile(path)
		if err != nil {
			return err
		}
		if err = volume.ValidateRecyclerPodTemplate(recyclerPod); err != nil {
			return fmt.Errorf("Pod specification (%v): %v", path, err)
		}
		config.RecyclerPodTemplate = recyclerPod
	}
	return nil
}
