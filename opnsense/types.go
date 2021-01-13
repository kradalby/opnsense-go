package opnsense

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	StatusSaved   = "saved"
	StatusDeleted = "deleted"
	StatusDone    = "done"
	StatusRunning = "running"
	StatusOK      = "ok"
)

var (
	ErrOpnsenseSave                              = errors.New("failed to save")
	ErrOpnsenseDelete                            = errors.New("failed to delete")
	ErrOpnsenseDone                              = errors.New("did not finish")
	ErrOpnsenseStatusNotOk                       = errors.New("status did not return ok")
	ErrOpnsenseEmptyListNotFound                 = errors.New("found empty array, most likely 404")
	ErrOpnsense500                               = errors.New("internal server error")
	ErrOpnsense401                               = errors.New("authentication failed")
	ErrOpnsenseBoolUnmarshal                     = errors.New("failed to unmarshal OPNsense bool")
	ErrOpnsenseBoolMarshal                       = errors.New("failed to marshal OPNsense bool")
	ErrOpnsenseInvalidPort                       = errors.New("port is invalid")
	ErrOpnsenseInvalidPortRange                  = errors.New("port range is invalid")
	ErrOpnsenseInvalidPortRangeToSmallerThanFrom = errors.New("port range is invalid, to smaller than from")
)

type SelectedMap map[string]Selected

// The OPNsense API returns a [] when there is no
// objects in the list of selected items. This is
// very inconvinient and this function tries to work
// around this by making the map pointer an empty map
// if the there is an empty array.
func (sm *SelectedMap) UnmarshalJSON(b []byte) error {
	*sm = SelectedMap{}

	type Alias SelectedMap

	var temp2 Alias

	err := json.Unmarshal(b, &temp2)
	if err != nil {
		var temp []string

		err := json.Unmarshal(b, &temp)
		if err != nil {
			return err
		}

		return nil
	}

	for key, value := range temp2 {
		(*sm)[key] = value
	}

	return nil
}

type Selected struct {
	Value    string `json:"value"`
	Selected int    `json:"selected"`
}

func (s *Selected) UnmarshalJSON(b []byte) error {
	*s = Selected{}

	type Alias Selected

	var temp Alias

	err := json.Unmarshal(b, &temp)
	if err != nil {
		type Selected2 struct {
			Value    string `json:"value"`
			Selected bool   `json:"selected"`
		}

		var temp2 Selected2

		err := json.Unmarshal(b, &temp2)
		if err != nil {
			return err
		}

		s.Value = temp2.Value
		if temp2.Selected {
			s.Selected = 1
		} else {
			s.Selected = 0
		}
	}

	s.Value = temp.Value
	s.Selected = temp.Selected

	return nil
}

func ListSelectedValues(m SelectedMap) []string {
	s := []string{}

	for _, value := range m {
		if value.Selected == 1 {
			s = append(s, value.Value)
		}
	}

	return s
}

func ListSelectedKeys(m SelectedMap) []string {
	s := []string{}

	for key, value := range m {
		if value.Selected == 1 {
			s = append(s, key)
		}
	}

	return s
}

type Bool bool

func (bit *Bool) UnmarshalJSON(b []byte) error {
	var txt string

	var boo bool

	err := json.Unmarshal(b, &txt)
	if err != nil {
		err2 := json.Unmarshal(b, &boo)
		if err2 != nil {
			return err
		}

		*bit = Bool(boo)

		return nil
	}

	*bit = Bool(txt == "1" || txt == "true")

	return nil
}

func (bit Bool) MarshalJSON() ([]byte, error) {
	switch bit {
	case true:
		return []byte("1"), nil
	case false:
		return []byte("0"), nil
	}

	return nil, ErrOpnsenseBoolMarshal
}

func (bit Bool) URLArgument() string {
	if bit {
		return "1"
	}

	return "0"
}

type Integer int

func (bit *Integer) UnmarshalJSON(b []byte) error {
	var txt string

	err := json.Unmarshal(b, &txt)
	if err != nil {
		return err
	}

	i, err := strconv.Atoi(txt)
	if err != nil {
		return err
	}

	*bit = Integer(i)

	return nil
}

func (bit Integer) MarshalJSON() ([]byte, error) {
	str := strconv.Itoa(int(bit))

	return []byte(str), nil
}

type Port int

func (p Port) Valid() bool {
	return p >= 1 && p <= 65535
}

func portFromString(portStr string) (Port, error) {
	number, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, err
	}

	port := Port(number)

	if !port.Valid() {
		return 0, ErrOpnsenseInvalidPort
	}

	return port, nil
}

type PortRange struct {
	From Port
	To   Port
}

func (pr *PortRange) UnmarshalJSON(b []byte) error {
	var txt string

	portRange := PortRange{}

	var singlePort Port

	err := json.Unmarshal(b, &singlePort)
	fmt.Printf("PortRange unmarsh err: %s", err)
	if err == nil {
		if singlePort.Valid() {
			portRange.From = singlePort
			portRange.To = singlePort

			*pr = portRange
		}
	}

	err = json.Unmarshal(b, &txt)
	if err != nil {
		return err
	}

	ports := strings.Split(txt, "-")

	if len(ports) != 2 {
		return ErrOpnsenseInvalidPortRange
	}

	fromStr := ports[0]
	toStr := ports[1]

	from, err := portFromString(fromStr)
	if err != nil {
		return err
	}

	to, err := portFromString(toStr)
	if err != nil {
		return err
	}

	if to < from {
		return ErrOpnsenseInvalidPortRangeToSmallerThanFrom
	}

	portRange.From = from
	portRange.To = to

	*pr = portRange

	return nil
}

func (pr *PortRange) MarshalJSON() ([]byte, error) {
	fromStr := strconv.Itoa(int(pr.From))
	toStr := strconv.Itoa(int(pr.To))

	r := fromStr + "-" + toStr

	return json.Marshal(r)
}

type NetworkOrAlias string

type Protocol string

type Interface string

type Option string
