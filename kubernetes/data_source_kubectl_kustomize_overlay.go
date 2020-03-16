package kubernetes

import (
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
)

func dataSourceKubectlKustomizeOverlay() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubectlKustomizeOverlayRead,
		Schema: map[string]*schema.Schema{
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"documents": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceKubectlKustomizeOverlayRead(d *schema.ResourceData, m interface{}) error {
	content := d.Get("content").(string)
	rm, err := runKustomizeBuild(content)
	if err != nil {
		return err
	}

	yaml, err := rm.AsYaml()
	if err != nil {
		return err
	}

	documents, err := splitMultiDocumentYAML(string(yaml))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%x", sha256.Sum256([]byte(content))))
	d.Set("documents", documents)
	return nil
}

func runKustomizeBuild(path string) (rm resmap.ResMap, err error) {
	fSys := filesys.MakeFsOnDisk()
	opts := &krusty.Options{
		DoLegacyResourceSort: true,
		LoadRestrictions:     types.LoadRestrictionsRootOnly,
		DoPrune:              false,
	}

	k := krusty.MakeKustomizer(fSys, opts)

	rm, err = k.Run(path)
	if err != nil {
		return nil, err
	}

	return rm, nil
}
