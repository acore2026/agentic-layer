package models

type ServiceClass string

const (
	ServiceClassBronze   ServiceClass = "BRONZE"
	ServiceClassSilver   ServiceClass = "SILVER"
	ServiceClassGold     ServiceClass = "GOLD"
	ServiceClassPlatinum ServiceClass = "PLATINUM"
)

type DeviceContainer struct {
	EnergyAvailabilityStatus int    `json:"energy_availability_status,omitempty"` // 0-100
	MobilityState           string `json:"mobility_state,omitempty"`            // e.g., "STATIONARY", "LOW_SPEED"
	ComputeResourceType      string `json:"compute_resource_type,omitempty"`
}

type NetworkContainer struct {
	NetworkLocality string   `json:"network_locality,omitempty"`
	ServiceArea     []string `json:"service_area,omitempty"`
	NFLoadStatus    int      `json:"nf_load_status,omitempty"`
}

type AppContainer struct {
	AppServiceCategory string `json:"app_service_category,omitempty"`
	MinBandwidthReq    string `json:"min_bandwidth_req,omitempty"`
	TransportLatencyReq int    `json:"transport_latency_req,omitempty"`
}

type SkillProfile struct {
	SkillID           string            `json:"skill_id"`
	Description       string            `json:"description"`
	EntityType        string            `json:"entity_type"` // "UE", "NF", "AF"
	ServiceClass      ServiceClass      `json:"service_class"`
	AgenticServiceURI string            `json:"agentic_service_uri"`
	Device            *DeviceContainer  `json:"device,omitempty"`
	Network           *NetworkContainer `json:"network,omitempty"`
	App               *AppContainer     `json:"app,omitempty"`
}
