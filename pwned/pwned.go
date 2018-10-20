package pwned

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
}

func NewClient(httpClient *http.Client, baseURL *url.URL, userAgent string) Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return Client{BaseURL: baseURL, UserAgent: userAgent, httpClient: httpClient}
}

type Hash struct {
	Value string
	Count int
}

func (p Hash) String() string {
	return fmt.Sprintf("%s:%d", p.Value, p.Count)
}

func (c Client) HashesInRange(r string) ([]Hash, error) {
	rel := &url.URL{Path: "/range/" + r}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute http request: %v", err)
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}
	rData := string(respData)
	pwnedPasswordsRaw := strings.Split(rData, "\r\n")

	p, err := parsePwnedPasswords(pwnedPasswordsRaw)
	if err != nil {
		return nil, fmt.Errorf("could not parse pwned passwords: %v", err)
	}

	return p, nil
}

func parsePwnedPasswords(pwnedPasswordsRaw []string) ([]Hash, error) {
	var pwnedPasswords []Hash

	for _, p := range pwnedPasswordsRaw {
		pp, err := parsePwnedPassword(p)
		if err != nil {
			return nil, fmt.Errorf("could not parse pwned password: %v", err)
		}
		pwnedPasswords = append(pwnedPasswords, *pp)
	}

	return pwnedPasswords, nil
}

func parsePwnedPassword(pwnedPasswordRaw string) (*Hash, error) {
	parts := strings.Split(pwnedPasswordRaw, ":")

	if len(parts) != 2 {
		return nil, fmt.Errorf("pwned password (%s) not in expected format", pwnedPasswordRaw)
	}

	count, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("could not retrive pwned password count: %v", err)
	}

	return &Hash{Value: parts[0], Count: count}, nil
}
