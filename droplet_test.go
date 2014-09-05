package digitalocean

import (
	"testing"

	. "github.com/motain/gocheck"
)

func TestDroplet(t *testing.T) {
	TestingT(t)
}

func (s *S) Test_CreateDroplet(c *C) {
	testServer.Response(202, nil, dropletExample)

	opts := CreateDroplet{
		Name:     "foobar",
		UserData: "baz",
	}

	id, err := s.client.CreateDroplet(&opts)

	req := testServer.WaitRequest()

	c.Assert(req.Form["name"], DeepEquals, []string{"foobar"})
	c.Assert(req.Form["user_data"], DeepEquals, []string{"baz"})
	c.Assert(err, IsNil)
	c.Assert(id, Equals, "25")
}

func (s *S) Test_RetrieveDroplet(c *C) {
	testServer.Response(200, nil, dropletExample)

	droplet, err := s.client.RetrieveDroplet("25")

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(droplet.StringId(), Equals, "25")
	c.Assert(droplet.RegionSlug(), Equals, "nyc1")
	c.Assert(droplet.IsLocked(), Equals, "false")
	c.Assert(droplet.NetworkingType(), Equals, "private")
	c.Assert(droplet.IPV4Address("public"), Equals, "127.0.0.20")
	c.Assert(droplet.IPV4Address("private"), Equals, "10.0.0.1")
	c.Assert(droplet.IPV6Address("public"), Equals, "")
	c.Assert(droplet.ImageSlug(), Equals, "foobar")
}

func (s *S) Test_RetrieveDroplet_noImage(c *C) {
	testServer.Response(200, nil, dropletExampleNoImage)

	droplet, err := s.client.RetrieveDroplet("25")

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(droplet.StringId(), Equals, "25")
	c.Assert(droplet.RegionSlug(), Equals, "nyc1")
	c.Assert(droplet.IsLocked(), Equals, "false")
	c.Assert(droplet.NetworkingType(), Equals, "public")
	c.Assert(droplet.IPV6Address("public"), Equals, "")
	c.Assert(droplet.IPV4Address("public"), Equals, "127.0.0.20")
	c.Assert(droplet.IPV4Address("private"), Equals, "")
	c.Assert(droplet.IPV6Address("private"), Equals, "")
	c.Assert(droplet.ImageSlug(), Equals, "")
	c.Assert(droplet.SizeSlug(), Equals, "512mb")
	c.Assert(droplet.ImageId(), Equals, "449676389")
}

func (s *S) Test_DestroyDroplet(c *C) {
	testServer.Response(204, nil, "")

	err := s.client.DestroyDroplet("25")

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
}

func (s *S) Test_Resize(c *C) {
	testServer.Response(200, nil, dropletExampleAction)

	err := s.client.Resize("25", "1gb")

	req := testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(req.Form["size"], DeepEquals, []string{"1gb"})
	c.Assert(req.Form["type"], DeepEquals, []string{"resize"})
}

func (s *S) Test_Rename(c *C) {
	testServer.Response(200, nil, dropletExampleAction)

	err := s.client.Rename("25", "foobar")

	req := testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(req.Form["name"], DeepEquals, []string{"foobar"})
	c.Assert(req.Form["type"], DeepEquals, []string{"rename"})
}

func (s *S) Test_EnableIPV6s(c *C) {
	testServer.Response(200, nil, dropletExampleAction)

	err := s.client.EnableIPV6s("25")

	req := testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(req.Form["type"], DeepEquals, []string{"enable_ipv6"})
}

func (s *S) Test_EnablePrivateNetworking(c *C) {
	testServer.Response(200, nil, dropletExampleAction)

	err := s.client.EnablePrivateNetworking("25")

	req := testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(req.Form["type"], DeepEquals, []string{"enable_private_networking"})
}

func (s *S) Test_ActionError(c *C) {
	testServer.Response(422, nil, dropletExampleActionError)

	err := s.client.EnablePrivateNetworking("25")

	req := testServer.WaitRequest()

	c.Assert(req.Form["type"], DeepEquals, []string{"enable_private_networking"})
	c.Assert(err.Error(), Equals, "Error processing droplet action: API Error: unprocessable_entity: You specified an invalid size for Droplet creation.")
}

func (s *S) Test_PowerOn(c *C) {
	testServer.Response(200, nil, dropletExampleAction)

	err := s.client.PowerOn("25")

	req := testServer.WaitRequest()

	c.Assert(req.Form["type"], DeepEquals, []string{"power_on"})
	c.Assert(err, IsNil)
}

