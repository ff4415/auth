package cmd

import "context"

func Serve(ctx context.Context) {
	serve(ctx)
}

func SetConfig(cfgPath, dir string) {
	configFile = cfgPath
	watchDir = dir
}
