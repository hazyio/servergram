package agents

import (
	_ "embed" // required when embedding into []byte or string
	"fmt"

	"hazyio.com/servergram/server/db"
)

//go:embed agent_linux_amd64.bin
var agentLinuxAMD64 []byte

//go:embed agent_linux_arm64.bin
var agentLinuxARM64 []byte

func GetAgentFor(os db.OS, arch db.Arch) ([]byte, error) {
	switch os {
	case db.OSLinux:
		switch arch {
		case db.ArchAMD64:
			return agentLinuxAMD64, nil
		case db.ArchARM64:
			return agentLinuxARM64, nil
		}
	}
	return nil, fmt.Errorf("unsupported OS/Arch combination")
}
