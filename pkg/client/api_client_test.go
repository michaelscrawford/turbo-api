package client

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/turbonomic/turbo-api/pkg/api"
	"crypto/tls"
)

func TestNewAPIClientWithBA(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost")
	secureURL, _ := url.Parse("https://localhost")
	apiPath := "path/to/api"
	table := []struct {
		config         *Config
		expectedClient *Client
		expectsError   bool
	}{
		{
			config:       &Config{baseURL, apiPath, nil},
			expectsError: true,
		},
		{
			config: &Config{baseURL, apiPath, &BasicAuthentication{"foo", "bar"}},
			expectedClient: &Client{
				&RESTClient{http.DefaultClient, baseURL, apiPath, &BasicAuthentication{"foo", "bar"}},
			},
			expectsError: false,
		},
		{
			config: &Config{secureURL, apiPath, &BasicAuthentication{"foo", "bar"}},
			expectedClient: &Client{
				&RESTClient{&http.Client{Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}}, secureURL, apiPath, &BasicAuthentication{"foo", "bar"}},
			},
			expectsError: false,
		},
	}
	for _, item := range table {
		client, err := NewAPIClientWithBA(item.config)
		if item.expectsError && err == nil {
			t.Error("Expects error, got no error")
		}
		if !reflect.DeepEqual(client, item.expectedClient) {
			t.Errorf("Expected client %++v, got %++v", item.expectedClient, client)
		}
	}
}

// Error is expected because of empty address
func TestClient_DiscoverTarget_WithError(t *testing.T) {
	address := ""
	baseURL, _ := url.Parse("http://localhost")
	apiPath := "path/to/api"
	config := &Config{baseURL, apiPath, &BasicAuthentication{"foo", "bar"}}
	client, _ := NewAPIClientWithBA(config)
	_, err := client.DiscoverTarget(address)
	if err == nil {
		t.Error("Expected error, but got no error.")
	}
}

func TestClient_AddExternalTarget_WithError(t *testing.T) {
	target := &api.ExternalTarget{}
	baseURL, _ := url.Parse("http://localhost")
	apiPath := "path/to/api"
	config := &Config{baseURL, apiPath, &BasicAuthentication{"foo", "bar"}}
	client, _ := NewAPIClientWithBA(config)
	_, err := client.AddExternalTarget(target)
	if err == nil {
		t.Error("Expected error, but got no error.")
	}
}

func TestClient_AddTarget_WithError(t *testing.T) {
	target := &api.Target{}
	baseURL, _ := url.Parse("http://localhost")
	apiPath := "path/to/api"
	config := &Config{baseURL, apiPath, &BasicAuthentication{"foo", "bar"}}
	client, _ := NewAPIClientWithBA(config)
	_, err := client.AddTarget(target)
	if err == nil {
		t.Error("Expected error, but got no error.")
	}
}
