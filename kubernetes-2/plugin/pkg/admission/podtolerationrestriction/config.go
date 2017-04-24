/*
Copyright 2017 The Kubernetes Authors.

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

package podtolerationrestriction

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/apimachinery/registered"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	internalapi "github.com/sourcegraph/monorepo-test-1/kubernetes-2/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-2/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction/install"
	versionedapi "github.com/sourcegraph/monorepo-test-1/kubernetes-2/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction/v1alpha1"
	"github.com/sourcegraph/monorepo-test-1/kubernetes-2/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction/validation"
)

var (
	groupFactoryRegistry = make(announced.APIGroupFactoryRegistry)
	registry             = registered.NewOrDie(os.Getenv("KUBE_API_VERSIONS"))
	scheme               = runtime.NewScheme()
	codecs               = serializer.NewCodecFactory(scheme)
)

func init() {
	install.Install(groupFactoryRegistry, registry, scheme)
}

// LoadConfiguration loads the provided configuration.
func loadConfiguration(config io.Reader) (*internalapi.Configuration, error) {
	// if no config is provided, return a default configuration
	if config == nil {
		externalConfig := &versionedapi.Configuration{}
		scheme.Default(externalConfig)
		internalConfig := &internalapi.Configuration{}
		if err := scheme.Convert(externalConfig, internalConfig, nil); err != nil {
			return nil, err
		}
		return internalConfig, nil
	}
	// we have a config so parse it.
	data, err := ioutil.ReadAll(config)
	if err != nil {
		return nil, err
	}
	decoder := codecs.UniversalDecoder()
	decodedObj, err := runtime.Decode(decoder, data)
	if err != nil {
		return nil, err
	}
	externalConfig, ok := decodedObj.(*internalapi.Configuration)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T", decodedObj)
	}

	if err := validation.ValidateConfiguration(externalConfig); err != nil {
		return nil, err
	}

	return externalConfig, nil
}
