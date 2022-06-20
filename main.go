package main

import (
	"log"
	"os"

	"github.com/hashicorp-dev-advocates/waypoint-plugin-access-control/builder"
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

func main() {
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		log.Fatal("VAULT_ADDR environment variable must be set on the runner")
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		log.Fatal("VAULT_TOKEN environment variable must be set on the runner")
	}

	vaultNamespace := os.Getenv("VAULT_NAMESPACE")

	grafanaAddr := os.Getenv("GRAFANA_ADDR")
	if vaultAddr == "" {
		log.Fatal("GRAFANA_ADDR environment variable must be set on the runner")
	}

	grafanaDatasource := os.Getenv("GRAFANA_DATASOURCE")
	if vaultAddr == "" {
		log.Fatal("GRAFANA_DATASOURCE environment variable must be set on the runner")
	}

	sdk.Main(sdk.WithComponents(
		&builder.Builder{
			VaultAddr:         vaultAddr,
			VaultToken:        vaultToken,
			VaultNamespace:    vaultNamespace,
			GrafanaAddr:       grafanaAddr,
			GrafanaDatasource: grafanaDatasource,
		},
	))
}
