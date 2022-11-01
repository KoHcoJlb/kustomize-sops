package main

import (
	"fmt"
	"go.mozilla.org/sops/v3/cmd/sops/formats"
	"go.mozilla.org/sops/v3/decrypt"
	"os"
	"sigs.k8s.io/kustomize/api/resource"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

var buildAnnotations = append(resource.BuildAnnotations, "kustomize.config.k8s.io/id")

func removeBuildAnnotations(node *yaml.RNode) {
	annotations := node.GetAnnotations()
	for _, name := range buildAnnotations {
		delete(annotations, name)
	}
	if err := node.SetAnnotations(annotations); err != nil {
		panic(err)
	}
}

func mergeBuildAnnotations(src, dst *yaml.RNode) {
	srcAnnotations := src.GetAnnotations()
	dstAnnotations := dst.GetAnnotations()
	for _, name := range buildAnnotations {
		dstAnnotations[name] = srcAnnotations[name]
	}
	if err := dst.SetAnnotations(dstAnnotations); err != nil {
		panic(err)
	}
}

func filter(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
	for idx, node := range nodes {
		value := node.Field("sops")
		if value != nil {
			newNode := node.Copy()
			removeBuildAnnotations(newNode)

			cleartext, err := decrypt.DataWithFormat([]byte(newNode.MustString()), formats.Yaml)
			if err != nil {
				return nil, err
			}

			newNode, err = yaml.Parse(string(cleartext))
			if err != nil {
				return nil, fmt.Errorf("redecode resource: %w", err)
			}

			mergeBuildAnnotations(node, newNode)
			nodes[idx] = newNode
		}
	}

	return nodes, nil
}

func main() {
	pipeline := kio.Pipeline{
		Inputs: []kio.Reader{
			&kio.ByteReader{
				Reader: os.Stdin,
			},
		},
		Outputs: []kio.Writer{
			&kio.ByteWriter{
				Writer: os.Stdout,
			},
		},
		Filters: []kio.Filter{
			kio.FilterFunc(filter),
		},
	}
	err := pipeline.Execute()
	if err != nil {
		panic(err)
	}
}
