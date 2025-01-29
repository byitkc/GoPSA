package gopsa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Credentials represents the credentiasl that are needed for the PSA client
// connectivity.
type Credentials struct {
	ClientID   string
	PublicKey  string
	PrivateKey string
}

// ConnectionParameters represents the configuration and some basic URLs that
// are used to connect to ConnectWise PSA.
type ConnectionParameters struct {
	Protocol   string
	CompanyID  string
	Host       string
	URLBase    string
	APIURLBase string
}

// Config is a representation of the configuration that is needed for the PSA
// client connectivity.
type Config struct {
	Credentials Credentials
	Parameters  ConnectionParameters
	Site        string
}

// CompanyInfoOutput is the output from the Company Information query that
// includes some information on the API used to connect to ConnectWise
// PSA along with a few other underlying items.
type CompanyInfoOutput struct {
	CompanyName    string   `json:"CompanyName"`
	Codebase       string   `json:"Codebase"`
	VersionCode    string   `json:"VersionCode"`
	VersionNumber  string   `json:"VersionNumber"`
	CompanyID      string   `json:"CompanyID"`
	IsCloud        bool     `json:"IsCloud"`
	SiteURL        string   `json:"SiteUrl"`
	Region         string   `json:"Region"`
	AllowedOrigins []string `json:"AllowedOrigins"`
}

// GetCompanyInfo is the first step to connect to the ConnectWise PSA. It
// reaches out to ConnectWise at the main site URL and returns information about
// the connection as a CompanyInfoOutput struct.
func GetCompanyInfo(c *http.Client, protocol, host, companyID string) (CompanyInfoOutput, error) {
	url := fmt.Sprintf("%s://%s/login/companyinfo/%s", protocol, host, companyID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CompanyInfoOutput{}, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return CompanyInfoOutput{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return CompanyInfoOutput{}, err
	}
	defer resp.Body.Close()
	bCompanyInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		return CompanyInfoOutput{}, err
	}
	var companyInfo CompanyInfoOutput
	if err := json.Unmarshal(bCompanyInfo, &companyInfo); err != nil {
		return CompanyInfoOutput{}, err
	}

	return companyInfo, nil
}
