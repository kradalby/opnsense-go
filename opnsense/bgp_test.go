package opnsense

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBgpNeighborUnmarshal(t *testing.T) {
	expectedJson := `
	{
	  "neighbor": {
	    "enabled": "1",
	    "address": "10.61.0.102",
	    "remoteas": "64461",
	    "updatesource": {
	      "": {
	        "value": "none",
	        "selected": false
	      },
	      "lan": {
	        "value": "LAN",
	        "selected": 0
	      },
	      "wan": {
	        "value": "WAN",
	        "selected": 1
	      },
	      "opt1": {
	        "value": "WIREGUARD",
	        "selected": 0
	      },
	      "wireguard": {
	        "value": "WireGuard",
	        "selected": 0
	      }
	    },
	    "nexthopself": "",
	    "multihop": "0",
	    "defaultoriginate": "",
	    "linkedPrefixlistIn": {
	      "": {
	        "value": "none",
	        "selected": 0
	      }
	    },
	    "linkedPrefixlistOut": {
	      "": {
	        "value": "none",
	        "selected": 0
	      }
	    },
	    "linkedRoutemapIn": {
	      "": {
	        "value": "none",
	        "selected": 0
	      }
	    },
	    "linkedRoutemapOut": {
	      "": {
	        "value": "none",
	        "selected": 0
	      }
	    }
	  }
	}`
	type Response struct {
		Neighbor BgpNeighborGet `json:"neighbor"`
	}

	var response Response

	err := json.Unmarshal([]byte(expectedJson), &response)
	if err != nil {
		t.Errorf("Received error from JSON unmarshal, %v", err)
	}

	fmt.Println(response)
}
