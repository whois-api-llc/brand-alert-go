package brandalert

import (
	"encoding/json"
	"testing"
)

// TestTime tests the time conversion functions.
func TestTime(t *testing.T) {
	tests := []struct {
		name   string
		decErr string
		encErr string
	}{
		{
			name:   `"2006-01-02"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"2006-01-02T15:04:05Z08:00"`,
			decErr: `parsing time "2006-01-02T15:04:05Z08:00": extra text: "T15:04:05Z08:00"`,
			encErr: "",
		},
		{
			name:   `""`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Time

			err := json.Unmarshal([]byte(tt.name), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.name {
				t.Errorf("got = %v, want %v", string(bb), tt.name)
			}
		})
	}
}

// TestMessages tests the Messages conversion functions.
func TestMessages(t *testing.T) {
	tests := []struct {
		name   string
		want   string
		decErr string
		encErr string
	}{
		{
			name:   `["Message1","Message2"]`,
			want:   `["Message1","Message2"]`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"Message"`,
			want:   `["Message"]`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `""`,
			want:   `[""]`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `{"key1":["Message1"],"key2":["Message2"]}`,
			want:   `["map[key1:[Message1] key2:[Message2]]"]`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `{"key1":["Message1"],"key2":"Message2"}`,
			want:   `["map[key1:[Message1] key2:Message2]"]`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Messages

			err := json.Unmarshal([]byte(tt.name), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.want {
				t.Errorf("got = %v, want %v", string(bb), tt.want)
			}
		})
	}
}

func checkErr(t *testing.T, err error, want string) {
	if (err != nil || want != "") && (err == nil || err.Error() != want) {
		t.Errorf("error = %v, wantErr %v", err, want)
	}
}
