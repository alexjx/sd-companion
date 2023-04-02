package cmd

import (
	"net/http"
	"os"

	"github.com/alexjx/image-browser/broswer"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var serveCmd = cli.Command{
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
			Value:   ":8088",
		},
	},
	Action: serveAction,
}

type ServeConfig struct {
	Root string
	Port string
}

// NewEngine create a gin engine
func NewEngine(cfg *ServeConfig, b *broswer.Broswer) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
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

		api.DELETE("/file/:path", func(c *gin.Context) {
			path := c.Param("path")
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

		api.GET("/metadata/:path", func(c *gin.Context) {
			path := c.Param("path")
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
	}

	return r
}

func serve(cfg *ServeConfig, engine *gin.Engine) {
	go engine.Run(cfg.Port)
}

func serveAction(cctx *cli.Context) error {
	logrus.SetLevel(logrus.DebugLevel)

	cfg := &ServeConfig{
		Root: cctx.String("root"),
		Port: cctx.String("port"),
	}

	fxApp := fx.New(
		fx.Supply(cfg),

		fx.Provide(
			NewEngine,
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
