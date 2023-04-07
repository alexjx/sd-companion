package cmd

import (
	"net/http"
	"os"
	"strconv"

	"github.com/alexjx/sd-companion/broswer"
	"github.com/alexjx/sd-companion/pages"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var ServeCmd = &cli.Command{
	Name:  "serve",
	Usage: "serve the web application",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "root",
			Aliases:  []string{"r"},
			Usage:    "the root directory of the broswer",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "listen",
			Aliases: []string{"l"},
			Usage:   "address and port of the web application",
			Value:   "127.0.0.1:9080",
		},
		&cli.StringSliceFlag{
			Name:    "extensions",
			Aliases: []string{"e"},
			Usage:   "the extensions of the files to broswer",
			Value:   []string{"jpg", "jpeg", "png", "gif"},
		},
		&cli.IntFlag{
			Name:    "quality",
			Aliases: []string{"q"},
			Usage:   "the quality of the jpeg image",
			Value:   80,
		},
		&cli.StringFlag{
			Name:     "trash",
			Aliases:  []string{"t"},
			Usage:    "the trash directory",
			Required: true,
		},
	},
	Action: serveAction,
}

type ServeConfig struct {
	Root    string
	Listen  string
	Ext     []string
	Quality int
	Transh  string
}

// NewEngine create a gin engine
func NewEngine(cfg *ServeConfig, b *broswer.Broswer) *gin.Engine {
	r := gin.New()

	// setup logger
	logger := logrus.WithFields(logrus.Fields{
		"component": "gin",
	})
	r.Use(gin.LoggerWithWriter(logger.WriterLevel(logrus.InfoLevel)))

	// recovery
	r.Use(gin.Recovery())

	// static
	staticFs := pages.EmbedFolder(pages.StaticFS, "image_broswer/dist", true)
	r.Use(static.Serve("/", staticFs))
	r.NoRoute(func(c *gin.Context) {
		c.FileFromFS("index.html", staticFs)
	})

	// cors
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	r.Use(cors.New(corsCfg))

	api := r.Group("/api")
	{
		// return the root path
		api.GET("/root", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"root": cfg.Root,
			})
		})

		api.GET("/files", func(c *gin.Context) {
			files, err := b.Files()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"files": files,
			})
		})

		api.DELETE("/file", func(c *gin.Context) {
			path := c.Query("path")
			if path == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "path is required",
				})
				return
			}

			err := b.Delete(path)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "file deleted",
			})
		})

		api.GET("/metadata", func(c *gin.Context) {
			path := c.Query("path")
			if path == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "path is required",
				})
				return
			}

			metadata, err := b.Metadata(path)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"metadata": metadata,
			})
		})

		api.GET("/encoded", func(c *gin.Context) {
			var (
				height int
				width  int
			)
			path := c.Query("path")
			if path == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "path is required",
				})
				return
			}
			if h := c.Query("height"); h != "" {
				hh, err := strconv.ParseInt(h, 10, 32)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}
				height = int(hh)
			}
			if w := c.Query("width"); w != "" {
				ww, err := strconv.ParseInt(w, 10, 32)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}
				width = int(ww)
			}

			encoded, err := b.Encoded(path, height, width)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.Data(http.StatusOK, "image/jpeg", encoded.Data)
		})
	}

	files := r.Group("/files")
	{
		files.StaticFS("/", http.Dir(cfg.Root))
	}

	return r
}

// NewBroswer create a broswer
func NewBroswer(cfg *ServeConfig) *broswer.Broswer {
	return broswer.NewBroswer(cfg.Root, cfg.Transh, cfg.Ext, cfg.Quality)
}

func serve(cfg *ServeConfig, engine *gin.Engine) {
	go engine.Run(cfg.Listen)
}

func serveAction(cctx *cli.Context) error {
	logrus.SetLevel(logrus.DebugLevel)

	cfg := &ServeConfig{
		Root:    cctx.String("root"),
		Listen:  cctx.String("listen"),
		Ext:     cctx.StringSlice("extensions"),
		Quality: cctx.Int("quality"),
		Transh:  cctx.String("trash"),
	}

	fxApp := fx.New(
		fx.Supply(cfg),

		fx.Provide(
			NewEngine,
			NewBroswer,
		),

		fx.Invoke(
			serve,
		),

		// misc
		fx.WithLogger(
			func() fxevent.Logger {
				return &fxevent.ConsoleLogger{
					W: os.Stderr,
				}
			},
		),
	)

	fxApp.Run()
	return nil
}
