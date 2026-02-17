package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestProvider_Metadata(t *testing.T) {
	p := New("1.0.0")()
	resp := &provider.MetadataResponse{}
	p.Metadata(context.Background(), provider.MetadataRequest{}, resp)

	if resp.TypeName != "localskills" {
		t.Errorf("expected type name 'localskills', got '%s'", resp.TypeName)
	}
	if resp.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", resp.Version)
	}
}

func TestProvider_Schema(t *testing.T) {
	p := New("test")()
	resp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %s", resp.Diagnostics)
	}

	if _, ok := resp.Schema.Attributes["base_url"]; !ok {
		t.Error("expected base_url attribute in schema")
	}
	if _, ok := resp.Schema.Attributes["api_token"]; !ok {
		t.Error("expected api_token attribute in schema")
	}
}

func configureProvider(t *testing.T, baseURL, apiToken string, baseURLNull, apiTokenNull bool) provider.ConfigureResponse {
	t.Helper()

	p := New("test")()

	schemaResp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

	attrTypes := map[string]tftypes.Type{
		"base_url":  tftypes.String,
		"api_token": tftypes.String,
	}

	var baseURLVal, apiTokenVal tftypes.Value
	if baseURLNull {
		baseURLVal = tftypes.NewValue(tftypes.String, nil)
	} else {
		baseURLVal = tftypes.NewValue(tftypes.String, baseURL)
	}
	if apiTokenNull {
		apiTokenVal = tftypes.NewValue(tftypes.String, nil)
	} else {
		apiTokenVal = tftypes.NewValue(tftypes.String, apiToken)
	}

	rawConfig := tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, map[string]tftypes.Value{
		"base_url":  baseURLVal,
		"api_token": apiTokenVal,
	})

	config, err := configToState(schemaResp.Schema, rawConfig)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	resp := provider.ConfigureResponse{}
	p.Configure(context.Background(), provider.ConfigureRequest{
		Config: *config,
	}, &resp)
	return resp
}

func configToState(s schema.Schema, raw tftypes.Value) (*tfsdk.Config, error) {
	fwSchema := s
	cfg := &tfsdk.Config{
		Schema: fwSchema,
		Raw:    raw,
	}
	return cfg, nil
}

func TestProvider_MissingToken(t *testing.T) {
	t.Setenv("LOCALSKILLS_API_TOKEN", "")
	resp := configureProvider(t, "", "", true, true)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for missing token")
	}

	found := false
	for _, d := range resp.Diagnostics.Errors() {
		if d.Summary() == "Missing API Token" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'Missing API Token' error, got: %s", resp.Diagnostics)
	}
}

func TestProvider_InvalidTokenPrefix(t *testing.T) {
	t.Setenv("LOCALSKILLS_API_TOKEN", "")
	resp := configureProvider(t, "", "invalid_token", true, false)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for invalid token prefix")
	}

	found := false
	for _, d := range resp.Diagnostics.Errors() {
		if d.Summary() == "Invalid API Token" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'Invalid API Token' error, got: %s", resp.Diagnostics)
	}
}

func TestProvider_EnvVarFallback(t *testing.T) {
	t.Setenv("LOCALSKILLS_API_TOKEN", "lsk_envtoken123")
	resp := configureProvider(t, "", "", true, true)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %s", resp.Diagnostics)
	}

	if resp.ResourceData == nil {
		t.Error("expected ResourceData to be set")
	}
	if resp.DataSourceData == nil {
		t.Error("expected DataSourceData to be set")
	}
}

func TestProvider_CustomBaseURL(t *testing.T) {
	t.Setenv("LOCALSKILLS_API_TOKEN", "")
	resp := configureProvider(t, "https://custom.localskills.sh", "lsk_customtoken", false, false)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %s", resp.Diagnostics)
	}

	if resp.ResourceData == nil {
		t.Error("expected ResourceData to be set")
	}
}

func TestProvider_EnvVarBaseURL(t *testing.T) {
	t.Setenv("LOCALSKILLS_BASE_URL", "https://staging.localskills.sh")
	t.Setenv("LOCALSKILLS_API_TOKEN", "lsk_test123")
	resp := configureProvider(t, "", "", true, true)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %s", resp.Diagnostics)
	}

	if resp.ResourceData == nil {
		t.Error("expected ResourceData to be set")
	}
}
