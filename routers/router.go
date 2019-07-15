// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"pluto-blog/controllers"
)

func init() {
	beego.Router("/", &controllers.ListArticleController{})
	beego.Router("/404.html", &controllers.BaseController{})

	beego.Router("/article", &controllers.ListArticleController{})
	beego.Router("/article/:id", &controllers.ShowArticleController{})
	beego.Router("/article/edit/:id", &controllers.UpdateArticleController{})
	beego.Router("/article/add", &controllers.AddArticleController{})

	//评论
	beego.Router("/comment/add", &controllers.AddCommentController{})

	//登录
	beego.Router("/login", &controllers.UserLoginController{})
	beego.Router("/logout", &controllers.UserLogoutController{})
}
