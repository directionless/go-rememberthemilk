package rememberthemilk

import "encoding/xml"

type ListResponse struct {
	XMLName xml.Name `xml:"rsp"`

	Stat  string `xml:"stat,attr"`
	Lists []List `xml:"lists>list"` // Uses the nested structure to
	// collapse this a little.
}

type GetTokenResponse struct {
	XMLName xml.Name `xml:"rsp"`

	Stat string `xml:"stat,attr"`
	Auth Auth   `xml:"auth"`
}

type GetFrobResponse struct {
	XMLName xml.Name `xml:"rsp"`

	Stat string `xml:"stat,attr"`
	Frob string `xml:"frob"`
}

type List struct {
	Name      string `xml:"name,attr"`
	ID        int    `xml:"id,attr"`
	Deleted   bool   `xml:"deleted,attr"`
	Locked    bool   `xml:"locked,attr"`
	Archived  bool   `xml:"archived,attr"`
	Position  int    `xml:"position,attr"`
	Smart     bool   `xml:"smart,attr"`
	SortOrder int    `xml:"sort_order,attr"`
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
