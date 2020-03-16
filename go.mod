module github.com/gavinbunney/terraform-provider-kubectl

go 1.13

require (
	github.com/cenkalti/backoff v2.1.1+incompatible
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/hcl2 v0.0.0-20190702185634-5b39d9ff3a9a
	github.com/hashicorp/terraform v0.12.5
	github.com/icza/dyno v0.0.0-20180601094105-0c96289f9585
	github.com/mitchellh/go-homedir v1.1.0
	github.com/stretchr/testify v1.4.0
	github.com/zclconf/go-cty v1.2.1
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0-20190831074750-7364b6bdad65
	k8s.io/apimachinery v0.0.0-20190831074630-461753078381
	k8s.io/cli-runtime v0.0.0-20190831080432-9d670f2021f4
	k8s.io/client-go v0.0.0-20190831074946-3fe2abece89e
	k8s.io/kube-aggregator v0.0.0-20190831115419-e81a1546b343
	k8s.io/kubectl v0.0.0-20190831163037-3b58a944563f
	sigs.k8s.io/kustomize/api v0.2.0
	sigs.k8s.io/yaml v1.1.0
)

replace github.com/Azure/go-autorest v10.15.4+incompatible => github.com/Azure/go-autorest v13.0.0+incompatible

replace github.com/russross/blackfriday v2.0.0+incompatible => github.com/russross/blackfriday v1.5.2
