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

// This file exists to force the desired plugin implementations to be linked.
// This should probably be part of some configuration fed into the build for a
// given binary target.
import (
	// Cloud providers
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/pkg/cloudprovider/providers"

	// Admission policies
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/admit"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/alwayspullimages"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/antiaffinity"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/defaulttolerationseconds"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/deny"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/exec"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/gc"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/imagepolicy"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/initialresources"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/limitranger"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/namespace/autoprovision"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/namespace/exists"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/namespace/lifecycle"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/persistentvolume/label"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/podnodeselector"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/podpreset"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/podtolerationrestriction"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/resourcequota"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/security/podsecuritypolicy"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/securitycontext/scdeny"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/serviceaccount"
	_ "github.com/sourcegraph/monorepo-test-1/kubernetes-13/plugin/pkg/admission/storageclass/default"
)
