package docs

import (
	"bytes"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @ host 127.0.0.1
// @basePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

type replacerResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (rrw replacerResponseWriter) Write(b []byte) (int, error) {
	rrw.Body.Write(b)
	return rrw.ResponseWriter.Write([]byte{})
}

func docImprove(c *gin.Context, title string) {
	lastModifiedFlag := "0"
	executablePath, _ := os.Executable()
	statInfo, err := os.Stat(executablePath)
	if err == nil {
		lastModifiedFlag = strconv.Itoa(int(statInfo.ModTime().Unix()))
	}
	if c.Request.RequestURI == "/doc/last-modified-flag" {
		c.PureJSON(http.StatusOK, gin.H{"lastModified": lastModifiedFlag})
		c.Abort()
	} else if c.Request.RequestURI == "/doc/index.html" {
		rrw := replacerResponseWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = rrw
		c.Next()
		headExt := `
			<link rel="stylesheet" type="text/css" href="/assets/doc/extend.css" >
		`
		bodyExt := `
			<script>window.lastModifiedFlag = '` + lastModifiedFlag + `';</script>
			<script src="/assets/doc/extend.js"></script>
		`
		hc := rrw.Body.String()
		var re = regexp.MustCompile(`<title>.*</title>`)
		hc = re.ReplaceAllString(hc, "<title>"+title+"</title>")
		hc = strings.ReplaceAll(hc, `<link href="https://fonts.googleapis.com/css?family=Open+Sans:400,700|Source+Code+Pro:300,600|Titillium+Web:400,600,700" rel="stylesheet">`, "")
		hc = strings.ReplaceAll(hc, "</head>", headExt+"</head>")
		hc = strings.ReplaceAll(hc, "</body>", bodyExt+"</body>")
		rrw.ResponseWriter.Write([]byte(hc))
	} else {
		c.Next()
	}
}

// HandlerDoc HandlerDoc
func HandlerDoc(r *gin.Engine, docConfig map[string]string) {

	title := docConfig["title"]
	description := docConfig["description"]
	version := docConfig["version"]
	listen := docConfig["listen"]
	username := docConfig["username"]
	password := docConfig["password"]
	auth := docConfig["auth"]

	SwaggerInfo.Title = title
	SwaggerInfo.Description = description
	SwaggerInfo.Version = version

	authHandler := gin.BasicAuth(gin.Accounts{username: password})
	if auth == "false" {
		authHandler = func(c *gin.Context) {}
	}

	r.GET("/doc/*any",
		authHandler,
		func(c *gin.Context) { docImprove(c, title) },
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
	r.GET("/doc", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/doc/index.html")
	})
	go func() {
		time.Sleep(time.Second)
		logrus.Info("doc on http://" + listen + "/doc")
	}()
}
