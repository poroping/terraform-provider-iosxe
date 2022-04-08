package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	examplesDir = "../../examples"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"iosxe": func() (*schema.Provider, error) {
		return New("acctest")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("acc-test")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = New("acc-test")()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("TF_IOSXE_HOST"); err == "" {
		t.Fatal("TF_IOSXE_HOST must be set for acceptance tests")
	}
	if err := os.Getenv("TF_IOSXE_USERNAME"); err == "" {
		t.Fatal("TF_IOSXE_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("TF_IOSXE_PASSWORD"); err == "" {
		t.Fatal("TF_IOSXE_PASSWORD must be set for acceptance tests")
	}
}

func testAccCreateResourceFromExampleStep(rName string) resource.TestStep {
	// skip test if no example is provided
	if testAccExampleResourceConfig(rName) == "" {
		return resource.TestStep{
			SkipFunc: testAccSkipTestStep,
		}
	}
	return resource.TestStep{
		Config: testAccExampleResourceConfig(rName),
		Check:  resource.ComposeTestCheckFunc(),
	}
}

func testAccExampleResourceConfig(rName string) string {
	b, err := os.ReadFile(fmt.Sprintf("%s/resources/%s/resource.tf", examplesDir, rName))
	if err != nil {
		return ""
	}
	return string(b)
}

func testAccSkipTestStep() (bool, error) {
	return true, nil
}
