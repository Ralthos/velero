/*
Copyright the Velero contributors.

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

package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"

	velerov1api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

func TestBackupPhaseHasNoArtifacts(t *testing.T) {
	tests := []struct {
		phase velerov1api.BackupPhase
		want  bool
	}{
		{velerov1api.BackupPhaseNew, true},
		{velerov1api.BackupPhaseQueued, true},
		{velerov1api.BackupPhaseReadyToStart, true},
		{velerov1api.BackupPhaseFailedValidation, true},

		// Anything from InProgress onwards may have written a partial log or other
		// artifacts, so these keep the behaviour callers have today.
		{velerov1api.BackupPhaseInProgress, false},
		{velerov1api.BackupPhaseWaitingForPluginOperations, false},
		{velerov1api.BackupPhaseWaitingForPluginOperationsPartiallyFailed, false},
		{velerov1api.BackupPhaseFinalizing, false},
		{velerov1api.BackupPhaseFinalizingPartiallyFailed, false},
		{velerov1api.BackupPhaseCompleted, false},
		{velerov1api.BackupPhasePartiallyFailed, false},
		{velerov1api.BackupPhaseFailed, false},
		{velerov1api.BackupPhaseDeleting, false},

		// A backup that has not been reconciled yet has an empty phase. It is left
		// alone deliberately: the state is transient and the caller can retry.
		{velerov1api.BackupPhase(""), false},
	}

	for _, tc := range tests {
		t.Run(string(tc.phase), func(t *testing.T) {
			assert.Equal(t, tc.want, backupPhaseHasNoArtifacts(tc.phase))
		})
	}
}

func TestRestorePhaseHasNoArtifacts(t *testing.T) {
	tests := []struct {
		phase velerov1api.RestorePhase
		want  bool
	}{
		{velerov1api.RestorePhaseNew, true},
		{velerov1api.RestorePhaseFailedValidation, true},

		{velerov1api.RestorePhaseInProgress, false},
		{velerov1api.RestorePhaseWaitingForPluginOperations, false},
		{velerov1api.RestorePhaseWaitingForPluginOperationsPartiallyFailed, false},
		{velerov1api.RestorePhaseFinalizing, false},
		{velerov1api.RestorePhaseFinalizingPartiallyFailed, false},
		{velerov1api.RestorePhaseCompleted, false},
		{velerov1api.RestorePhasePartiallyFailed, false},
		{velerov1api.RestorePhaseFailed, false},

		{velerov1api.RestorePhase(""), false},
	}

	for _, tc := range tests {
		t.Run(string(tc.phase), func(t *testing.T) {
			assert.Equal(t, tc.want, restorePhaseHasNoArtifacts(tc.phase))
		})
	}
}

// Every phase the API defines should be covered by the tables above, so a phase added
// later fails here instead of silently falling into the permissive branch.
func TestPhaseTablesCoverAllPhases(t *testing.T) {
	backupPhases := []velerov1api.BackupPhase{
		velerov1api.BackupPhaseNew,
		velerov1api.BackupPhaseQueued,
		velerov1api.BackupPhaseReadyToStart,
		velerov1api.BackupPhaseFailedValidation,
		velerov1api.BackupPhaseInProgress,
		velerov1api.BackupPhaseWaitingForPluginOperations,
		velerov1api.BackupPhaseWaitingForPluginOperationsPartiallyFailed,
		velerov1api.BackupPhaseFinalizing,
		velerov1api.BackupPhaseFinalizingPartiallyFailed,
		velerov1api.BackupPhaseCompleted,
		velerov1api.BackupPhasePartiallyFailed,
		velerov1api.BackupPhaseFailed,
		velerov1api.BackupPhaseDeleting,
	}
	assert.Len(t, backupPhases, 13, "BackupPhase count changed; update backupPhaseHasNoArtifacts and its table")

	restorePhases := []velerov1api.RestorePhase{
		velerov1api.RestorePhaseNew,
		velerov1api.RestorePhaseFailedValidation,
		velerov1api.RestorePhaseInProgress,
		velerov1api.RestorePhaseWaitingForPluginOperations,
		velerov1api.RestorePhaseWaitingForPluginOperationsPartiallyFailed,
		velerov1api.RestorePhaseFinalizing,
		velerov1api.RestorePhaseFinalizingPartiallyFailed,
		velerov1api.RestorePhaseCompleted,
		velerov1api.RestorePhasePartiallyFailed,
		velerov1api.RestorePhaseFailed,
	}
	assert.Len(t, restorePhases, 10, "RestorePhase count changed; update restorePhaseHasNoArtifacts and its table")
}
