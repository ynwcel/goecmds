package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/urfave/cli/v2"
	"github.com/ynwcel/goecmds/goss/util"
)

var buildCmd = &cli.Command{
	Name:  "build",
	Usage: "build html file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			DefaultText: getBuildDir(),
			Value:       getBuildDir(),
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Value:   false,
		},
		&cli.StringSliceFlag{
			Name:        "view-exts",
			DefaultText: util.GetViewExtsString(),
			Value:       cli.NewStringSlice(util.GetViewExtsSlice()...),
		},
	},
	Action: func(ctx *cli.Context) error {
		//创建生成目录
		build_dir := ctx.String("output")
		if ctx.Bool("force") {
			if err := gfile.Remove(build_dir); err != nil {
				return err
			}
		}
		if gfile.IsDir(build_dir) {
			return fmt.Errorf("`%s` folder exists!", build_dir)
		}
		return buildAction(ctx, ".", build_dir)
	},
}

func getBuildDir() string {
	return fmt.Sprintf("./_build.%s", time.Now().Format("20060102.1504"))
}

func buildAction(cliCtx *cli.Context, src, target string) error {
	curfs, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if !gfile.IsDir(target) {
		if err := gfile.Mkdir(target); err != nil {
			return err
		}
	}
	view := gview.New(src)
	for _, f := range curfs {
		fname := f.Name()
		fpath := filepath.Clean(fmt.Sprintf("%s/%s", src, fname))
		bfpath := filepath.Clean(fmt.Sprintf("%s/%s", target, fname))

		if f.IsDir() {
			if err := buildAction(cliCtx, fpath, bfpath); err != nil {
				return err
			}
		} else {
			if util.CheckExt(cliCtx.StringSlice("view-exts"), filepath.Ext(fname)) {
				log.Printf("Build: %-60s ==> %s", fpath, bfpath)
				content, err := view.Parse(cliCtx.Context, fname)
				if err != nil {
					return fmt.Errorf("`%s` build error:%w", fpath, err)
				}
				if err := gfile.PutContents(bfpath, content); err != nil {
					return fmt.Errorf("`%s` build error:%s", fname, err)
				}
			} else {
				log.Printf("Copy : %-60s ==> %s", fpath, bfpath)
				if err := gfile.Copy(fpath, bfpath); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
