project = "api"

app "access" {
  deploy {
    use "access-control" {
      trigger {
        config {
          platform = "alertmanager"
          event    = "hello-world-down"
          // project  = "api" # After waypoint v0.9.0
          teams    = [
            "hello-world-devs",
            "hello-world-ops",
          ]
        }
        environment = "production"
        approval = "alert"
      }

      trigger {
        config {}
        approval    = "manual" # Manual, automatic or alert
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