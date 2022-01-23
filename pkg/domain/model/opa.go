package model

type ActionInquiryInput struct {
	Alert  *Alert  `json:"alert"`
	Job    *Job    `json:"job"`
	Action *Action `json:"action"`
}

type ActionInquiryResult struct {
	Cancel bool         `json:"cancel"`
	Args   []*Attribute `json:"args"`
}

type AlertInquiryInput struct {
	Alert *Alert `json:"alert"`
}

type AlertInquiryResult struct {
	Severity string `json:"severity"`
	Status   string `json:"status"`
}
