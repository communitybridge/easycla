// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// a mixed ACL edit adds only the added manager and removes only the removed one (lfx-self-serve-ops#346 item 3)
func TestUpdateCLAPermissionsAppliesEachACLEntrySeparately(t *testing.T) {
	type call struct {
		managers []string
		action   string
	}
	var calls []call
	previous := claManagerPermissionUpdater
	claManagerPermissionUpdater = func(_ *service, _ Signature, managers []string, action string) error {
		calls = append(calls, call{managers: append([]string(nil), managers...), action: action})
		return nil
	}
	defer func() { claManagerPermissionUpdater = previous }()

	image := func(acl []string) map[string]events.DynamoDBAttributeValue {
		return map[string]events.DynamoDBAttributeValue{
			"signature_id":             events.NewStringAttribute("sig-ccla"),
			"signature_type":           events.NewStringAttribute("ccla"),
			"signature_reference_type": events.NewStringAttribute("company"),
			"signature_reference_id":   events.NewStringAttribute("company-1"),
			"signature_project_id":     events.NewStringAttribute("cla-group-1"),
			"signature_acl":            events.NewStringSetAttribute(acl),
		}
	}
	s := &service{}
	err := s.UpdateCLAPermissions(events.DynamoDBEventRecord{Change: events.DynamoDBStreamRecord{OldImage: image([]string{"keep", "removed"}), NewImage: image([]string{"keep", "added"})}})
	require.NoError(t, err)
	require.Len(t, calls, 2)
	for _, c := range calls {
		require.Len(t, c.managers, 1, "every ACL difference is applied on its own")
		switch c.managers[0] {
		case "added":
			assert.Equal(t, AddCLAManager, c.action)
		case "removed":
			assert.Equal(t, DeleteCLAManager, c.action)
		default:
			t.Fatalf("unexpected manager %q", c.managers[0])
		}
	}
}
