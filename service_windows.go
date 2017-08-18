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
		logrus.Debugf("Found ednpoint %s", eid)

		driver, err := n.driver(true)
		if err != nil {
			logrus.Debugf("Error happend %v", err)
			continue
		}

		data, err := driver.EndpointOperInfo(n.ID(), eid)

		if err != nil {
			logrus.Debugf("Error happend %s", err)
			continue
		}

		if data["hnsid"] != nil {
			logrus.Debugf("Data hnsid %s", data["hnsid"])
			epList = append(epList, data["hnsid"].(string))
		} else {
			logrus.Debugf("hnsid not found")
		}
	}

	if lb.policyList != nil {
		lb.policyList.Delete()
		lb.policyList = nil
	}

	var elbPolicies []hcsshim.ELBPolicy
	var sourceVIP string

	c := n.ctrlr

	c.Lock()
	sb, ok := c.sandboxes[n.id]
	c.Unlock()
	logrus.Debugf("Finding lb backend")
	// If the load balancer sandbox is there
	if ok {
		logrus.Debugf("Found lb backend %v", sb)
		for _, ep := range sb.getConnectedEndpoints() {
			logrus.Debugf("Lb backend has connected ips")
			if ep.getNetwork().ID() == n.id {
				if ip := ep.getFirstInterfaceAddress(); ip != nil {
					sourceVIP = ip.String()
					logrus.Debugf("Lb backend found source VIP %s", sourceVIP)
				}
			}
		}
	}
	for _, port := range ingressPorts {

		elbPolicy := hcsshim.ELBPolicy{
			//VIPs:      []string{vip.String()},
			SourceVIP: sourceVIP,
			ILB:       false,
		}

		elbPolicy.Type = hcsshim.ExternalLoadBalancer
		elbPolicy.InternalPort = uint16(port.TargetPort)
		elbPolicy.ExternalPort = uint16(port.PublishedPort)

		elbPolicies = append(elbPolicies, elbPolicy)
	}

	elbPolicy := hcsshim.ELBPolicy{
		VIPs:      []string{vip.String()},
		SourceVIP: sourceVIP,
		ILB:       true,
	}

	elbPolicy.Type = hcsshim.ExternalLoadBalancer
	elbPolicies = append(elbPolicies, elbPolicy)

	if len(elbPolicies) > 0 {
		lb.policyList, _ = hcsshim.AddLoadBalancer(epList, true, vip.String(), elbPolicies)
	}
}

func (n *network) rmLBBackend(ip, vip net.IP, lb *loadBalancer, ingressPorts []*PortConfig, rmService bool) {
	logrus.Debugf("Removing lb backend %v %v", ip, vip)

	if rmService {
		if lb.policyList != nil {
			lb.policyList.Delete()
			lb.policyList = nil
		}
	} else {
		n.addLBBackend(ip, vip, lb, ingressPorts)
	}
}

func (sb *sandbox) populateLoadbalancers(ep *endpoint) {
	logrus.Debugf("Populating lb for ep %v %v", ep)
}

func arrangeIngressFilterRule() {
}
