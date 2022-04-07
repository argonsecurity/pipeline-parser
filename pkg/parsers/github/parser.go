package github

import "gopkg.in/yaml.v3"

func Parse(data []byte) (*Root, error) {
	root := &Root{}
	if err := yaml.Unmarshal(data, root); err != nil {
		return nil, err
	}
	return root, nil
}
