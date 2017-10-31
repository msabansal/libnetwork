// +build !windows

package libnetwork

import "github.com/docker/libnetwork/osl"

func (sb *sandbox) createKey() error {
	if sb.config.useDefaultSandBox {
		sb.key = osl.GenerateKey("default")
	}
	sb.key = osl.GenerateKey(sb.id)
	return nil
}

func (sb *sandbox) deleteKey() error {
	return nil
}
