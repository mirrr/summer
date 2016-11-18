package summer

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_init(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		init()
	}
}

func TestPackagePath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := PackagePath(); got != tt.want {
			t.Errorf("%q. PackagePath() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_dot(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := dot(tt.args.name); got != tt.want {
			t.Errorf("%q. dot() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_jsoner(t *testing.T) {
	type args struct {
		object interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := jsoner(tt.args.object); got != tt.want {
			t.Errorf("%q. jsoner() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_postBind(t *testing.T) {
	type args struct {
		c   *gin.Context
		ret interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := postBind(tt.args.c, tt.args.ret); got != tt.want {
			t.Errorf("%q. postBind() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_indexOf(t *testing.T) {
	type args struct {
		arr interface{}
		v   interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := indexOf(tt.args.arr, tt.args.v); got != tt.want {
			t.Errorf("%q. indexOf() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_getJSON(t *testing.T) {
	type args struct {
		url    string
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := getJSON(tt.args.url, tt.args.target); (err != nil) != tt.wantErr {
			t.Errorf("%q. getJSON() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestExtend(t *testing.T) {
	type args struct {
		to   interface{}
		from interface{}
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		Extend(tt.args.to, tt.args.from)
	}
}

func TestH3hash(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := H3hash(tt.args.s); got != tt.want {
			t.Errorf("%q. H3hash() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_dummy(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		dummy(tt.args.c)
	}
}

func Test_setCookie(t *testing.T) {
	type args struct {
		c     *gin.Context
		name  string
		value string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		setCookie(tt.args.c, tt.args.name, tt.args.value)
	}
}

func TestEnv(t *testing.T) {
	type args struct {
		envName      string
		defaultValue string
	}
	tests := []struct {
		name      string
		args      args
		wantValue string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if gotValue := Env(tt.args.envName, tt.args.defaultValue); gotValue != tt.wantValue {
			t.Errorf("%q. Env() = %v, want %v", tt.name, gotValue, tt.wantValue)
		}
	}
}
