package response

type WebResponse struct {
	Token      string      `json:"token,omitempty"`
	Status     int         `json:"-"`
	Message    string      `json:"message"`
	Error      interface{} `json:"-"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
	RawData    interface{} `json:"-"`
}

type ListResponse struct {
	Token   string      `json:"token"`
	Status  int         `json:"-"`
	Message string      `json:"message"`
	Error   interface{} `json:"-"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}
