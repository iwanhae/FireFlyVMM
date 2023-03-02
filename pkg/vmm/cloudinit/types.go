package cloudinit

func NewDefaultCloudConfig() CloudConfig {
	return CloudConfig{
		Hostname:         "firefly",
		DisableRoot:      false,
		PreserveHostname: false,
		SystemInfo:       SystemInfoConfig{Distro: "ubuntu"},
		Users:            []UserCoinfig{},
		GrowPartition: GrowPartitionConfig{
			Mode:    GrowPartitionMode_Auto,
			Devices: []string{"/"},
		},
	}
}

// https://cloudinit.readthedocs.io/en/latest/reference/modules.html
type CloudConfig struct {
	Hostname         string              `yaml:"hostname" json:"hostname"`
	DisableRoot      bool                `yaml:"disable_root" json:"disable_root"`
	PreserveHostname bool                `yaml:"preserve_hostname" json:"preserve_hostname"`
	SystemInfo       SystemInfoConfig    `yaml:"system_info" json:"system_info"`
	Users            []UserCoinfig       `yaml:"users" json:"users"`
	GrowPartition    GrowPartitionConfig `yaml:"growpart" json:"growpart"`
}

type SystemInfoConfig struct {
	Distro string `yaml:"distro" json:"distro"`
}

type UserCoinfig struct {
	Name              string   `yaml:"name" json:"name"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys" json:"ssh_authorized_keys"`
}

type GrowPartitionConfig struct {
	Mode    GrowPartitionMode `yaml:"mode" json:"mode"`
	Devices []string          `yaml:"devices" json:"devices"`
}

type GrowPartitionMode string

const (
	GrowPartitionMode_Auto GrowPartitionMode = "auto"
)
