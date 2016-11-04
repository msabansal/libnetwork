package libnetwork

import "github.com/docker/libnetwork/types"

const libnGWNetwork = "nat"

func (c *controller) createGWNetwork() (Network, error) {
	return nil, types.NotImplementedErrorf("default gateway functionality is not implemented in windows")
}
