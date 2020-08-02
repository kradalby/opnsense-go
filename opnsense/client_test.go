package opnsense

import (
	"testing"
)

func TestSelectedMap_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}

	tests := []struct {
		name    string
		sm      *SelectedMap
		args    args
		wantErr bool
	}{
		{
			name: "no special",
			sm:   &SelectedMap{},
			args: args{b: []byte(`{
	    					  "lan": {
	    					    "value": "LAN",
	    					    "selected": 0
	    					  },
	    					  "wan": {
	    					    "value": "WAN",
	    					    "selected": 1
	    					  }
	    					}`)},
			wantErr: false,
		},
		{
			name: "with selected as bool",
			sm:   &SelectedMap{},
			args: args{b: []byte(`{
	    					  "lan": {
	    					    "value": "LAN",
	    					    "selected": 0
	    					  },
	    					  "wan": {
	    					    "value": "WAN",
	    					    "selected": 0
	    					  },
	    					  "wan2": {
	    					    "value": "WAN",
	    					    "selected": true
	    					  }
	    					}`)},
			wantErr: false,
		},
		{
			name: "with no key for none",
			sm:   &SelectedMap{},
			args: args{b: []byte(`{
						  "": {
	      					    "value": "none",
	      					    "selected": false
	      					  },
	    					  "lan": {
	    					    "value": "LAN",
	    					    "selected": 0
	    					  },
	    					  "wan2": {
	    					    "value": "WAN",
	    					    "selected": true
	    					  }
	    					}`)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sm.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("SelectedMap.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSelected_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Value    string
		Selected int
	}

	type args struct {
		b []byte
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "pass with Selected",
			fields: fields{
				Value:    "",
				Selected: 0,
			},
			args: args{
				b: []byte(`{
	      				     "value": "LAN",
	      				     "selected": 0
	      				   }`),
			},
			wantErr: false,
		},
		{
			name: "pass with Selected2 (selected is bool)",
			fields: fields{
				Value:    "",
				Selected: 0,
			},
			args: args{
				b: []byte(`{
	      				     "value": "LAN",
	      				     "selected": false
	      				   }`),
			},
			wantErr: false,
		},
		{
			name: "fail with string",
			fields: fields{
				Value:    "",
				Selected: 0,
			},
			args: args{
				b: []byte(`{
	      				     "value": "LAN",
	      				     "selected": "fail"
	      				   }`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Selected{
				Value:    tt.fields.Value,
				Selected: tt.fields.Selected,
			}
			if err := s.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Selected.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
