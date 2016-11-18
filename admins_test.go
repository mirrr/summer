package summer

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAdmins_Init(t *testing.T) {
	type args struct {
		panel *Panel
	}
	tests := []struct {
		name string
		a    *Admins
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.a.Init(tt.args.panel)
	}
}

func TestAdmins_Page(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		a    *Admins
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.a.Page(tt.args.c)
	}
}

func TestAdmins_Auth(t *testing.T) {
	type args struct {
		g *gin.RouterGroup
	}
	tests := []struct {
		name string
		a    *Admins
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.a.Auth(tt.args.g)
	}
}

func TestAdmins_Logout(t *testing.T) {
	type args struct {
		panelPath string
	}
	tests := []struct {
		name string
		a    *Admins
		args args
		want gin.HandlerFunc
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.a.Logout(tt.args.panelPath); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Admins.Logout() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAdmins_Login(t *testing.T) {
	type args struct {
		panelPath string
	}
	tests := []struct {
		name string
		a    *Admins
		args args
		want gin.HandlerFunc
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.a.Login(tt.args.panelPath); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Admins.Login() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAdmins_AddRaw(t *testing.T) {
	type args struct {
		admin AdminsStruct
	}
	tests := []struct {
		name    string
		a       *Admins
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.a.AddRaw(tt.args.admin); (err != nil) != tt.wantErr {
			t.Errorf("%q. Admins.AddRaw() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestAdmins_GetArr(t *testing.T) {
	tests := []struct {
		name       string
		a          *Admins
		wantAdmins []AdminsStruct
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if gotAdmins := tt.a.GetArr(); !reflect.DeepEqual(gotAdmins, tt.wantAdmins) {
			t.Errorf("%q. Admins.GetArr() = %v, want %v", tt.name, gotAdmins, tt.wantAdmins)
		}
	}
}
