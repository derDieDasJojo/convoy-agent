package volume

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/rancher/convoy-agent/cattle"
)

var rootUuidFileName = "UUID"

var Commands = []cli.Command{
	{
		Name:   "volume",
		Usage:  "Start convoy-agent as a volume agent",
		Action: volumeAgent,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "socket, s",
				Value: "/var/run/convoy/convoy.sock",
				Usage: "specify unix domain socket for communicating with convoy server",
			},
			cli.StringFlag{
				Name:   "host-uuid",
				Usage:  "set the host uuid for the host",
				EnvVar: "CATTLE_HOST_UUID",
			},
		},
		ShortName: "v",
	},
}

func volumeAgent(c *cli.Context) {
	socket := c.String("socket")
	cattleUrl := c.GlobalString("url")
	cattleAccessKey := c.GlobalString("access-key")
	cattleSecretKey := c.GlobalString("secret-key")
	if c.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	healthCheckInterval := c.GlobalInt("healthcheck-interval")
	healthCheckBaseDir := c.GlobalString("healthcheck-basedir")

	controlChan := make(chan bool, 1)

	driver := c.GlobalString("storagepool-driver")
	if driver == "" {
		logrus.Fatal("required field storagepool-driver has not been set")
	}

	hostUuid := c.String("host-uuid")
	if hostUuid == "" {
		logrus.Fatal("required field host-uuid has not been set")
	}

	cattleClient, err := cattle.NewCattleClient(cattleUrl, cattleAccessKey, cattleSecretKey)
	if err != nil {
		logrus.Fatal(err)
	}

	volAgent := NewVolumeAgent(healthCheckBaseDir, socket, hostUuid, healthCheckInterval, cattleClient, driver)

	if err := volAgent.Run(controlChan); err != nil {
		logrus.Fatal(err)
	}
}
