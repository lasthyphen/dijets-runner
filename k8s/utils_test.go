package k8s

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/lasthyphen/dijets-runner/network"
	"github.com/lasthyphen/dijets-runner/network/node"
	"github.com/lasthyphen/dijets-runner/utils"
	"github.com/lasthyphen/dijetsgo/utils/logging"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

// TestBuildNodeEnv tests the internal buildNodeEnv method which creates the env vars for the avalanche nodes
func TestBuildNodeEnv(t *testing.T) {
	genesis := defaultTestGenesis
	testConfig := `
	{
		"network-peer-list-gossip-frequency": "250ms",
		"network-max-reconnect-delay": "1s",
		"health-check-frequency": "2s"
	}`
	c := node.Config{
		ConfigFile: testConfig,
	}

	envVars, err := buildNodeEnv(logging.NoLog{}, genesis, c)
	assert.NoError(t, err)
	controlVars := []v1.EnvVar{
		{
			Name:  "AVAGO_NETWORK_PEER_LIST_GOSSIP_FREQUENCY",
			Value: "250ms",
		},
		{
			Name:  "AVAGO_NETWORK_MAX_RECONNECT_DELAY",
			Value: "1s",
		},
		{
			Name:  "AVAGO_HEALTH_CHECK_FREQUENCY",
			Value: "2s",
		},
		{
			Name:  "AVAGO_NETWORK_ID",
			Value: fmt.Sprint(defaultTestNetworkID),
		},
	}

	assert.ElementsMatch(t, envVars, controlVars)
}

// TestConvertKey tests the internal convertKey method which is used
// to convert from the avalanchego config file format to env vars
func TestConvertKey(t *testing.T) {
	testKey := "network-peer-list-gossip-frequency"
	controlKey := "AVAGO_NETWORK_PEER_LIST_GOSSIP_FREQUENCY"
	convertedKey := convertKey(testKey)
	assert.Equal(t, convertedKey, controlKey)
}

// TestCreateDeploymentConfig tests the internal createDeploymentFromConfig method which creates the k8s objects
func TestCreateDeploymentConfig(t *testing.T) {
	assert := assert.New(t)
	genesis := defaultTestGenesis

	nodeConfigs := []node.Config{
		{
			Name:        "test1",
			IsBeacon:    true,
			StakingKey:  "fooKey",
			StakingCert: "fooCert",
			ConfigFile:  "{}",
			ImplSpecificConfig: utils.NewK8sNodeConfigJsonRaw(
				"v1",
				"test11",
				"img1",
				"Avalanchego",
				"test01",
				"t1",
			),
		},
		{
			Name:        "test2",
			IsBeacon:    false,
			StakingKey:  "barKey",
			StakingCert: "barCert",
			ConfigFile:  "{}",
			ImplSpecificConfig: utils.NewK8sNodeConfigJsonRaw(
				"v2",
				"test22",
				"img2",
				"Avalanchego",
				"test02",
				"t2",
			),
		},
	}
	params := networkParams{
		conf: network.Config{
			Genesis:     string(genesis),
			NodeConfigs: nodeConfigs,
		},
	}

	beacons, nonBeacons, err := createDeploymentFromConfig(params)
	assert.NoError(err)
	assert.Len(beacons, 1)
	assert.Len(nonBeacons, 1)

	b := beacons[0]
	n := nonBeacons[0]

	assert.Equal(b.Name, "test11")
	assert.Equal(n.Name, "test22")
	assert.Equal(b.Kind, "Avalanchego")
	assert.Equal(n.Kind, "Avalanchego")
	assert.Equal(b.APIVersion, "v1")
	assert.Equal(n.APIVersion, "v2")
	assert.Equal(b.Namespace, "test01")
	assert.Equal(n.Namespace, "test02")
	assert.Equal(b.Spec.DeploymentName, "test11")
	assert.Equal(n.Spec.DeploymentName, "test22")
	assert.Equal(b.Spec.Image, "img1")
	assert.Equal(n.Spec.Image, "img2")
	assert.Equal(b.Spec.Tag, "t1")
	assert.Equal(n.Spec.Tag, "t2")
	assert.Equal(b.Spec.BootstrapperURL, "")
	assert.Equal(n.Spec.BootstrapperURL, "")
	assert.Equal(b.Spec.Env[0].Name, "AVAGO_NETWORK_ID")
	assert.Equal(n.Spec.Env[0].Name, "AVAGO_NETWORK_ID")
	assert.Equal(b.Spec.Env[0].Value, fmt.Sprint(defaultTestNetworkID))
	assert.Equal(n.Spec.Env[0].Value, fmt.Sprint(defaultTestNetworkID))
	assert.Equal(b.Spec.NodeCount, 1)
	assert.Equal(n.Spec.NodeCount, 1)
	assert.Equal(b.Spec.Certificates[0].Cert, base64.StdEncoding.EncodeToString([]byte("fooCert")))
	assert.Equal(b.Spec.Certificates[0].Key, base64.StdEncoding.EncodeToString([]byte("fooKey")))
	assert.Equal(n.Spec.Certificates[0].Cert, base64.StdEncoding.EncodeToString([]byte("barCert")))
	assert.Equal(n.Spec.Certificates[0].Key, base64.StdEncoding.EncodeToString([]byte("barKey")))
	assert.Equal(n.Spec.NodeCount, 1)
	assert.Equal(b.Spec.Genesis, string(genesis))
	assert.Equal(n.Spec.Genesis, string(genesis))
}
