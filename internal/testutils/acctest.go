package testutils

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/localskills/terraform-provider-localskills/internal/provider"
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"localskills": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("LOCALSKILLS_API_TOKEN"); v == "" {
		t.Fatal("LOCALSKILLS_API_TOKEN must be set for acceptance tests")
	}
}

func RandomName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, rand.Intn(10000))
}
