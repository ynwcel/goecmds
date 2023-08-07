package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/urfave/cli/v2"
	"github.com/ynwcel/goecmds/goss/util"
)

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "run web site in current directory",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "drop-log",
			DefaultText: "false",
			Value:       false,
		},
		&cli.Int64Flag{
			Name:        "port",
			DefaultText: "8080",
			Value:       8080,
		},
		&cli.StringSliceFlag{
			Name:        "view-exts",
			DefaultText: util.GetViewExtsString(),
			Value:       cli.NewStringSlice(util.GetViewExtsSlice()...),
		},
	},
	Action: func(ctx *cli.Context) error {
		myLogger := log.Default()
		if ctx.Bool("drop-log") {
			myLogger.SetOutput(io.Discard)
		}

		port := ctx.Int("port")
		fmt.Println("Listen:", port)
		http.HandleFunc("/", buildHandler(ctx, myLogger))
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	},
}

func buildHandler(cliCtx *cli.Context, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			gCtx             = gctx.WithCtx(cliCtx.Context)
			gviewObj         = gview.New(".")
			status           = 200
			result           = ""
			err        error = nil
			req_path         = path.Clean(fmt.Sprintf("./%s", r.URL.Path))
			absCurPath       = ""
			absUrlPath       = ""
		)
		if gfile.IsDir(req_path) {
			req_path = filepath.Clean(fmt.Sprintf("%s/index.html", req_path))
		}

		if util.CheckExt(cliCtx.StringSlice("view-exts"), filepath.Ext(req_path)) {
			absCurPath, _ = filepath.Abs(".")
			absUrlPath, _ = filepath.Abs(req_path)

			if !strings.Contains(absUrlPath, absCurPath) {
				w.Header().Add("goss-path-err", err.Error())
				status = 500
				result = "url path error"
			} else {
				if !gfile.Exists(absUrlPath) {
					status = 404
					view_404 := "404.html"
					if gfile.Exists(view_404) {
						result, err = gviewObj.Parse(gCtx, view_404)
					} else {
						result = "404 Not Found"
					}
				} else {
					result, err = gviewObj.Parse(gCtx, absUrlPath)
				}
				if err != nil {
					w.Header().Add("goss-gview-error", err.Error())
					status = 500
					result = "Server Error!"
				}
			}
			w.WriteHeader(status)
			w.Write([]byte(result))
			log.Print(strings.Repeat("-", 91))
			log.Printf("| %3d | %-30s| %-50s|", status, r.RemoteAddr, req_path)
		} else {
			http.FileServer(http.Dir(".")).ServeHTTP(w, r)
			return
		}
	}
}
