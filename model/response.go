package model

type Response struct {
	Id           *int    `json:"id"                      validate:"required,number"`
	Code         *int    `json:"code"                    validate:"required,number"`
	ContentType  *string `json:"content_type,omitempty"  validate:"required_with=BodyFilename"`
	BodyFilename *string `json:"body_filename,omitempty" validate:"required_with=ContentType,file"`
	Delay        int     `json:"delay"                   validate:"number"`
}

func (r *Response) HasBody() bool {
	return r.BodyFilename != nil
}