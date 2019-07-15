package controllers

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/utils/pagination"
	"pluto-blog/models"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

//文章列表
// ArticleController operations for Article
type ListArticleController struct {
	BaseController
}

func (this *ListArticleController) Get() {
	page, err1 := this.GetInt("p")
	title := this.GetString("title")
	keywords := this.GetString("keywords")
	status := this.GetString("status")
	offset, err2 := beego.AppConfig.Int("pageoffset")
	if err1 != nil {
		page = 1
	}
	if err2 != nil {
		offset = 10
	}
	condArr := make(map[string]string)
	condArr["title"] = title
	condArr["keywords"] = keywords
	condArr["status"] = status
	countArticle := models.CountArticle(condArr)
	paginator := pagination.SetPaginator(this.Ctx, offset, countArticle)
	_, _, art := models.ListArticle(condArr, page, offset)
	this.Data["paginator"] = paginator
	this.Data["art"] = art
	//userLogin := this.GetSession("userLogin")
	//this.Data["isLogin"] = this.isLogin
	this.Data["isLogin"] = true
	this.TplName = "article.tpl"
	//this.Data["json"] = map[string]interface{}{"code": 0, "message": "ok", "data": art}
	//this.ServeJSON()

}

//显示文章
type ShowArticleController struct {
	BaseController
}

func (this *ShowArticleController) Get() {
	idstr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idstr)
	art, err := models.GetArticleById(id)
	if err != nil {
		this.Redirect("/404.html", 302)
	}
	this.Data["art"] = art

	//评论分页
	page, err1 := this.GetInt("p")
	if err1 != nil {
		page = 1
	}
	offset, err2 := beego.AppConfig.Int("pageoffset")
	if err2 != nil {
		offset = 9
	}
	where := make(map[string]string)
	where["article_id"] = idstr
	comment_count := models.CommentCount(where)
	paginator := pagination.SetPaginator(this.Ctx, offset, comment_count)
	_, _, ret := models.ListComment(where, page, offset)
	this.Data["paginator"] = paginator
	this.Data["coms"] = ret

	this.TplName = "article-detail.tpl"
	//
	//this.Data["json"] = art
	//this.ServeJSON()
}

//修改文章
type UpdateArticleController struct {
	BaseController
}

func (this *UpdateArticleController) Get() {
	idstr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idstr)
	ret_one, err := models.GetArticleById(id)
	if err != nil {
		this.Redirect("/404.html", 302)
	}
	this.Data["art"] = ret_one
	this.TplName = "article-form.tpl"
}

func (this *UpdateArticleController) Post() {
	id, _ := this.GetInt("id")
	title := this.GetString("title")
	author := this.GetString("author")
	keywords := this.GetString("keywords")
	uri := this.GetString("uri")
	summary := this.GetString("summary")
	content := this.GetString("content")
	status, _ := this.GetInt("status")
	var data models.Article
	data.Title = title
	data.Status = int8(status)
	data.Content = content
	data.Summary = summary
	data.Uri = uri
	data.Keywords = keywords
	data.Author = author
	err := models.UpdateArticle(id, data)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 200, "message": "博客修改成功", "id": id}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 500, "message": "error", "id": id}
	}
	this.ServeJSON()
}

//新增文章

type AddArticleController struct {
	BaseController
}

func (this *AddArticleController) Get() {
	var art models.Article
	art.Status = 1
	this.Data["art"] = art
	this.TplName = "article-form.tpl"
}

func (this *AddArticleController) Post() {
	title := this.GetString("title")
	author := this.GetString("author")
	keywords := this.GetString("keywords")
	uri := this.GetString("uri")
	summary := this.GetString("summary")
	content := this.GetString("content")
	status, _ := this.GetInt("status")
	var data models.Article
	data.Title = title
	data.Status = int8(status)
	data.Content = content
	data.Summary = summary
	data.Uri = uri
	data.Keywords = keywords
	data.Author = author
	id, err := models.AddArticle(data)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "ok", "id": id}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "error", "id": id}
	}
	this.ServeJSON()
}

//===================================分割线=========================
type ArticleController struct {
	beego.Controller
}

// URLMapping ...
func (c *ArticleController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// GetOne ...
// @Title Get One
// @Description get Article by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Article
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ArticleController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetArticleById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Article
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Article
// @Failure 403
// @router / [get]
func (c *ArticleController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllArticle(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Article
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Article	true		"body for Article content"
// @Success 200 {object} models.Article
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ArticleController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Article{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateArticleById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Article
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ArticleController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteArticle(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
