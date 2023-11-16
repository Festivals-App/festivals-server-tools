package servertools_test

import (
	"net/http"
	"reflect"
	"testing"

	servertools "github.com/Festivals-App/festivals-server-tools"
)

func TestLoadServerCertificateHandler(t *testing.T) {

	/*
		handler := festivalspki.LoadServerCertificateHandler("certificates/*.festivalsapp.dev.crt", "certificates/*.festivalsapp.dev.key", "certificates/festivalsapp-development-root-ca.crt")
		_, err := handler(&tls.ClientHelloInfo{})
		if err != nil {
			t.Errorf("Handler failed to load server certificates.")
		}*/

	t.Log("running TestLoadX509Certificate")
}

func TestIsAuthenticated(t *testing.T) {
	type args struct {
		keys     []string
		endpoint func(http.ResponseWriter, *http.Request)
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := servertools.IsEntitled(tt.args.keys, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsAuthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}
