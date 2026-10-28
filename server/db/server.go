package db

type Server struct {
	ID       uint `gorm:"primaryKey"`
	GroupID  uint
	SSHKeyID uint
	Group    Group  `gorm:"foreignKey:GroupID"`
	SSHKey   SSHKey `gorm:"foreignKey:SSHKeyID"`
}
