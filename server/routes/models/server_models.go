package models

import "time"

type AddServerRequest struct {
	Body struct {
		Name     string `json:"name" required:"true" minLength:"1" maxLength:"100" pattern:"^[A-Za-z0-9_ .@-]+$" doc:"The name of the server"`
		Host     string `json:"host" required:"true" minLength:"1" maxLength:"100" doc:"The host of the server"`
		UserName string `json:"username" required:"true" minLength:"1" maxLength:"100" doc:"The username for the server"`
		Port     uint16 `json:"port" required:"true" minimum:"1" maximum:"65535" doc:"The port of the server"`
		GroupID  uint   `json:"group_id" required:"true" minimum:"1" doc:"The group ID of the server"`
		SSHKeyID uint   `json:"ssh_key_id" required:"true" minimum:"1" doc:"The SSH key for the server"`
	}
}

type AddServerResponse struct {
	Body struct {
		ID uint `json:"id" doc:"The ID of the newly created server"`
	}
}

type ListServersRequest struct {
	Body struct {
	}
}
type ServerItem struct {
	ID        uint      `json:"id" doc:"The ID of the server"`
	Name      string    `json:"name" doc:"The name of the server"`
	Host      string    `json:"host" doc:"The host of the server"`
	Trusted   bool      `json:"trusted" doc:"Whether the server is trusted"`
	CreatedAt time.Time `json:"created_at" doc:"The creation time of the server"`
}

type ListServersResponse struct {
	Body struct {
		Items []ServerItem `json:"items" doc:"The list of servers"`
	}
}

type UpdateServerFingerPrintRequest struct {
	ID   uint `path:"id" minimum:"1" doc:"The ID of the server"`
	Body struct {
		Fingerprint string `json:"fingerprint" required:"true" minLength:"1" maxLength:"100" doc:"The new host fingerprint for the server"`
	}
}

type WatchServerResourceRequest struct {
	ID        uint `path:"id" minimum:"1" doc:"The ID of the server"`
	Frequency uint `query:"frequency" required:"true" minimum:"1" doc:"The frequency of watching the server resource in seconds"`
}
type ServerUsageItem struct {
	Usage float64 `json:"usage" doc:"The usage statistics of the server"`
	Total float64 `json:"total" doc:"The total resources of the server"`
}
type ServerCPUUsageItem struct {
	Usage float64 `json:"usage" doc:"The CPU usage of the server"`
	Cores int     `json:"cores" doc:"The number of CPU cores"`
}
type WatchServerResourceResponse struct {
	Body struct {
		CPU    ServerCPUUsageItem `json:"cpu" doc:"The CPU usage of the server"`
		Memory ServerUsageItem    `json:"memory" doc:"The memory usage of the server"`
		Disk   ServerUsageItem    `json:"disk" doc:"The disk usage of the server"`
		Swap   *ServerUsageItem   `json:"swap" doc:"The swap usage of the server"`
	}
}
