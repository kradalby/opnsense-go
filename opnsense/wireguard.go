package opnsense

func (c *Client) WireGuardRestart() (*GenericResponse, error) {
	api := "wireguard/service/restart"

	var response GenericResponse
	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) WireGuardStart() (*GenericResponse, error) {
	api := "wireguard/service/start"

	var response GenericResponse
	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) WireGuardStop() (*GenericResponse, error) {
	api := "wireguard/service/stop"

	var response GenericResponse
	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) WireGuardGetConfig() (*GenericResponse, error) {
	api := "wireguard/service/showconf"

	var response GenericResponse
	err := c.GetAndUnmarshal(api, &response)

	return &response, err
}

func (c *Client) WireGuardGetHandshake() (*GenericResponse, error) {
	api := "wireguard/service/showhandshake"

	var response GenericResponse
	err := c.GetAndUnmarshal(api, &response)

	return &response, err
}
