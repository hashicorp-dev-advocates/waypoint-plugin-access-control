package builder

type Config struct {
	Trigger []Trigger `hcl:"trigger,block"`
	Access  []Access  `hcl:"access,block"`
}

type Trigger struct {
	Config      TriggerConfig `hcl:"config,block"`
	Environment string        `hcl:"environment,optional"`
	Approval    string        `hcl:"approval"`
}

type TriggerConfig struct {
	Platform string   `hcl:"platform,optional"`
	Event    string   `hcl:"event,optional"`
	Project  string   `hcl:"project,optional"`
	Teams    []string `hcl:"teams,optional"`
}

type Access struct {
	Database string `hcl:"database,optional"`
	Cloud    string `hcl:"cloud,optional"`
	Account  string `hcl:"account,optional"`
	Role     string `hcl:"role,optional"`
}
