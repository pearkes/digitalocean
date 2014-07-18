package digitalocean

import (
	"fmt"
	"strconv"
	"strings"
)

type DropletResponse struct {
	Droplet *Droplet `json:"droplet"`
}

// Droplet is used to represent a retrieved Droplet. All properties
// are set as strings.
type Droplet struct {
	Id       int64                               `json:"id"`
	Name     string                              `json:"name"`
	Region   map[string]interface{}              `json:"region"`
	Image    map[string]interface{}              `json:"image"`
	Size     map[string]interface{}              `json:"size"`
	Locked   bool                                `json:"locked"`
	Status   string                              `json:"status"`
	Networks map[string][]map[string]interface{} `json:"networks"`
}

// Returns the slug for the region
func (d *Droplet) RegionSlug() string {
	return d.Region["slug"].(string)
}

// Returns the slug for the region
func (d *Droplet) StringId() string {
	return strconv.FormatInt(d.Id, 10)
}

// Returns the string for Locked
func (d *Droplet) IsLocked() string {
	return strconv.FormatBool(d.Locked)
}

// Returns the slug for the image
func (d *Droplet) ImageSlug() string {
	if _, ok := d.Image["slug"]; ok {
		return d.Image["slug"].(string)
	} else {
		return strconv.FormatInt(d.Image["id"].(int64), 10)
	}
}

// Returns the slug for the size
func (d *Droplet) SizeSlug() string {
	return d.Size["slug"].(string)
}

// Returns the ipv4 address
func (d *Droplet) IPV4Address() string {
	if _, ok := d.Networks["v4"]; ok {
		return d.Networks["v4"][0]["ip_address"].(string)
	}

	return ""
}

// Returns the ipv6 adddress
func (d *Droplet) IPV6Address() string {
	if arr, ok := d.Networks["v6"]; ok && len(arr) > 0 {
		return d.Networks["v6"][0]["ip_address"].(string)
	}

	return ""
}

// Currently DO only has a network type per droplet,
// so we just takes ipv4s
func (d *Droplet) NetworkingType() string {
	if _, ok := d.Networks["v4"]; ok {
		return d.Networks["v4"][0]["type"].(string)
	}

	return ""
}

// CreateDroplet contains the request parameters to create a new
// droplet.
type CreateDroplet struct {
	Name              string   // Name of the droplet
	Region            string   // Slug of the region to create the droplet in
	Size              string   // Slug of the size to use for the droplet
	Image             string   // Slug of the image, if using a public image
	SSHKeys           []string // Array of SSH Key IDs that should be added
	Backups           string   // 'true' or 'false' if backups are enabled
	IPV6              string   // 'true' or 'false' if IPV6 is enabled
	PrivateNetworking string   // 'true' or 'false' if Private Networking is enabled
}

// CreateDroplet creates a droplet from the parameters specified and
// returns an error if it fails. If no error and an ID is returned,
// the Droplet was succesfully created.
func (c *Client) CreateDroplet(opts *CreateDroplet) (string, error) {
	// Make the request parameters
	params := make(map[string]string)

	params["name"] = opts.Name
	params["region"] = opts.Region
	params["size"] = opts.Size

	if opts.Image != "" {
		params["image"] = opts.Image
	}

	if len(opts.SSHKeys) > 0 {
		params["ssh_keys"] = strings.Join(opts.SSHKeys, ",")
	}

	if opts.Backups == "" {
		params["backups"] = "false"
	} else {
		params["backups"] = opts.Backups
	}

	if opts.IPV6 == "" {
		params["ipv6"] = "false"
	} else {
		params["ipv6"] = opts.IPV6
	}

	if opts.PrivateNetworking == "" {
		params["private_networking"] = "false"
	} else {
		params["private_networking"] = opts.PrivateNetworking
	}

	req, err := c.NewRequest(params, "POST", "/droplets")
	if err != nil {
		return "", err
	}

	resp, err := checkResp(c.Http.Do(req))

	if err != nil {
		return "", fmt.Errorf("Error creating droplet: %s", err)
	}

	droplet := new(DropletResponse)

	err = decodeBody(resp, &droplet)

	if err != nil {
		return "", fmt.Errorf("Error parsing droplet response: %s", err)
	}

	// The request was successful
	return droplet.Droplet.StringId(), nil
}

// DestroyDroplet destroys a droplet by the ID specified and
// returns an error if it fails. If no error is returned,
// the Droplet was succesfully destroyed.
func (c *Client) DestroyDroplet(id string) error {
	req, err := c.NewRequest(map[string]string{}, "DELETE", fmt.Sprintf("/droplets/%s", id))

	if err != nil {
		return err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return fmt.Errorf("Error destroying droplet: %s", parseErr(resp))
	}

	// The request was successful
	return nil
}

// RetrieveDroplet gets  a droplet by the ID specified and
// returns a Droplet and an error. An error will be returned for failed
// requests with a nil Droplet.
func (c *Client) RetrieveDroplet(id string) (*Droplet, error) {
	req, err := c.NewRequest(map[string]string{}, "GET", fmt.Sprintf("/droplets/%s", id))

	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, fmt.Errorf("Error destroying droplet: %s", parseErr(resp))
	}

	droplet := new(DropletResponse)

	err = decodeBody(resp, &droplet)

	if err != nil {
		return nil, fmt.Errorf("Error decoding droplet response: %s", err)
	}

	// The request was successful
	return droplet.Droplet, nil
}
