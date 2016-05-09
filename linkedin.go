// package linkedin provides basic functionality to access the LinkedIn API
package linkedin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ErrInvalidClient = "Invalid Client"

// Client manages communication with the LinkedIn API
type Client struct {
	accessToken string
}

// Profile is a LinkedIn user's profile information
type Profile struct {
	Id                    string
	FirstName             string
	LastName              string
	MaidenName            string
	FormattedName         string
	PhoneticFirstName     string
	PhoneticLastName      string
	FormattedPhoneticName string
	Headline              string
	Industry              string
	NumConnections        int
	NumConnectionsCapped  bool
	Summary               string
	Specialties           string
	Positions             positions
	PictureURL            string
	PictureURLs           pictures
	PublicProfileURL      string
	EmailAddress          string
	Location              Location
	// APIStandardProfileRequest  string
	// SiteStandardProfileRequest string
	// CurrentShare               string
}

// Original pictures associated with a LinkedIn profile
type pictures struct {
	Total int      `json:"_total"`
	URLs  []string `json:"values"`
}

// Work positions associated with a LinkedIn profile
type positions struct {
	Total  int `json:"_total"`
	Values []Position
}

// Work position
type Position struct {
	Id        int
	IsCurrent bool
	Title     string
	Summary   string
	Location  Location
	StartDate Date
	EndDate   Date
	Company   Company
}

type Company struct {
	Id       int
	Industry string
	Name     string
	Size     string
	Type     string
}

type Location struct {
	Name    string
	Country Country
}

type Country struct {
	Name string
	Code string
}

type Date struct {
	Month int
	Year  int
}

// Returns a new Client using the provided accessToken
func NewClient(accessToken string) (c *Client, err error) {
	if accessToken == "" {
		return nil, errors.New(ErrInvalidClient)
	}
	return &Client{accessToken: accessToken}, nil
}

var allFields = []string{"first-name", "last-name", "maiden-name",
	"formatted-name", "phonetic-first-name", "phonetic-last-name", "formatted-phonetic-name",
	"headline", "current-share", "num-connections", "num-connections-capped", "summary",
	"specialties", "picture-url", "site-standard-profile-request", "api-standard-profile-request",
	"public-profile-url", "email-address", "industry", "picture-urls::(original)", "location", "positions", "id"}

// Get the profile for the user associated with the accessToken. Supply no fields
// to receive the minimal profile, or supply a list of field names per the LinkedIn API
// (https://developer.linkedin.com/docs/fields/basic-profile) to receive those fields,
// or supply "all" as the field to receive all the possible basic profile fields.
func (c *Client) BasicProfile(fields ...string) (p *Profile, err error) {
	if c == nil || c.accessToken == "" {
		return nil, errors.New(ErrInvalidClient)
	}
	fieldStr := ""
	if len(fields) == 1 && fields[0] == "all" {
		fields = allFields
	}
	for i, field := range fields {
		if i == 0 {
			fieldStr = ":("
		}
		if i > 0 {
			fieldStr += ","
		}
		fieldStr += field
		if i == len(fields)-1 {
			fieldStr += ")"
		}
	}
	profileURL := fmt.Sprintf("https://api.linkedin.com/v1/people/~%s?format=json", fieldStr)
	request, err := http.NewRequest("GET", profileURL, nil)
	request.Header.Add("Authorization", "Bearer "+c.accessToken)
	client := http.Client{}
	resp, err := client.Do(request)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	p = &Profile{}
	err = json.Unmarshal(body, p)
	return
}
