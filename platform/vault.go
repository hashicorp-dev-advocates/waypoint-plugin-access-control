package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

type Payload struct {
	Policy string `json:"policy"`
}

func (p *Platform) policyWriter() (status string) {

	for i, t := range p.config.Access {

		var mountPoint string

		switch t.Database != "" {
		case true:
			mountPoint = t.Database
		}

		switch t.Cloud != "" {
		case true:
			mountPoint = t.Cloud
		}

		role := t.Role

		policy := `path "` + mountPoint + `/creds/` + role + `" {capabilities = ["read"] }`
		newPol := Payload{
			Policy: policy,
		}

		var jsonData []byte
		jsonData, err := json.Marshal(newPol)
		if err != nil {
			fmt.Printf(err.Error())
		}

		vaultAddr := os.Getenv("VAULT_ADDR")
		if vaultAddr == "" {
			fmt.Errorf("VAULT_ADDR environment variable must be set on the runner")
		}

		vaultToken := os.Getenv("VAULT_TOKEN")
		if vaultToken == "" {
			fmt.Errorf("VAULT_TOKEN environment variable must be set on the runner")
		}

		vaultEndpoint := vaultAddr + "/v1/sys/policy/" + p.config.Access[i].Role

		request, err := http.NewRequest("POST", vaultEndpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}

		request.Header.Set("X-Vault-Token", vaultToken)

		vaultNamespace := os.Getenv("VAULT_NAMESPACE")
		if vaultNamespace != "" {
			request.Header.Set("X-Vault-Namespace", vaultNamespace)
		}

		resp, err := Client.Do(request)
		if err != nil {
			log.Fatal(err)
		}

		_, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

	}

	status = "Vault mount point created successfully"

	return status
}
