package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Payload struct {
	Policy string `json:"policy"`
}

func (b *Builder) CreatePolicy(access Access) error {
	var mount string
	if access.Database != "" {
		mount = access.Database
	} else if access.Cloud != "" {
		mount = access.Cloud
	} else {
		return fmt.Errorf("no access type defined")
	}

	payload := Payload{
		Policy: fmt.Sprintf(`path %s/creds/%s {capabilities = ["read"]}`, mount, access.Role),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/v1/sys/policy/%s", b.VaultAddr, access.Role)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	request.Header.Set("X-Vault-Token", b.VaultToken)
	request.Header.Set("Content-Type", "application/json")

	if b.VaultNamespace != "" {
		request.Header.Set("X-Vault-Namespace", b.VaultNamespace)
	}

	client := http.Client{
		Timeout: time.Duration(30) * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf(string(body))
	}

	return nil
}
