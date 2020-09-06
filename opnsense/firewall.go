package opnsense

import (
	"path"

	uuid "github.com/satori/go.uuid"
)

// Requires: os-wireguard

// Docs:
// https://docs.opnsense.org/development/api/plugins/firewall.html

// TODO: Save/Apply function that handles save, check if we locked out, roll back or cancel rollback

// TODO: argument $rollback_revision=null.
func (c *Client) FirewallFilterApply() (*GenericResponse, error) {
	api := "firewall/filter/apply"

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// TODO: argument $rollback_revision.
func (c *Client) FirewallFilterCancelRollback() (*GenericResponse, error) {
	api := "firewall/filter/cancelRollback"

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// TODO: argument $revision.
func (c *Client) FirewallFilterRevert() (*GenericResponse, error) {
	api := "firewall/filter/revert"

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) FirewallFilterSavepoint() (*GenericResponse, error) {
	api := "firewall/filter/savepoint"

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) FirewallFilterRuleAdd() (*GenericResponse, error) {
	api := "firewall/filter/addRule"

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) FirewallFilterRuleDelete(uuid uuid.UUID) (*GenericResponse, error) {
	api := path.Join("firewall/filter/delRule", uuid.String())

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) FirewallFilterRuleGet(uuid uuid.UUID) (*GenericResponse, error) {
	api := path.Join("firewall/filter/getRule", uuid.String())

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) FirewallFilterRuleSearch() (*GenericResponse, error) {
	api := "firewall/filter/searchRule"

	var response GenericResponse

	err := c.GetAndUnmarshal(api, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// TODO: some sort of payload
func (c *Client) FirewallFilterRuleSet(uuid uuid.UUID) (*GenericResponse, error) {
	api := path.Join("firewall/filter/setRule", uuid.String())

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) FirewallFilterRuleToggle(uuid uuid.UUID, enabled Bool) (*GenericResponse, error) {
	api := path.Join("firewall/filter/toggleRule", uuid.String(), enabled.URLArgument())

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
