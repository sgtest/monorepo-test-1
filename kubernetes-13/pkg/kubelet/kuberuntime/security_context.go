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

package kuberuntime

import (
	"fmt"

	"github.com/sourcegraph/monorepo-test-1/kubernetes-13/pkg/api/v1"
	runtimeapi "github.com/sourcegraph/monorepo-test-1/kubernetes-13/pkg/kubelet/api/v1alpha1/runtime"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-13/pkg/securitycontext"
)

// determineEffectiveSecurityContext gets container's security context from v1.Pod and v1.Container.
func (m *kubeGenericRuntimeManager) determineEffectiveSecurityContext(pod *v1.Pod, container *v1.Container, uid *int64, username string) *runtimeapi.LinuxContainerSecurityContext {
	effectiveSc := securitycontext.DetermineEffectiveSecurityContext(pod, container)
	synthesized := convertToRuntimeSecurityContext(effectiveSc)
	if synthesized == nil {
		synthesized = &runtimeapi.LinuxContainerSecurityContext{}
	}

	// set RunAsUser.
	if synthesized.RunAsUser == nil {
		if uid != nil {
			synthesized.RunAsUser = &runtimeapi.Int64Value{Value: *uid}
		}
		synthesized.RunAsUsername = username
	}

	// set namespace options and supplemental groups.
	synthesized.NamespaceOptions = &runtimeapi.NamespaceOption{
		HostNetwork: pod.Spec.HostNetwork,
		HostIpc:     pod.Spec.HostIPC,
		HostPid:     pod.Spec.HostPID,
	}
	podSc := pod.Spec.SecurityContext
	if podSc != nil {
		if podSc.FSGroup != nil {
			synthesized.SupplementalGroups = append(synthesized.SupplementalGroups, *podSc.FSGroup)
		}

		if podSc.SupplementalGroups != nil {
			synthesized.SupplementalGroups = append(synthesized.SupplementalGroups, podSc.SupplementalGroups...)
		}
	}
	if groups := m.runtimeHelper.GetExtraSupplementalGroupsForPod(pod); len(groups) > 0 {
		synthesized.SupplementalGroups = append(synthesized.SupplementalGroups, groups...)
	}

	return synthesized
}

// verifyRunAsNonRoot verifies RunAsNonRoot.
func verifyRunAsNonRoot(pod *v1.Pod, container *v1.Container, uid int64) error {
	effectiveSc := securitycontext.DetermineEffectiveSecurityContext(pod, container)
	if effectiveSc == nil || effectiveSc.RunAsNonRoot == nil {
		return nil
	}

	if effectiveSc.RunAsUser != nil {
		if *effectiveSc.RunAsUser == 0 {
			return fmt.Errorf("container's runAsUser breaks non-root policy")
		}
		return nil
	}

	if uid == 0 {
		return fmt.Errorf("container has runAsNonRoot and image will run as root")
	}

	return nil
}

// convertToRuntimeSecurityContext converts v1.SecurityContext to runtimeapi.SecurityContext.
func convertToRuntimeSecurityContext(securityContext *v1.SecurityContext) *runtimeapi.LinuxContainerSecurityContext {
	if securityContext == nil {
		return nil
	}

	sc := &runtimeapi.LinuxContainerSecurityContext{
		Capabilities:   convertToRuntimeCapabilities(securityContext.Capabilities),
		SelinuxOptions: convertToRuntimeSELinuxOption(securityContext.SELinuxOptions),
	}
	if securityContext.RunAsUser != nil {
		sc.RunAsUser = &runtimeapi.Int64Value{Value: *securityContext.RunAsUser}
	}
	if securityContext.Privileged != nil {
		sc.Privileged = *securityContext.Privileged
	}
	if securityContext.ReadOnlyRootFilesystem != nil {
		sc.ReadonlyRootfs = *securityContext.ReadOnlyRootFilesystem
	}

	return sc
}

// convertToRuntimeSELinuxOption converts v1.SELinuxOptions to runtimeapi.SELinuxOption.
func convertToRuntimeSELinuxOption(opts *v1.SELinuxOptions) *runtimeapi.SELinuxOption {
	if opts == nil {
		return nil
	}

	return &runtimeapi.SELinuxOption{
		User:  opts.User,
		Role:  opts.Role,
		Type:  opts.Type,
		Level: opts.Level,
	}
}

// convertToRuntimeCapabilities converts v1.Capabilities to runtimeapi.Capability.
func convertToRuntimeCapabilities(opts *v1.Capabilities) *runtimeapi.Capability {
	if opts == nil {
		return nil
	}

	capabilities := &runtimeapi.Capability{
		AddCapabilities:  make([]string, len(opts.Add)),
		DropCapabilities: make([]string, len(opts.Drop)),
	}
	for index, value := range opts.Add {
		capabilities.AddCapabilities[index] = string(value)
	}
	for index, value := range opts.Drop {
		capabilities.DropCapabilities[index] = string(value)
	}

	return capabilities
}
