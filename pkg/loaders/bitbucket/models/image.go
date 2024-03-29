package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"gopkg.in/yaml.v3"
)

type Image struct {
	ImageData *ImageData
}

type ImageData struct {
	Name      *string `yaml:"name"`
	RunAsUser *int64  `yaml:"run-as-user,omitempty"`
	Email     *string `yaml:"email,omitempty"`    // Email to use to fetch the Docker image
	Password  *string `yaml:"password,omitempty"` // Password to use to fetch the Docker image
	Username  *string `yaml:"username,omitempty"` // Username to use to fetch the Docker image
	Aws       *Aws    `yaml:"aws,omitempty"`      // AWS credentials
}

type Aws struct {
	AccessKey *string `yaml:"access-key"` // AWS Access Key
	SecretKey *string `yaml:"secret-key"` // AWS Secret Key
}

func (i *Image) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.StringTag {
		*i = Image{&ImageData{Name: &node.Value}}
		return nil
	}
	var image ImageData
	if err := node.Decode(&image); err != nil {
		return err
	}
	*i = Image{ImageData: &image}
	return nil
}
