/*
Copyright 2019 The Kubernetes Authors.

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

package kubeconfig

import (
	"os"
	"path"
	"path/filepath"

	"k8s.io/apimachinery/pkg/util/sets"
)

const kubeconfigEnv = "KUBECONFIG"

/*
paths returns the list of paths to be considered for kubeconfig files
where explicitPath is the value of --kubeconfig

Logic based on kubectl

https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands

- If the --kubeconfig flag is set, then only that file is loaded. The flag may only be set once and no merging takes place.

- If $KUBECONFIG environment variable is set, then it is used as a list of paths (normal path delimiting rules for your system). These paths are merged. When a value is modified, it is modified in the file that defines the stanza. When a value is created, it is created in the first file that exists. - If no files in the chain exist, then it creates the last file in the list.

- Otherwise, ${HOME}/.kube/config is used and no merging takes place.
*/
func paths(explicitPath string, getEnv func(string) string) []string {
	if explicitPath != "" {
		return []string{explicitPath}
	}

	paths := discardEmptyAndDuplicates(
		filepath.SplitList(getEnv(kubeconfigEnv)),
	)
	if len(paths) != 0 {
		return paths
	}

	return []string{path.Join(homeDir(), ".kube", "config")}
}

// pathForMerge returns the file that kubectl would merge into
func pathForMerge(explicitPath string) string {
	// find the first file that exists
	p := paths(explicitPath, os.Getenv)
	if len(p) == 1 {
		return p[0]
	}
	for _, filename := range p {
		if fileExists(filename) {
			return filename
		}
	}
	// otherwise the last file
	return p[len(p)-1]
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func discardEmptyAndDuplicates(paths []string) []string {
	seen := sets.NewString()
	kept := 0
	for _, p := range paths {
		if p != "" && !seen.Has(p) {
			paths[kept] = p
			kept++
			seen.Insert(p)
		}
	}
	return paths[:kept]
}
