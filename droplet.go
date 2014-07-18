package digitalocean

import (
	"fmt"
	"strings"
)

// Droplet is used to represent a retrieved Droplet. All properties
// are set as strings.
type Droplet struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Region      string `json:"region>slug"`
	Image       string `json:"image>slug"`
	Size        string `json:"size>slug"`
	Locked      string `json:"locked"`
	Status      string `json:"status"`
	IPV6Address string `json:"networks>v6>ip_address"`
	IPV4Address string `json:"networks>v4>ip_address"`
	IPV6Type    string `json:"networks>v6>type"`
	IPV4Type    string `json:"networks>v4>type"`
}

// CreateDroplet contains the request parameters to create a new
// droplet.
type CreateDroplet struct {
	Name              string   // Name of the droplet
	Region            string   // Slug of the region to create the droplet in
	Size              string   // Slug of the size to use for the droplet
	Image             string   // Slug of the image, if using a public image
	ImageID           string   // ID of the image, if using a private image
	SSH_keys          []string // Array of SSH Key IDs that should be added
	Backups           string   // 'true' or 'false' if backups are enabled
	IPV6              string   // 'true' or 'false' if IPV6 is enabled
	PrivateNetworking string   // 'true' or 'false' if Private Networking is enabled
}

// Create creates a droplet from the parameters specified and
// returns an error if it fails. If no error is returned,
// the Droplet was succesfully created.
func (c *Client) Create(opts *CreateDroplet) error {
	// Make the request parameters
	params := make(map[string]string)

	params["name"] = opts.Name
	params["region"] = opts.Region
	params["size"] = opts.Size

	if opts.Image != "" {
		params["image"] = opts.Image
	}
	// If we specify the image_id, we override the
	// image slug
	if opts.ImageID != "" {
		params["image"] = opts.ImageID
	}

	if len(opts.SSH_keys) > 0 {
		params["ssh_keys"] = strings.Join(opts.SSH_keys, ",")
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
		return err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return fmt.Errorf("Error creating droplet: %s", parseErr(resp))
	}

	// The request was successful
	return nil
}

// Destroy destroys a droplet by the ID specified and
// returns an error if it fails. If no error is returned,
// the Droplet was succesfully destroyed.
func (c *Client) Destroy(id string) error {
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

// Retrieve gets  a droplet by the ID specified and
// returns a Droplet and an error. An error will be returned for failed
// requests with a nil Droplet.
func (c *Client) Retrieve(id string) (*Droplet, error) {
	req, err := c.NewRequest(map[string]string{}, "GET", fmt.Sprintf("/droplets/%s", id))

	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, fmt.Errorf("Error destroying droplet: %s", parseErr(resp))
	}

	droplet := &Droplet{}
	err = decodeBody(resp, droplet)

	if err != nil {
		return nil, fmt.Errorf("Error decoding droplet response: %s", err)
	}

	// The request was successful
	return droplet, nil
}
