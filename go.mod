module github.com/TNK-Studio/lazykube

go 1.14

replace (
	github.com/jroimartin/gocui v0.4.0 => github.com/elfgzp/gocui v0.4.1-0.20201118030412-21fac610f2e0
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
)

require (
	github.com/atotto/clipboard v0.1.2
	github.com/docker/distribution v2.7.1+incompatible
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/fatih/camelcase v1.0.0
	github.com/go-logr/logr v0.2.0
	github.com/gookit/color v1.3.2
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jesseduffield/asciigraph v0.0.0-20190605104717-6d88e39309ee
	github.com/jroimartin/gocui v0.4.0
	github.com/kr/pretty v0.2.1 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/nsf/termbox-go v0.0.0-20200418040025-38ba6e5628f1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spkg/bom v1.0.0
	github.com/tebeka/strftime v0.1.5 // indirect
	golang.org/x/net v0.0.0-20200927032502-5d4f70055728 // indirect
	golang.org/x/sys v0.0.0-20200926100807-9d91bd62050c // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/cli-runtime v0.19.3
	k8s.io/client-go v0.19.3
	k8s.io/klog/v2 v2.2.0
	k8s.io/kubectl v0.19.3
	k8s.io/metrics v0.19.3
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
)
