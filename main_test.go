package summer

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		s Settings
	}
	tests := []struct {
		name string
		args args
		want *Panel
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := Init(tt.args.s); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Init() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestWait(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		Wait()
	}
}
