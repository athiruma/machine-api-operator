/*
Copyright 2018 The Kubernetes Authors.

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

package machine

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2/textlogger"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func init() {
	textLoggerConfig := textlogger.NewConfig()
	textLoggerConfig.AddFlags(flag.CommandLine)
	logf.SetLogger(textlogger.NewLogger(textLoggerConfig))
}

const (
	timeout = time.Second * 10
)

var (
	cfg       *rest.Config
	ctx       = context.Background()
	k8sClient client.Client
	testEnv   *envtest.Environment
)

func TestMachineController(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Machine Controller Suite")
}

var _ = BeforeSuite(func(ctx SpecContext) {
	By("bootstrapping test environment")
	var err error
	cfg, testEnv, err = StartEnvTest()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	Expect(testEnv.Stop()).To(Succeed())
})

func TestMain(m *testing.M) {
	// Register required object kinds with global scheme.
	if err := machinev1.Install(scheme.Scheme); err != nil {
		log.Fatalf("cannot add scheme: %v", err)
	}
	if err := configv1.Install(scheme.Scheme); err != nil {
		log.Fatalf("cannot add scheme: %v", err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func StartEnvTest() (*rest.Config, *envtest.Environment, error) {
	testEnv := &envtest.Environment{

		CRDInstallOptions: envtest.CRDInstallOptions{
			Paths: []string{
				filepath.Join("..", "..", "..", "vendor", "github.com", "openshift", "api", "machine", "v1beta1", "zz_generated.crd-manifests", "0000_10_machine-api_01_machines-CustomNoUpgrade.crd.yaml"),
			},
		},
	}

	var err error
	if cfg, err = testEnv.Start(); err != nil {
		return nil, nil, err
	}

	return cfg, testEnv, nil
}
