package begger

import (
	"fmt"
	"net/url"
	"strings"
)

type Url struct {
	Actual *string
	Parts  *UrlParts
}

func (u *Url) Get() string {
	if u.Actual != nil {
		return *u.Actual
	} else if u.Parts != nil {
		return u.Parts.GetUrl()
	}
	panic("Either the full url or the url parts must be supplied.")
}

type UrlParts struct {
	Host        string
	Port        *int
	PathFormat  string
	PathParams  PathParams
	QueryParams QueryParams
}

func (u *UrlParts) GetUrl() string {
	url := u.Host
	if u.Port != nil && *u.Port > 0 {
		url += fmt.Sprintf(":%d", *u.Port)
	}
	url += u.PathParams.ActualPath(u.PathFormat)
	if qs := u.QueryParams.ToEncodedString(); qs != "" {
		url += "?" + qs
	}
	return url
}

type QueryParams map[string]string

func (q *QueryParams) ToEncodedString() string {
	params := url.Values{}
	for key, val := range *q {
		params.Add(key, val)
	}
	return params.Encode()
}

type PathParams map[string]string

/*
	Make sure to use the path param's placeholder structure as the map key.
	For example,
	- If pathFormat uses {id}, then it must be PathParams{"{id}": 123}
	- If pathFormat uses :id, then it must be PathParams{":id": 123}

	** NOTE: This method will make sure that the actual path will contain
	a leading slash (/). For example, if the pathFormat is either "users/:id"
	or "/users/:id", the return value will always be like "/users/123".
*/
func (p *PathParams) ActualPath(pathFormat string) string {
	var oldNew []string
	for key, val := range *p {
		oldNew = append(oldNew, key, val)
	}
	return "/" + strings.NewReplacer(oldNew...).Replace(strings.TrimLeft(pathFormat, "/"))
}
