// © 2023 SolarWinds Worldwide, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package host

import (
	"github.com/google/uuid"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/host/azure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockedHostID(t *testing.T) {
	hostname := "example.com"
	p := 12345
	ec2Id := "ewr234ldsfl342-1"
	ec2Zone := "us-east-7"
	dockerId := "23423jlksl4j2l"
	mac := []string{"72:00:07:e5:23:51", "c6:61:8b:53:d6:b5", "72:00:07:e5:23:50"}
	herokuId := "heroku-test"
	azureAppInstId := "azure-test"
	azureMetadata := &azure.Metadata{
		Compute: azure.Compute{
			Location:          "westus",
			Name:              "my name",
			ResourceGroupName: "my resource group name",
			SubscriptionID:    "my subscription id",
			VMID:              "my vmid",
			VMScaleSetName:    "my vmscaleset name",
			VMSize:            "my vmsize",
		},
		Other: nil,
	}

	lh := newLockedID()
	assert.False(t, lh.ready())
	// try partial update
	lh.fullUpdate(withHostname(hostname))
	assert.Equal(t, "", lh.copyID().Hostname())

	lh.fullUpdate(
		withHostname(hostname),
		withPid(p), // pid doesn't change, but we fullUpdate it anyways
		withEC2Id(ec2Id),
		withEC2Zone(ec2Zone),
		withContainerId(dockerId),
		withMAC(mac),
		withHerokuId(herokuId),
		withAzureAppInstId(azureAppInstId),
		withAzureMetadata(azureMetadata),
	)

	assert.True(t, lh.ready())
	lh.setReady()

	lh.waitForReady()
	h := lh.copyID()
	assert.Equal(t, hostname, h.Hostname())
	assert.Equal(t, p, h.Pid())
	assert.Equal(t, ec2Id, h.EC2Id())
	assert.Equal(t, ec2Zone, h.EC2Zone())
	assert.Equal(t, dockerId, h.ContainerId())
	assert.Equal(t, mac, h.MAC())
	assert.EqualValues(t, herokuId, h.HerokuId())
	assert.EqualValues(t, azureAppInstId, h.AzureAppInstId())
	assert.Equal(t, azureMetadata, h.AzureMetadata())
	assert.Equal(t, instanceId, h.InstanceID())
	assert.Len(t, h.InstanceID(), 36)
	uid, err := uuid.Parse(h.InstanceID())
	assert.NoError(t, err)
	assert.Equal(t, uuid.Version(4), uid.Version())
}
