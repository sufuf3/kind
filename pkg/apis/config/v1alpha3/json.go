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

package v1alpha3

import (
	"encoding/json"
	"fmt"
	"strings"
)

/*
Custom JSON (de)serialization for these types
TODO: just use yaml in v1alpha4 ...
*/

// MarshalJSON implements custom encoding for JSON
// https://golang.org/pkg/encoding/json/
func (m *Mount) MarshalJSON() ([]byte, error) {
	type Alias Mount
	name, ok := MountPropagationValueToName[m.Propagation]
	if !ok {
		return nil, fmt.Errorf("unknown propagation value: %v", m.Propagation)
	}
	return json.Marshal(&struct {
		Propagation string `json:"propagation"`
		*Alias
	}{
		Propagation: name,
		Alias:       (*Alias)(m),
	})
}

// UnmarshalJSON implements custom decoding for JSON and Yaml
// https://golang.org/pkg/encoding/json/
func (m *Mount) UnmarshalJSON(data []byte) error {
	type Alias Mount
	aux := &struct {
		Propagation string `json:"propagation"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// if unset, will fallback to the default (0)
	if aux.Propagation != "" {
		val, ok := MountPropagationNameToValue[aux.Propagation]
		if !ok {
			return fmt.Errorf("unknown propagation value: %s", aux.Propagation)
		}
		m.Propagation = MountPropagation(val)
	}
	return nil
}

// MarshalJSON implements custom encoding for JSON
// https://golang.org/pkg/encoding/json/
func (p *PortMapping) MarshalJSON() ([]byte, error) {
	type Alias PortMapping
	name, ok := PortMappingProtocolValueToName[p.Protocol]
	if !ok {
		return nil, fmt.Errorf("unknown protocol value: %v", p.Protocol)
	}
	return json.Marshal(&struct {
		Protocol string `json:"protocol"`
		*Alias
	}{
		Protocol: name,
		Alias:    (*Alias)(p),
	})
}

// UnmarshalJSON implements custom decoding for JSON
// https://golang.org/pkg/encoding/json/
func (p *PortMapping) UnmarshalJSON(data []byte) error {
	type Alias PortMapping
	aux := &struct {
		Protocol string `json:"protocol"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Protocol != "" {
		val, ok := PortMappingProtocolNameToValue[strings.ToUpper(aux.Protocol)]
		if !ok {
			return fmt.Errorf("unknown protocol value: %s", aux.Protocol)
		}
		p.Protocol = PortMappingProtocol(val)
	}
	return nil
}
