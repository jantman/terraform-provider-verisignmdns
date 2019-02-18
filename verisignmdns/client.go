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
  account_id            string
  zone_name             string
  timeout               int
  debug                 bool
}

type api_response struct {
  body                    string
  resp_code               int
  location                string
}

// Make a new api client for RESTful calls
func NewAPIClient (i_token string, i_base_url string, i_account_id string, i_zone_name string, i_debug bool, i_timeout int) (*api_client, error) {
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
    account_id: i_account_id,
    zone_name: i_zone_name,
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
  buffer.WriteString(fmt.Sprintf("timeout: %d\n", obj.timeout))
  buffer.WriteString(fmt.Sprintf("accound_id: %s\n", obj.account_id))
  buffer.WriteString(fmt.Sprintf("zone_name: %s\n", obj.zone_name))
  return buffer.String()
}

func (client *api_client) get_rr (resourceRecordId string) (map[string]interface{}, error) {
  var data map[string]interface{}
  path := fmt.Sprintf("%s/api/v1/accounts/%s/zones/%s/rr/%s", client.base_url, client.account_id, client.zone_name, resourceRecordId)
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

func (client *api_client) delete_rr (resourceRecordId string) (error) {
  path := fmt.Sprintf("%s/api/v1/accounts/%s/zones/%s/rr/%s", client.base_url, client.account_id, client.zone_name, resourceRecordId)

  type NewRr struct {
    Comments string `json:"comments"`
  }

  recData := NewRr{
    Comments: "deleted by terraform-provider-verisignmdns",
  }

  sendData, err2 := json.Marshal(recData)
  if err2 != nil {
    return err2
  }

  resp, err := client.send_request("DELETE", path, string(sendData))

  if err != nil { return err }

  if resp.resp_code != 204 {
    return errors.New(fmt.Sprintf("Error DELETEing RR - HTTP %d - %s", resp.resp_code, resp.body))
  }

  return nil
}

func (client *api_client) create_rr (recordName string, recordType string, recordData string) (map[string]interface{}, error) {
  var data map[string]interface{}
  path := fmt.Sprintf("%s/api/v1/accounts/%s/zones/%s/rr", client.base_url, client.account_id, client.zone_name)

  type NewRr struct {
    Owner    string `json:"owner"`
    Type     string `json:"type"`
    Rdata    string `json:"rdata"`
    Comments string `json:"comments"`
  }

  recData := NewRr{
    Owner:    recordName,
    Type:     recordType,
    Rdata:    recordData,
    Comments: "created by terraform-provider-verisignmdns",
  }

  sendData, err2 := json.Marshal(recData)
  if err2 != nil {
    return nil, err2
  }

  resp, err := client.send_request("POST", path, string(sendData))
  if err != nil { return make(map[string]interface{}), err }
  if client.debug {
    log.Printf("api_client.go: create_rr got send_request response: +%v", resp)
  }

  if resp.resp_code != 201 {
    return make(map[string]interface{}), errors.New(fmt.Sprintf("Error Creating RR - HTTP %d - %s", resp.resp_code, resp.body))
  }

  if resp.location == "" {
    return make(map[string]interface{}), errors.New(fmt.Sprintf("Error Creating RR - HTTP %d - %s", resp.resp_code, resp.body))
  }

  // Ok, created, follow the Location
  resp2, err4 := client.send_request("GET", resp.location, "")
  if client.debug {
    log.Printf("api_client.go: create_rr GET got send_request response: +%v", resp2)
  }

  if err4 != nil { return make(map[string]interface{}), err4 }

  if resp2.resp_code != 200 {
    return make(map[string]interface{}), errors.New(fmt.Sprintf("Error GETting RR - HTTP %d - %s", resp2.resp_code, resp2.body))
  }

  err3 := json.Unmarshal([]byte(resp2.body), &data)
  if err3 != nil {
    return nil, err3
  }
  return data, nil
}

func (client *api_client) send_request (method string, full_uri string, data string) (api_response, error) {
  var req *http.Request
  var err error
  var locStr string

  log.Printf("client.debug=%s", client.debug)
  if client.debug {
    log.Printf("api_client.go: method='%s', full_uri='%s', data='%s'\n", method, full_uri, data)
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
  req.Header.Set("Authorization", fmt.Sprintf("Token %s", client.token))

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

  locStr = ""
  loc, err4 := resp.Location()
  if err4 != nil {
    log.Printf("api_client.go - Error parsing resp.Location()")
  } else {
    locStr = loc.String()
    if client.debug { log.Printf("api_client.go response Location header: %s", locStr)}
  }
  return api_response{
    body: body,
    resp_code: resp.StatusCode,
    location: locStr,
  }, nil
}
