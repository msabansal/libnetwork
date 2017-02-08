package overlay

import (
	"encoding/json"

	"github.com/Microsoft/hcsshim"
	"github.com/Sirupsen/logrus"
)

const overlayNetworkPrefix = "overlay/network"

func (d *driver) createHnsNetwork(n *network) error {

	subnets := []hcsshim.Subnet{}

	for _, s := range n.subnets {
		subnet := hcsshim.Subnet{
			AddressPrefix: s.subnetIP.String(),
		}

		if s.gwIP != nil {
			subnet.GatewayAddress = s.gwIP.String()
		}

		vsidPolicy, err := json.Marshal(hcsshim.VsidPolicy{
			Type: "VSID",
			VSID: uint(s.vni),
		})

		if err != nil {
			return err
		}

		subnet.Policies = append(subnet.Policies, vsidPolicy)
		subnets = append(subnets, subnet)
	}

	network := &hcsshim.HNSNetwork{
		Name:               n.name,
		Type:               d.Type(),
		Subnets:            subnets,
		NetworkAdapterName: n.interfaceName,
	}

	configurationb, err := json.Marshal(network)
	if err != nil {
		return err
	}

	configuration := string(configurationb)
	logrus.Infof("HNSNetwork Request =%v", configuration)

	hnsresponse, err := hcsshim.HNSNetworkRequest("POST", "", configuration)
	if err != nil {
		return err
	}

	n.hnsId = hnsresponse.Id
	n.providerAddress = hnsresponse.ManagementIP

	return nil
}
