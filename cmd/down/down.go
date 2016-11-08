package down

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/distribution/digest"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/joshwget/lay/utils"
)

func Action(c *cli.Context) error {
	image := c.Args()[0]
	registryURL := c.GlobalString("registry")
	dir := c.String("dir")

	hub, err := registry.New(registryURL, "", "")
	if err != nil {
		return err
	}
	hub.Logf = registry.Quiet

	return down(hub, dir, image)
}

func down(hub *registry.Registry, dir string, images ...string) error {
	for _, image := range images {
		manifest, err := hub.Manifest(image, "latest")
		if err != nil {
			return err
		}

		layers := []string{}
		for _, layer := range manifest.FSLayers {
			split := strings.Split(fmt.Sprint(layer.BlobSum), ":")[1]
			layers = append(layers, split)
		}

		for _, layer := range layers {
			digest := digest.NewDigestFromHex(
				"sha256",
				layer,
			)

			reader, err := hub.DownloadLayer(image, digest)
			if err != nil {
				return err
			}
			pkg, err := utils.FindPackage(reader)
			if err != nil {
				return err
			}
			if pkg == nil {
				continue
			}
			reader.Close()

			for _, dependency := range pkg.Dependencies {
				if err = down(hub, dir, dependency); err != nil {
					return err
				}
			}

			reader, err = hub.DownloadLayer(image, digest)
			if err != nil {
				return err
			}

			log.Infof("Installing package %s", image)
			if err = utils.ExtractTar(reader, dir); err != nil {
				return err
			}
			reader.Close()
		}
	}

	return nil
}
