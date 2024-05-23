package response

type WebResponse struct {
	Token      string      `json:"token"`
	Status     int         `json:"-"`
	Message    string      `json:"message"`
	Error      interface{} `json:"-"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
}
