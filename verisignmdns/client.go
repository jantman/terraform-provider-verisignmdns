package verisignmdns

import (
  "log"
  "net/http"
  "crypto/tls"
  "errors"
  "fmt"
  "io/ioutil"
  "strings"
  "bytes"
  "time"
  "encoding/json"
)

type api_client struct {
  http_client           *http.Client
  base_url              string
  token                 string
  timeout               int
  debug                 bool
}

type api_response struct {
  body                    string
  resp_code               int
  location                string
}

// Make a new api client for RESTful calls
func NewAPIClient (i_token string, i_base_url string, i_debug bool, i_timeout int) (*api_client, error) {
  if i_debug {
    log.Printf("api_client.go: Constructing debug api_client\n")
  }

  if i_base_url == "" {
    return nil, errors.New("base URL must be set to construct a client")
  }

  /* Remove any trailing slashes since we will append
     to this URL with our own root-prefixed location */
  if strings.HasSuffix(i_base_url, "/") {
    i_base_url = i_base_url[:len(i_base_url)-1]
  }

  tr := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
  }

  client := api_client{
    http_client: &http.Client{
      Timeout: time.Second * time.Duration(i_timeout),
      Transport: tr,
      },
    timeout: i_timeout,
    base_url: i_base_url,
    token: i_token,
    debug: i_debug,
  }

  if i_debug {
    log.Printf("api_client.go: Constructed object:\n%s", client.toString())
  }
  return &client, nil
}

// Convert the important bits about this object to string representation
// This is useful for debugging.
func (obj *api_client) toString() string {
  var buffer bytes.Buffer
  buffer.WriteString(fmt.Sprintf("base_url: %s\n", obj.base_url))
  buffer.WriteString(fmt.Sprintf("token: %s\n", obj.token))
  buffer.WriteString(fmt.Sprintf("Timeout: %d\n", obj.timeout))
  return buffer.String()
}

func (client *api_client) get_rr (accountId string, zoneName string, resourceRecordId string) (map[string]interface{}, error) {
  var data map[string]interface{}
  path := fmt.Sprintf("/api/v1/accounts/%s/zones/%s/rr/%s", accountId, zoneName, resourceRecordId)
  resp, err := client.send_request("GET", path, "")

  if err != nil { return make(map[string]interface{}), err }

  if resp.resp_code != 200 {
    return make(map[string]interface{}), errors.New(fmt.Sprintf("Error GETting RR - HTTP %d - %s", resp.resp_code, resp.body))
  }

  err3 := json.Unmarshal([]byte(resp.body), &data)
  if err3 != nil {
    return nil, err3
  }
  return data, nil
}

func (client *api_client) send_request (method string, path string, data string) (api_response, error) {
  full_uri := client.base_url + path
  var req *http.Request
  var err error

  if client.debug {
    log.Printf("api_client.go: method='%s', path='%s', full uri (derived)='%s', data='%s'\n", method, path, full_uri, data)
  }

  buffer := bytes.NewBuffer([]byte(data))

  if data == "" {
    req, err = http.NewRequest(method, full_uri, nil)
  } else {
    req, err = http.NewRequest(method, full_uri, buffer)
  }

  if err != nil {
    log.Fatal(err)
    return api_response{}, err
  }

  if client.debug {
    log.Printf("api_client.go: Sending HTTP request to %s...\n", req.URL)
  }

  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Accept", "application/json")
  req.Header.Set("User-Agent", userAgentFormat)
  req.Header.Set("Authorization", fmt.Sprintf("token %s", client.token))

  if client.debug {
    log.Printf("api_client.go: Request headers:\n")
    for name, headers := range req.Header {
      for _, h := range headers {
       log.Printf("api_client.go:   %v: %v", name, h)
      }
    }

    log.Printf("api_client.go: BODY:\n")
    body := "<none>"
    if req.Body != nil {
      body = string(data)
    }
    log.Printf("%s\n", body)
  }

  resp, err := client.http_client.Do(req)

  if err != nil {
    log.Printf("api_client.go: Error detected: %s\n", err)
    return api_response{}, err
  }

  if client.debug {
    log.Printf("api_client.go: Response code: %d\n", resp.StatusCode)
    log.Printf("api_client.go: Response headers:\n")
    for name, headers := range resp.Header {
      for _, h := range headers {
       log.Printf("api_client.go:   %v: %v", name, h)
      }
    }
  }

  bodyBytes, err2 := ioutil.ReadAll(resp.Body)
  resp.Body.Close()

  if err2 != nil { return api_response{}, err2 }
  body := string(bodyBytes)
  if client.debug { log.Printf("api_client.go: BODY:\n%s\n", body) }

  loc, err4 := resp.Location()
  if err4 != nil {
    log.Printf("api_client.go - Error parsing resp.Location()")
  }
  if client.debug { log.Printf("api_client.go response Location header: %s", loc)}
  return api_response{
    body: body,
    resp_code: resp.StatusCode,
    location: loc.String(),
  }, nil
}
