package controller

// IPValidateRequest is the struct provided to validate the ip
type IPValidateRequest struct {
	IPAddress string   `json:"ip_address"`
	Countries []string `json:"countries"`
}
