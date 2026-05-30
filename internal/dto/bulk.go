package dto

type BulkCreateLeadRequest struct {
	Leads []CreateLeadRequest `json:"leads" binding:"required,dive"`
}

type BulkUpdateLeadItem struct {
	ID string `json:"id" binding:"required"`
	UpdateLeadRequest
}

type BulkUpdateLeadRequest struct {
	Leads []BulkUpdateLeadItem `json:"leads" binding:"required,dive"`
}

type BulkResult struct {
	Index   int           `json:"index"`
	Success bool          `json:"success"`
	Lead    *LeadResponse `json:"lead,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type BulkResponse struct {
	Total      int          `json:"total"`
	Successful int          `json:"successful"`
	Failed     int          `json:"failed"`
	Results    []BulkResult `json:"results"`
}