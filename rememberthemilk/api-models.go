package rememberthemilk

import (
	"encoding/xml"

	"github.com/pkg/errors"
)

type ResponseInterface interface {
	HasError() error
}

type ResponseObject struct {
	XMLName xml.Name `xml:"rsp"`

	Stat  string        `xml:"stat,attr"`
	Error ResponseError `xml:"err"`
}

type ResponseError struct {
	Code string `xml:"code,attr"`
	Msg  string `xml:"msg,attr"`
}

func (ro *ResponseObject) HasError() error {
	if ro.Stat != "ok" {
		return errors.New(ro.Error.Msg)
	}
	return nil
}

type listResponse struct {
	*ResponseObject
	Lists []List `xml:"lists>list"` // Uses the nested structure to collapse these attributes
}

type TasklistResponse struct {
	*ResponseObject
	Lists []List `xml:"tasks>list"`
}

type GetTokenResponse struct {
	*ResponseObject
	Auth Auth `xml:"auth"`
}

type GetFrobResponse struct {
	*ResponseObject
	Frob string `xml:"frob"`
}

type checkTokenResponse struct {
	*ResponseObject
	Auth Auth `xml:"auth"`
}

type List struct {
	Name       string       `xml:"name,attr"`
	ID         int          `xml:"id,attr"`
	Deleted    bool         `xml:"deleted,attr"`
	Locked     bool         `xml:"locked,attr"`
	Archived   bool         `xml:"archived,attr"`
	Position   int          `xml:"position,attr"`
	Smart      bool         `xml:"smart,attr"`
	SortOrder  int          `xml:"sort_order,attr"`
	Taskseries []Taskseries `xml:"taskseries,omitempty"`
}

type User struct {
	Name     string `xml:"username,attr"`
	ID       int    `xml:"id,attr"`
	FullName string `xml:"fullname,attr"`
}

type Auth struct {
	User       User   `xml:"user,attr"`
	Token      string `xml:"token"`
	Permission string `xml:"perms"`
}

type Taskseries struct {
	ID        int    `xml:"id,attr"`
	Name      string `xml:"name,attr"`
	Tasks     []Task `xml:"task"`
	Completed string `xml:"completed,attr"`
}

type Task struct {
	ID        int    `xml:"id,attr"`
	Completed string `xml:"completed,attr"`
}
