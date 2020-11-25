package dockerhub

import (
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type SearchImageResponse struct {
	PageSize  int              `json:"page_size"`
	Next      string           `json:"next"`
	Previous  string           `json:"previous"`
	Page      int              `json:"page"`
	Count     int              `json:"count"`
	Summaries []ImageSummaries `json:"summaries"`
}

type ImagePublisher struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ImageCategories struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type ImageOperatingSystems struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type ImageArchitectures struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type ImageLogoURL struct {
	Small   string `json:"small"`
	Small2X string `json:"small@2x"`
}

type ImageSummaries struct {
	ID                  string                  `json:"id"`
	Name                string                  `json:"name"`
	Slug                string                  `json:"slug"`
	Type                string                  `json:"type"`
	Publisher           ImagePublisher          `json:"publisher"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
	ShortDescription    string                  `json:"short_description"`
	Source              string                  `json:"source"`
	Popularity          interface{}             `json:"popularity"`
	Categories          []ImageCategories       `json:"categories"`
	OperatingSystems    []ImageOperatingSystems `json:"operating_systems"`
	Architectures       []ImageArchitectures    `json:"architectures"`
	LogoURL             ImageLogoURL            `json:"logo_url"`
	CertificationStatus string                  `json:"certification_status"`
	StarCount           int                     `json:"star_count"`
	FilterType          string                  `json:"filter_type"`
}

func SearchImage(q string, page, pageSize int) (*SearchImageResponse, error) {
	u, err := url.Parse(searchImageURL)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("q", q)
	query.Set("page", strconv.Itoa(page))
	query.Set("page_size", strconv.Itoa(pageSize))

	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	searchImageResp := &SearchImageResponse{}
	if err := json.Unmarshal(body, searchImageResp); err != nil {
		return nil, err
	}
	return searchImageResp, nil
}
