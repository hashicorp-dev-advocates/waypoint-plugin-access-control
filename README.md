# Waypoint Plugin Vault policy

This plugin for HashiCorp Waypoint creates policies for HashiCorp Vault so security teams can configure Vault roles to assist developer with incident response. The idea behind this plugin is that developers can specify which underlying pieces of infrastructure they would need access to in the event of an incident or outage. Security teams can then take these policies and create Vault roles according to their organisation's security posture.

The plugin has two main constructs:
1. **Triggers** - This specifies what should trigger access to the underlying infrastructure. This could be event driven for specific environments, manual for other environments and automatic for any scenarios that your organisation sees fit. It also specifies which teams should be granted access to the underlying infrastructure for the event driven options.
2. **Access** - This specifies what pieces of underlying infrastructure that developers will need access to in an incident scenario.

## Example usage

```hcl
project = "hello-world"

app "access" {
  deploy {
    use "access-control" {
      # Which event should trigger the access configurator?
      trigger {
        config {
          platform = "pager-duty"
          event    = "hello-world-down" # Event to listen out
          teams = [
            "hello-world-devs", # The person that deployed app
            "hello-world-ops",
          ]

        }
        environment = "production"
        approval = "alert"
      
      }

      trigger {
        config {}
        approval = "manual" # Manual, automatic or alert
        environment = "dev"
      
      }
  
      #What access should be granted

      access {
        database = "bye-world-db"
        role     = "bye-world-db-read"
      }

      access {
        database = "hello-world-db"
        role     = "hello-world-db-read"
      }

      access {
        cloud = "aws"
        role = "aws-cloud-read"
      }

    }
  }

  build {
    use "noop" {}
  }
}
```

The above example will create and write the following policies to Vault:

```hcl
# bye-world-db-read.hcl

path "bye-world-db/creds/bye-world-db-read" {
   capabilities = [
      "read"
   ]
}
```

```hcl
# hello-world-db-read.hcl

path "hello-world-db/creds/hello-world-db-read" {
   capabilities = [
      "read"
   ]
}
```

```hcl
# aws-cloud-read.hcl

path "aws-cloud/creds/aws-cloud-read" {
   capabilities = [
      "read"
   ]
}
```

## Waypoint runner dependencies

This plugin will need to be able to authenticate with Vault to write these policies so it will need the Vault address and an access token with the below policy attached:

```hcl
# waypoint-runner.hcl

path "/sys/policy/*" {
   capabilities = [
      "write"
   ]
}
```

Once you have created a Vault token with the above policy, you will need to pass the following environment variables to the Waypoint runner:
- `VAULT_ADDR`
- `VAULT_TOKEN`
- `VAULT_NAMESPACE` (optional)

The Waypoint runner can also be configured using our [Terraform Provider for Waypoint.](https://github.com/hashicorp-dev-advocates/terraform-provider-waypoint)

