package opnsense

import (
	"fmt"
	"log"
	"path"

	uuid "github.com/satori/go.uuid"
)

// Requires: os-wireguard

// Docs:
// https://docs.opnsense.org/development/api/plugins/firewall.html

// TODO: Save/Apply function that handles save, check if we locked out, roll back or cancel rollback

// I think apply will make the changes live, and then revert back to rollbackRevision
// after 60s if not FirewallFilterCancelRollback is called with rollbackRevision.
type ApplyStatus struct {
	Status string `json:"status"`
}

func (c *Client) FirewallFilterApply(rollbackRevision *string) error {
	api := "firewall/filter/apply"
	if rollbackRevision != nil {
		api = path.Join("firewall/filter/apply", *rollbackRevision)
	}

	var response ApplyStatus

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return err
	}

	if response.Status != "OK\n\n" {
		log.Printf("[TRACE] FirewallFilterApply response: %#v", response)

		return fmt.Errorf("FirewallFilterApply failed: %w", ErrOpnsenseStatusNotOk)
	}

	return nil
}

func (c *Client) FirewallFilterCancelRollback(rollbackRevision string) (*GenericResponse, error) {
	api := path.Join("firewall/filter/cancelRollback", rollbackRevision)

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) FirewallFilterRevert(revision string) (*GenericResponse, error) {
	api := path.Join("firewall/filter/revert", revision)

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type Savepoint struct {
	Status    string `json:"status"`
	Retention string `json:"retention"`
	Revision  string `json:"revision"`
}

// Fetch the previous revision/create a revision _before_ we do stuff?
func (c *Client) FirewallFilterSavepoint() (*Savepoint, error) {
	api := "firewall/filter/savepoint"

	var response Savepoint

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type FilterRule struct {
	UUID            *uuid.UUID     `json:"uuid,omitempty"`
	Enabled         Bool           `json:"enabled,omitempty"`
	Sequence        Integer        `json:"sequence,omitempty"`
	Action          Option         `json:"action,omitempty"`
	Quick           Bool           `json:"quick,omitempty"`
	Interface       Interface      `json:"interface,omitempty"` // InterfaceField
	Direction       Option         `json:"direction,omitempty"`
	IPProtocol      Option         `json:"ipprotocol,omitempty"`
	Protocol        Protocol       `json:"protocol,omitempty"`   // ProtocolField
	SourceNet       NetworkOrAlias `json:"source_net,omitempty"` // NetworkAliasField
	SourceNot       Bool           `json:"source_not,omitempty"`
	SourcePort      *PortRange     `json:"source_port,omitempty"`     // Custom port range type PortField
	DestinationNet  NetworkOrAlias `json:"destination_net,omitempty"` // NetworkAliasField
	DestinationNot  Bool           `json:"destination_not,omitempty"`
	DestinationPort *PortRange     `json:"destination_port,omitempty"`
	Gateway         string         `json:"gateway,omitempty"` // JsonKeyValueStoreField
	Log             Bool           `json:"log,omitempty"`
	Description     string         `json:"description,omitempty"`
}

func (c *Client) FirewallFilterRuleGet(uuid uuid.UUID) (*FilterRule, error) {
	api := path.Join("firewall/filter/getRule", uuid.String())

	type Response struct {
		Rule FilterRule `json:"rule"`
	}

	var response Response

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Rule, nil
}

func (c *Client) FirewallFilterRuleSet(rule *FilterRule) error {
	api := path.Join("firewall/filter/setRule", rule.UUID.String())

	request := map[string]interface{}{
		"rule": rule,
	}

	var response GenericResponse

	err := c.PostAndMarshal(api, request, &response)
	if err != nil {
		return err
	}

	if response.Result != StatusSaved {
		log.Printf("[TRACE] FirewallFilterRuleSet response: %#v", response)

		return fmt.Errorf("FirewallFilterRuleSet failed: %w", ErrOpnsenseSave)
	}

	return nil
}

func (c *Client) FirewallFilterRuleAdd(rule *FilterRule) error {
	api := "firewall/filter/addRule"

	var response GenericResponse

	request := map[string]interface{}{
		"rule": rule,
	}

	err := c.PostAndMarshal(api, request, &response)
	if err != nil {
		return err
	}

	fmt.Printf("Add: %#v", response)

	if response.Result != StatusSaved {
		log.Printf("[TRACE] FirewallFilterRuleAdd response: %#v", response)

		return fmt.Errorf("FirewallFilterRuleAdd failed: %w", ErrOpnsenseSave)
	}

	return nil
}

func (c *Client) FirewallFilterRuleDelete(uuid uuid.UUID) error {
	api := path.Join("firewall/filter/delRule", uuid.String())

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return err
	}

	if response.Result != StatusDeleted {
		log.Printf("[TRACE] FirewallFilterRuleDelete response: %#v", response)

		return fmt.Errorf("FirewallFilterRuleDelete failed: %w", ErrOpnsenseDelete)
	}

	return nil
}

func (c *Client) FirewallFilterRuleSearch() ([]*FilterRule, error) {
	api := "firewall/filter/searchRule"

	type SearchResult struct {
		Rows []*FilterRule `json:"rows"`
		SearchResultPartial
	}

	var response SearchResult

	err := c.GetAndUnmarshal(api, &response)
	if err != nil {
		return nil, err
	}

	return response.Rows, nil
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

func (c *Client) FirewallSourceNatRuleSearch() (*GenericResponse, error) {
	api := "firewall/source_nat/searchRule"

	var response GenericResponse

	err := c.GetAndUnmarshal(api, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
