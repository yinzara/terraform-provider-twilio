module github.com/Preskton/terraform-provider-twilio

go 1.12

require (
	github.com/fatih/structs v1.1.0
	github.com/hashicorp/terraform v0.12.10
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/kevinburke/twilio-go v0.0.0-20190630185733-fe05957cdaf8
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/onsi/ginkgo v1.10.2
	github.com/onsi/gomega v1.7.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cast v1.3.0
	github.com/stretchr/testify v1.4.0 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/kevinburke/twilio-go => github.com/yinzara/twilio-go v0.0.0-20200819210423-0111d7a4b98c
