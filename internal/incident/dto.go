package incident

type GetNearbyIncidentRequest struct {
	Latitude  float64 `form:"lat"  binding:"required,min=-90,max=90"`
	Longitude float64 `form:"long" binding:"required,min=-180,max=180"`
	Radius    int     `form:"rad"    binding:"required,min=1,max=50000"`
}

type GetNearbyIncidentResponse struct {
	NearbyIncidents []Incident `json:"incidents"`
}

type GetIncidentRequest struct {
	Id string `form:"id" binding:"required"`
}

type GetIncidentResponse struct {
	Incident      Incident               `json:"incident"`
	Confirmations []IncidentConfirmation `json:"confirmations"`
}
