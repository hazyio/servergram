package models

import (
	"time"

	"hazyio.com/servergram/shared/gen/proto"
)

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

}
type ListServersByGroupRequest struct {
	GroupID uint `path:"group_id" required:"true" minimum:"1" doc:"The group ID to filter servers by"`
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
	IntervalMs int  `query:"interval" default:"1000" minimum:"1000" maximum:"60000" doc:"Update interval in ms"`
	ID         uint `path:"id" minimum:"1" doc:"The ID of the server"`
}

type WatchProcessesRequest struct {
	IntervalMs int                                                 `query:"interval" default:"1000" minimum:"1000" maximum:"60000" doc:"Update interval in ms"`
	ID         uint                                                `path:"id" minimum:"1" doc:"The ID of the server"`
	Columns    []proto.WatchProcessesRequest_WatchProcessesColumns `query:"columns" required:"true" doc:"The columns to include in the process watch, comma-separated 1: memory, 2: command, 3: user, 4: cpu, 5: state, 6: uptime"`
}

type WatchTopProcessesRequest struct {
	ID         uint                   `path:"id" minimum:"1" doc:"The ID of the server"`
	SortBy     proto.TopProcessSortBy `query:"sortBy" required:"true" doc:"The field to sort top processes by"`
	Limit      uint32                 `query:"limit" required:"true" minimum:"1" doc:"The maximum number of top processes to return"`
	IntervalMs int                    `query:"interval" default:"1000" minimum:"1000" maximum:"60000" doc:"Update interval in ms"`
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
	Body *proto.GetServerResourcesResponse
}

type ServerInfoResponse struct {
	Body *proto.GetServerInfoResponse
}

type StopProcessRequest struct {
	ID   uint `path:"id" minimum:"1" doc:"The ID of the server"`
	Body struct {
		PID uint `json:"pid" required:"true" minimum:"1" doc:"The process ID to stop"`
	}
}

type SignalProcessRequest struct {
	ID   uint `path:"id" minimum:"1" doc:"The ID of the server"`
	Body struct {
		PID    uint                                                   `json:"pid" required:"true" minimum:"1" doc:"The process ID to signal"`
		SigNum proto.SignalProcessRequest_SignalProcessRequestSignals `json:"sig_num" required:"true" doc:"The signal to send to the process"`
	}
}
