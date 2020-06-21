package grillhttp

type MatchCondition interface{}

type EqualToCondition struct {
	Value string `json:"equalTo,omitempty"`
}

type BasicAuthCredentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Request struct {
	Method               string                    `json:"method,omitempty"`
	Url                  string                    `json:"url,omitempty"`
	UrlPath              string                    `json:"urlPath,omitempty"`
	UrlPathPattern       string                    `json:"urlPathPattern,omitempty"`
	UrlPattern           string                    `json:"urlPattern,omitempty"`
	QueryParameters      map[string]MatchCondition `json:"queryParameters,omitempty"`
	Headers              map[string]MatchCondition `json:"headers,omitempty"`
	BodyPatterns         map[string]MatchCondition `json:"bodyPatterns,omitempty"`
	Cookies              map[string]MatchCondition `json:"cookies,omitempty"`
	BasicAuthCredentials *BasicAuthCredentials     `json:"basicAuthCredentials,omitempty"`
}

type Response struct {
	Status                 int               `json:"status,omitempty"`
	StatusMessage          string            `json:"statusMessage,omitempty"`
	Headers                map[string]string `json:"headers,omitempty"`
	Body                   string            `json:"body,omitempty"`
	JsonBody               interface{}       `json:"jsonBody,omitempty"`
	FixedDelayMilliseconds int               `json:"fixedDelayMilliseconds,omitempty"`
}

type Stub struct {
	Id       string   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	Request  Request  `json:"request,omitempty"`
	Response Response `json:"response,omitempty"`
	Priority int      `json:"priority,omitempty"`
}
