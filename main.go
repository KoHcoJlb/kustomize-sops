package main

import (
	"fmt"
	"go.mozilla.org/sops/v3/aes"
	"go.mozilla.org/sops/v3/cmd/sops/common"
	"go.mozilla.org/sops/v3/cmd/sops/formats"
	"os"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func decrypt(data []byte) (cleartext []byte, err error) {
	store := common.StoreForFormat(formats.Yaml)

	// Load SOPS file and access the data key
	tree, err := store.LoadEncryptedFile(data)
	if err != nil {
		return nil, err
	}
	key, err := tree.Metadata.GetDataKey()
	if err != nil {
		return nil, err
	}

	// Decrypt the tree
	cipher := aes.NewCipher()
	_, err = tree.Decrypt(key, cipher)
	if err != nil {
		return nil, err
	}

	return store.EmitPlainFile(tree.Branches)
}

func filter(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
	for idx, node := range nodes {
		value := node.Field("sops")
		if value != nil {
			cleartext, err := decrypt([]byte(node.MustString()))
			if err != nil {
				return nil, err
			}

			nodes[idx], err = yaml.Parse(string(cleartext))
			if err != nil {
				return nil, fmt.Errorf("redecode resource: %w", err)
			}
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
