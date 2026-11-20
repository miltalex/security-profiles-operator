/*
Copyright 2025 The Kubernetes Authors.

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

package v1alpha1

import (
	"encoding/json"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/conversion"

	secprofnodestatusv1 "sigs.k8s.io/security-profiles-operator/api/secprofnodestatus/v1"
)

func (src *SecurityProfileNodeStatus) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*secprofnodestatusv1.SecurityProfileNodeStatus)
	if !ok {
		return fmt.Errorf("expected *secprofnodestatusv1.SecurityProfileNodeStatus, got %T", dstRaw)
	}

	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.NodeName = src.Spec.NodeName
	dst.Status.Status = secprofnodestatusv1.ProfileState(src.Status.Status)

	return nil
}

func (dst *SecurityProfileNodeStatus) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*secprofnodestatusv1.SecurityProfileNodeStatus)
	if !ok {
		return fmt.Errorf("expected *secprofnodestatusv1.SecurityProfileNodeStatus, got %T", srcRaw)
	}

	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.NodeName = src.Spec.NodeName
	dst.Status.Status = ProfileState(src.Status.Status)

	return nil
}

// UnmarshalJSON also accepts objects stored by 0.10.x, which keep nodeName
// and the state string at the top level instead of in spec and status.
func (s *SecurityProfileNodeStatus) UnmarshalJSON(data []byte) error {
	type plain SecurityProfileNodeStatus

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var state string
	if raw, ok := fields["status"]; ok && json.Unmarshal(raw, &state) == nil {
		status, err := json.Marshal(SecurityProfileNodeStatusStatus{Status: ProfileState(state)})
		if err != nil {
			return err
		}

		fields["status"] = status

		// The API server prunes the top-level nodeName before conversion,
		// the node name label of 0.10.x keeps it.
		var nodeName string
		if raw, ok := fields["nodeName"]; ok {
			if err := json.Unmarshal(raw, &nodeName); err != nil {
				return err
			}
		}

		if nodeName == "" {
			var meta struct {
				Labels map[string]string `json:"labels"`
			}
			if raw, ok := fields["metadata"]; ok {
				if err := json.Unmarshal(raw, &meta); err != nil {
					return err
				}
			}

			nodeName = meta.Labels["spo.x-k8s.io/node-name"]
		}

		if nodeName != "" {
			spec := SecurityProfileNodeStatusSpec{}
			if raw, ok := fields["spec"]; ok {
				if err := json.Unmarshal(raw, &spec); err != nil {
					return err
				}
			}

			if spec.NodeName == "" {
				spec.NodeName = nodeName
			}

			if fields["spec"], err = json.Marshal(spec); err != nil {
				return err
			}
		}

		delete(fields, "nodeName")

		if data, err = json.Marshal(fields); err != nil {
			return err
		}
	}

	return json.Unmarshal(data, (*plain)(s))
}
