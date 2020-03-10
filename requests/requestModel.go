package requests

type Request struct {
	Id int `json:"id"`
	Method `json:"method"`
	URL `json:"url"`
	Proto `json:"url"`
	ProtoMajor `json:"url"`
	ProtoMinor `json:"url"`
	Header `json:"url"`
	Body `json:"url"`
	GetBody `json:"url"`
	ContentLength `json:"url"`
	TransferEncoding `json:"url"`
	Close `json:"url"`
	Host `json:"url"`
	Form `json:"url"`
	PostForm `json:"url"`
	MultipartForm `json:"url"`
	Trailer `json:"url"`
	RemoteAddr `json:"url"`
	RequestURI `json:"url"`
	TLS `json:"url"`
	Cancel `json:"url"`
	Response `json:"url"`
}
