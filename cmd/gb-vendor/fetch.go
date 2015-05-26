package main

import (
	"fmt"
	"path/filepath"

	"github.com/constabulary/gb"
	"github.com/constabulary/gb/cmd"
	"github.com/constabulary/gb/cmd/gb-vendor/vendor"
)

func init() {
	registerCommand("fetch", FetchCmd)
}

var FetchCmd = &cmd.Command{
	ShortDesc: "fetch a remote dependency",
	Run: func(ctx *gb.Context, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("fetch: import path missing")
		}
		path := args[0]

		m, err := vendor.ReadManifest(manifestFile(ctx))
		if err != nil {
			return fmt.Errorf("could not load manifest: %T %v", err, err)
		}

		repo, extra, err := vendor.RepositoryFromPath(path)
		if err != nil {
			return err
		}

		wc, err := repo.Clone()
		if err != nil {
			return err
		}

		rev, err := wc.Revision()
		if err != nil {
			return err
		}

		branch, err := wc.Branch()
		if err != nil {
			return err
		}

		dep := vendor.Dependency{
			Importpath: path,
			Repository: repo.(*vendor.GitRepo).URL,
			Revision:   rev,
			Branch:     branch,
			Path:       extra,
		}

		if err := m.AddDependency(dep); err != nil {
			return err
		}

		dst := filepath.Join(ctx.Projectdir(), "vendor", "src", dep.Importpath)
		src := filepath.Join(wc.Dir(), dep.Path)

		if err := copypath(dst, src); err != nil {
			return err
		}

		if err := vendor.WriteManifest(manifestFile(ctx), m); err != nil {
			return err
		}
		return wc.Destroy()
	},
}
