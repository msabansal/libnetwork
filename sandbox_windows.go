package libnetwork

import "github.com/Microsoft/hcsshim"

func (sb *sandbox) createKey() error {
	var name string
	if sb.config.useDefaultSandBox {
		name = "default"
	} else {
		name = sb.id
	}
	compartment := &hcsshim.Compartment{
		Name: name,
	}
	compartment, err := compartment.Create()
	if err == nil {
		sb.key = compartment.ID
	}
	return err
}

func (sb *sandbox) deleteKey() error {
	if sb.key == "" {
		return nil
	}

	compartment := &hcsshim.Compartment{
		ID: sb.key,
	}

	compartment, err := compartment.Delete()
	if err == nil {
		sb.key = ""
	}
	return err
}
