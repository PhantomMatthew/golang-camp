module newline.com/dzzs

go 1.13

require (
	newline.com/newline v0.0.1
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gin-gonic/gin v1.6.2
	github.com/jinzhu/gorm v1.9.12
	github.com/sirupsen/logrus v1.5.0
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.5
	moul.io/zapgorm v1.0.0
)

replace newline.com/newline => ./../../..
