package libnetwork

import (
	"net"

	"github.com/Microsoft/hcsshim"
	"github.com/Sirupsen/logrus"
)

func (n *network) addLBBackend(ip, vip net.IP, lb *loadBalancer, ingressPorts []*PortConfig) {
	logrus.Debugf("Adding lb backend %v %v portconfig %v", ip, vip, ingressPorts)
	// n.WalkEndpoints(func(e Endpoint) bool {
	// 	ep := e.(*endpoint)
	// 	if sb, ok := ep.getSandbox(); ok {
	// 		if !sb.isEndpointPopulated(ep) {
	// 			return false
	// 		}

	// 		var gwIP net.IP
	// 		if ep := sb.getGatewayEndpoint(); ep != nil {
	// 			gwIP = ep.Iface().Address().IP
	// 		}

	// 		sb.addLBBackend(ip, vip, fwMark, ingressPorts, ep.Iface().Address(), gwIP, n.ingress)
	// 	}

	// 	return false
	// })

	var epList []string

	for eid, _ := range lb.backEnds {
		ep, err := n.EndpointByID(eid)
		if err != nil {
			continue
		}
		data, err := ep.DriverInfo()
		if err != nil {
			continue
		}

		if data["hnsid"] != nil {
			epList = append(epList, data["hnsid"].(string))
		}
	}

	if lb.policyList != nil {
		lb.policyList.Delete()
		lb.policyList = nil
	}

	var elbPolicies []hcsshim.ELBPolicy

	for _, port := range ingressPorts {

		elbPolicy := hcsshim.ELBPolicy{
			VIPs: []string{vip.String()},
			ILB:  true,
		}

		elbPolicy.Type = hcsshim.ExternalLoadBalancer
		elbPolicy.InternalPort = uint16(port.PublishedPort)
		elbPolicy.ExternalPort = uint16(port.TargetPort)

		elbPolicies = append(elbPolicies, elbPolicy)
	}

	lb.policyList, _ = hcsshim.AddLoadBalancer(epList, true, vip.String(), elbPolicies)

}

func (n *network) rmLBBackend(ip, vip net.IP, fwMark uint32, ingressPorts []*PortConfig, rmService bool) {
	logrus.Debugf("Removing lb backend %v %v", ip, vip)
}

func (sb *sandbox) populateLoadbalancers(ep *endpoint) {
	logrus.Debugf("Populating lb for ep %v %v", ep)
}

func arrangeIngressFilterRule() {
}
