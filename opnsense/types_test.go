package opnsense

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type SampleStruct struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type EmbededStruct struct {
	FieldStruct `json:"field"`
	Hello       string `json:"hello"`
}

type FieldStruct struct {
	OnePoint string        `json:"one_point"`
	Sample   *SampleStruct `json:"sample"`
}

func TestStructToMap_Normal(t *testing.T) {
	sample := SampleStruct{
		Name: "John Doe",
		ID:   "12121",
	}

	res := StructToMap(sample)
	require.NotNil(t, res)

	fmt.Printf("%+v \n", res)
	// Output: map[name:John Doe id:12121]
	jbyt, err := json.Marshal(res)
	require.NoError(t, err)
	fmt.Println(string(jbyt))
	// Output: {"id":"12121","name":"John Doe"}
}

func TestStructToMap_FieldStruct(t *testing.T) {
	sample := &SampleStruct{
		Name: "John Doe",
		ID:   "12121",
	}
	field := FieldStruct{
		Sample:   sample,
		OnePoint: "yuhuhuu",
	}

	res := StructToMap(field)
	require.NotNil(t, res)
	fmt.Printf("%+v \n", res)
	// Output: map[sample:0xc4200f04a0 one_point:yuhuhuu]
	jbyt, err := json.Marshal(res)
	require.NoError(t, err)
	fmt.Println(string(jbyt))
	// Output: {"one_point":"yuhuhuu","sample":{"name":"John Doe","id":"12121"}}
}

func TestStructToMap_EmbeddedStruct(t *testing.T) {
	sample := &SampleStruct{
		Name: "John Doe",
		ID:   "12121",
	}
	field := FieldStruct{
		Sample:   sample,
		OnePoint: "yuhuhuu",
	}

	embed := EmbededStruct{
		FieldStruct: field,
		Hello:       "WORLD!!!!",
	}

	res := StructToMap(embed)
	require.NotNil(t, res)
	fmt.Printf("%+v \n", res)
	// Output: map[field:map[one_point:yuhuhuu sample:0xc420106420] hello:WORLD!!!!]

	jbyt, err := json.Marshal(res)
	require.NoError(t, err)
	fmt.Println(string(jbyt))
	// Output: {"field":{"one_point":"yuhuhuu","sample":{"name":"John Doe","id":"12121"}},"hello":"WORLD!!!!"}
}

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

func TestBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		bit     Bool
		want    []byte
		wantErr bool
	}{
		{
			name:    "check true",
			bit:     true,
			want:    []byte("1"),
			wantErr: false,
		},
		{
			name:    "check false",
			bit:     false,
			want:    []byte("0"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.bit.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bool.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		bit     *Bool
		args    args
		wantErr bool
	}{
		{
			name: "check fails",
			bit:  nil,
			args: args{
				data: []byte("12"),
			},
			wantErr: true,
		},
		// These tests does not work as it does not check bit
		// as a "want" value. It is checked by extended tests
		// in bgp_test.go
		// {
		// 	name: "check '1'",
		// 	bit:  NewBoolPointer(true),
		// 	args: args{
		// 		data: []byte("1"),
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "check 'true'",
		// 	bit:  NewBoolPointer(true),
		// 	args: args{
		// 		data: []byte("true"),
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "check '0'",
		// 	bit:  NewBoolPointer(false),
		// 	args: args{
		// 		data: []byte("0"),
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "check 'false'",
		// 	bit:  NewBoolPointer(false),
		// 	args: args{
		// 		data: []byte("false"),
		// 	},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bit.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Bool.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPortRangeUnmarshal(t *testing.T) {
	type c struct {
		input    string
		expected PortRange
	}

	type wrap struct {
		Range PortRange `json:"range"`
	}

	validRanges := []c{
		{
			input: "{\"range\": \"20-20\"}",
			expected: PortRange{
				From: 20,
				To:   20,
			},
		},
		{
			input: "{\"range\": \"20-80\"}",
			expected: PortRange{
				From: 20,
				To:   80,
			},
		},
		{
			input: "{\"range\": \"1-65535\"}",
			expected: PortRange{
				From: 1,
				To:   65535,
			},
		},
	}
	invalidRanges := []string{
		"{\"range\": \"80-20\"}",
		"{\"range\": \"-20-20\"}",
		"{\"range\": \"-20\"}",
		"{\"range\": \"-80\"}",
		"{\"range\": \"-0\"}",
		"{\"range\": \"-65536\"}",
		"{\"range\": \"999999-9999999\"}",
	}

	for _, valid := range validRanges {
		w := wrap{}

		err := json.Unmarshal([]byte(valid.input), &w)
		if err != nil {
			t.Errorf("Expected valid result from %s, error: %s", valid.input, err)
		}

		if w.Range != valid.expected {
			t.Errorf("Actual does not match expected: %v vs %v", w.Range, valid.expected)
		}
	}

	for _, invalid := range invalidRanges {
		w := wrap{}

		err := json.Unmarshal([]byte(invalid), &w)
		if err == nil {
			t.Errorf("Expected error and an invalid outout %s, %s", invalid, err)
		}

		fmt.Printf("from: %d, to: %d\n", w.Range.From, w.Range.To)
	}
}

func TestPortValid(t *testing.T) {
	valids := []Port{1, 2, 3, 66, 9999, 65535}
	invalids := []Port{0, 65536, 99999999}

	for _, valid := range valids {
		if !valid.Valid() {
			t.Errorf("Expected valid port: %d", valid)
		}
	}

	for _, invalid := range invalids {
		if invalid.Valid() {
			t.Errorf("Expected invalid port: %d", invalid)
		}
	}
}

func TestPortFromString(t *testing.T) {
	valids := []string{"1", "2", "3", "66", "9999", "65535"}
	invalids := []string{"", "asdf", "7a", "0", "65536", "99999999"}

	for _, valid := range valids {
		_, err := portFromString(valid)
		if err != nil {
			t.Errorf("Expected valid port: %s, %s", valid, err)
		}
	}

	for _, invalid := range invalids {
		val, err := portFromString(invalid)
		if err == nil {
			t.Errorf("Expected invalid port: %s, returned: %d", invalid, val)
		}
	}
}

func TestPortRangeMarshal(t *testing.T) {
	pr := PortRange{
		From: 22,
		To:   23,
	}
	expected := "\"22-23\""

	actual, err := json.Marshal(&pr)
	if err != nil {
		t.Errorf("Failed to marshal: %s", err)
	}

	if string(actual) != expected {
		t.Errorf("Actual is not the same as expected, %s != %s", actual, expected)
	}
}
