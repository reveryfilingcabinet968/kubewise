// Copyright 2026 KubeWise Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kube

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// NewClientSet builds a Kubernetes clientset from the given kubeconfig path and context.
// If kubeconfig is empty, it falls back to $KUBECONFIG, then ~/.kube/config.
// If context is empty, the current context from the kubeconfig is used.
func NewClientSet(kubeconfig, context string) (*kubernetes.Clientset, error) {
	config, err := buildConfig(kubeconfig, context)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes clientset: %w", err)
	}
	return clientset, nil
}

// NewMetricsClientSet builds a metrics-server clientset from the given kubeconfig path and context.
// If kubeconfig is empty, it falls back to $KUBECONFIG, then ~/.kube/config.
// If context is empty, the current context from the kubeconfig is used.
func NewMetricsClientSet(kubeconfig, context string) (*metricsv.Clientset, error) {
	config, err := buildConfig(kubeconfig, context)
	if err != nil {
		return nil, err
	}

	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating metrics clientset: %w", err)
	}
	return metricsClient, nil
}

func buildConfig(kubeconfig, context string) (*rest.Config, error) {
	path := resolveKubeconfig(kubeconfig)

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: path}
	overrides := &clientcmd.ConfigOverrides{}
	if context != "" {
		overrides.CurrentContext = context
	}

	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
	restConfig, err := config.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("building kubeconfig from %q: %w", path, err)
	}

	return restConfig, nil
}

func resolveKubeconfig(explicit string) string {
	if explicit != "" {
		return explicit
	}
	if env := os.Getenv("KUBECONFIG"); env != "" {
		return env
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}
