package db

import "time"

type OS string
type Arch string

const (
	OSLinux   OS = "linux"
	OSWindows OS = "windows"
	OSDarwin  OS = "darwin"

	ArchAMD64 Arch = "amd64"
	ArchARM64 Arch = "arm64"
)

type Server struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:100;not null"`
	Host     string `gorm:"size:255;not null;uniqueIndex:idx_server_target,priority:1"`
	Port     uint16 `gorm:"not null;default:22;uniqueIndex:idx_server_target,priority:3"`
	Username string `gorm:"size:64;not null;uniqueIndex:idx_server_target,priority:2"`

	OS   OS   `gorm:"not null;size:20;check:chk_servers_os,os IN ('linux','windows','darwin')"`
	Arch Arch `gorm:"not null;size:20;check:chk_servers_arch,arch IN ('amd64','arm64')"`

	GroupID            uint   `gorm:"not null;index"`
	SSHKeyID           uint   `gorm:"not null;index"`
	Group              Group  `gorm:"foreignKey:GroupID"`
	SSHKey             SSHKey `gorm:"foreignKey:SSHKeyID"`
	HostKeyFingerprint string `gorm:"size:100"`
	Trusted            bool   `gorm:"not null;default:false"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
