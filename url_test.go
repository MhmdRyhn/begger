package begger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlParts(t *testing.T) {
	port := 443
	parts := UrlParts{
		Host:        "https://example.com",
		Port:        &port,
		PathFormat:  "users/{UserId}",
		PathParams:  PathParams{"{UserId}": "123"},
		QueryParams: QueryParams{"type": "admin", "dept": "engg"},
	}
	assert.Equal(t, parts.GetUrl(), "https://example.com:443/users/123?dept=engg&type=admin")
}

func TestQueryParams(t *testing.T) {
	q := QueryParams{"type": "admin", "dept": "engg"}
	assert.Equal(t, q.ToEncodedString(), "dept=engg&type=admin")
}

func TestPathParams(t *testing.T) {
	format := "countries/{CountryId}/cities/{CityId}"
	p := PathParams{"{CountryId}": "2", "{CityId}": "4"}
	assert.Equal(t, p.ActualPath(format), "/countries/2/cities/4")
}
