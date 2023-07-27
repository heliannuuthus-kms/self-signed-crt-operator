/*
Copyright 2020 The cert-manager Authors

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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	piv1alpha1api "github.com/heliannuuthus/privateca-issuer/api/v1alpha1"
)

func TestSetReadyCondition(t *testing.T) {
	var issuerStatus piv1alpha1api.IssuerStatus

	SetReadyCondition(&issuerStatus, piv1alpha1api.ConditionTrue, "reason1", "message1")
	assert.Equal(t, "message1", GetReadyCondition(&issuerStatus).Message)

	SetReadyCondition(&issuerStatus, piv1alpha1api.ConditionFalse, "reason2", "message2")
	assert.Equal(t, "message2", GetReadyCondition(&issuerStatus).Message)
}
