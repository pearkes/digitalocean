package digitalocean

import (
	"github.com/pearkes/digitalocean/testutil"
	"testing"

	. "github.com/motain/gocheck"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct {
	client *Client
}

var _ = Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *C) {
	testServer.Start()
	var err error
	s.client, err = NewClient("foobar")
	s.client.URL = "http://localhost:4444"
	if err != nil {
		panic(err)
	}
}

func (s *S) TearDownTest(c *C) {
	testServer.Flush()
}

func (s *S) Test_CreateDroplet(c *C) {
	testServer.Response(202, nil, dropletExample)

	opts := CreateDroplet{
		Name: "foobar",
	}

	id, err := s.client.CreateDroplet(&opts)

	req := testServer.WaitRequest()

	c.Assert(req.Form["name"], DeepEquals, []string{"foobar"})
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
	c.Assert(droplet.NetworkingType(), Equals, "public")
	c.Assert(droplet.IPV6Address(), Equals, "")
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
	c.Assert(droplet.IPV6Address(), Equals, "")
	c.Assert(droplet.ImageSlug(), Equals, "")
	c.Assert(droplet.ImageId(), Equals, "449676389")
}

func (s *S) Test_DestroyDroplet(c *C) {
	testServer.Response(204, nil, "")

	err := s.client.DestroyDroplet("25")

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
}

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
