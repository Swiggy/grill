package grilldp

type matchRequest struct {
	Body string `json:"body"`
}

type requestJournal struct {
	Requests []matchRequest `json:"requests"`
}

type eventHeader struct {
	Name          string `json:"name"`
	AppName       string `json:"appName"`
	SchemaVersion string `json:"schemaVersion"`
}

type event struct {
	Header *eventHeader           `json:"header"`
	Event  map[string]interface{} `json:"event"`
}

type uuidRegister struct {
	AppName string `json:"appName"`
}

type request struct {
	Method string `json:"method,omitempty"`
	Url    string `json:"url,omitempty"`
}

type response struct {
	Status int    `json:"status,omitempty"`
	Body   string `json:"body,omitempty"`
}

type stub struct {
	Request  request  `json:"request,omitempty"`
	Response response `json:"response,omitempty"`
}

var registerEventStub = stub{
	Request: request{
		Method: "POST",
		Url:    "/register",
	},
	Response: response{
		Status: 200,
		Body:   "{}",
	},
}

var messageSetStub = stub{
	Request: request{
		Method: "POST",
		Url:    "/message-set",
	},
	Response: response{
		Status: 200,
		Body:   "{}",
	},
}

var messageStub = stub{
	Request: request{
		Method: "POST",
		Url:    "/message",
	},
	Response: response{
		Status: 200,
		Body:   "{}",
	},
}
