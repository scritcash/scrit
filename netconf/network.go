package netconf

import (
	"encoding/json"
	"io/ioutil"
)

// Network defines a Scrit network.
type Network struct {
	NetworkEpochs []NetworkEpoch // global list of signing epochs
}

// Load a network configuration from filename and return the Network struct.
func Load(filename string) (*Network, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var net Network
	if err := json.Unmarshal(data, &net); err != nil {
		return nil, err
	}
	return &net, err
}

// Validate the net configuration.
func (net *Network) Validate() error {
	// validate each network epoch
	for _, e := range net.NetworkEpochs {
		if err := e.Validate(); err != nil {
			return err
		}
	}
	// validate network epoch transitions
	// TODO

	// sign start n+1 == sign end n

	// validation end n <= sign end n+1

	return nil
}

// Marshal net as string.
func (net *Network) Marshal() string {
	jsn, err := json.MarshalIndent(net, "", "  ")
	if err != nil {
		panic(err) // should never happen
	}
	return string(jsn)
}
