package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "0.2.0"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-chapic %v\n", version)
		return
	}

	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range plugin.Files {
			if !f.Generate {
				continue
			}
			generateFile(plugin, f)
		}
		return nil
	})
}