func (s *S) Test_PowerOff(c *C) {
	testServer.Response(200, nil, dropletExampleAction)

	err := s.client.PowerOff("25")

	req := testServer.WaitRequest()

	c.Assert(req.Form["type"], DeepEquals, []string{"power_off"})
	c.Assert(err, IsNil)
}

var dropletExampleActionError = `{
  "id": "unprocessable_entity",
  "message": "You specified an invalid size for Droplet creation."
}`

var dropletExample = `{
  "droplet": {
    "id": 25,
    "name": "My-Droplet",
    "region": {
      "slug": "nyc1",
      "name": "New York",
      "sizes": [
        "1024mb",
        "512mb"
      ],
      "available": true,
      "features": [
        "virtio",
        "private_networking",
        "backups",
        "ipv6"
      ]
    },
    "image": {
      "id": 449676389,
      "name": "Ubuntu 13.04",
      "distribution": "ubuntu",
      "slug": "foobar",
      "public": true,
      "regions": [
        "nyc1"
      ],
      "created_at": "2014-07-18T16:20:40Z"
    },
    "size": {
      "slug": "512mb",
      "memory": 512,
      "vcpus": 1,
      "disk": 20,
      "transfer": null,
      "price_monthly": "5.0",
      "price_hourly": "0.00744",
      "regions": [
        "nyc1",
        "sfo1",
        "ams1"
      ]
    },
    "locked": false,
    "status": "new",
    "networks": {
      "v4": [
        {
          "ip_address": "10.0.0.1",
          "netmask": "255.255.255.0",
          "gateway": "10.0.0.0",
          "type": "private"
        },
        {
          "ip_address": "127.0.0.20",
          "netmask": "255.255.255.0",
          "gateway": "127.0.0.21",
          "type": "public"
        }
      ],
      "v6": []
    },
    "kernel": {
      "id": 485432972,
      "name": "Ubuntu 14.04 x64 vmlinuz-3.13.0-24-generic (1221)",
      "version": "3.13.0-24-generic"
    },
    "created_at": "2014-07-18T16:20:40Z",
    "features": [
      "virtio"
    ],
    "backup_ids": [

    ],
    "snapshot_ids": [

    ],
    "action_ids": [
      20
    ]
  },
  "links": {
    "actions": [
      {
        "id": 20,
        "rel": "create",
        "href": "http://example.org/v2/actions/20"
      }
    ]
  }
}`

var dropletExampleNoImage = `{
  "droplet": {
    "id": 25,
    "name": "My-Droplet",
    "region": {
      "slug": "nyc1",
      "name": "New York",
      "sizes": [
        "1024mb",
        "512mb"
      ],
      "available": true,
      "features": [
        "virtio",
        "private_networking",
        "backups",
        "ipv6"
      ]
    },
    "image": {
      "id": 449676389,
      "name": "Ubuntu 13.04",
      "distribution": "ubuntu",
      "slug": null,
      "public": true,
      "regions": [
        "nyc1"
      ],
      "created_at": "2014-07-18T16:20:40Z"
    },
    "size": {
      "slug": "512mb",
      "memory": 512,
      "vcpus": 1,
      "disk": 20,
      "transfer": null,
      "price_monthly": "5.0",
      "price_hourly": "0.00744",
      "regions": [
        "nyc1",
        "sfo1",
        "ams1"
      ]
    },
    "locked": false,
    "status": "new",
    "networks": {
      "v4": [
        {
          "ip_address": "127.0.0.20",
          "netmask": "255.255.255.0",
          "gateway": "127.0.0.21",
          "type": "public"
        }
      ],
      "v6": []
    },
    "kernel": {
      "id": 485432972,
      "name": "Ubuntu 14.04 x64 vmlinuz-3.13.0-24-generic (1221)",
      "version": "3.13.0-24-generic"
    },
    "created_at": "2014-07-18T16:20:40Z",
    "features": [
      "virtio"
    ],
    "backup_ids": [

    ],
    "snapshot_ids": [

    ],
    "action_ids": [
      20
    ]
  },
  "links": {
    "actions": [
      {
        "id": 20,
        "rel": "create",
        "href": "http://example.org/v2/actions/20"
      }
    ]
  }
}`

var dropletExampleAction = `{
  "action": {
    "id": 15,
    "status": "in-progress",
    "type": "enable_ipv6",
    "started_at": "2014-07-18T16:20:37Z",
    "completed_at": null,
    "resource_id": 15,
    "resource_type": "droplet",
    "region": "nyc1"
  }
}`
