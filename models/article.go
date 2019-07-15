package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Article struct {
	Id       int    `orm:"column(id);auto"`
	Title    string `orm:"column(title);size(255)" description:"文章标题"`
	Uri      string `orm:"column(uri);size(255);null" description:"URL"`
	Keywords string `orm:"column(keywords);size(2550);null" description:"关键词"`
	Summary  string `orm:"column(summary);size(255);null"`
	Content  string `orm:"column(content)" description:"正文"`
	Author   string `orm:"column(author);size(20);null" description:"作者"`
	Created  int64  `orm:"column(created);null" description:"发布时间"`
	Viewnum  int    `orm:"column(viewnum);null" description:"阅读次数"`
	Status   int8   `orm:"column(status);null" description:"状态: 0草稿，1已发布"`
}

func (t *Article) TableName() string {
	return "article"
}

func init() {
	orm.RegisterModel(new(Article))
}

func ListArticle(condArr map[string]string, page int, offset int) (num int64, err error, art []Article) {
	o := orm.NewOrm()
	qs := o.QueryTable("article")
	cond := orm.NewCondition()
	if condArr["title"] != "" {
		cond = cond.And("title_icontains", condArr["title"])
	}
	if condArr["keywords"] != "" {
		cond = cond.Or("keywords__icontains", condArr["title"])
	}
	if condArr["status"] != "" {
		cond = cond.And("status", condArr["title"])
	}
	qs = qs.SetCond(cond)
	if page < 1 {
		page = 1
	}
	if offset < 1 {
		offset = 9
	}
	start := (page - 1) * offset
	var articles []Article
	num, err1 := qs.OrderBy("-created").Limit(offset, start).All(&articles)
	return num, err1, articles
}

func CountArticle(condArr map[string]string) int64 {
	o := orm.NewOrm()
	qs := o.QueryTable("article")
	cond := orm.NewCondition()
	if condArr["title"] != "" {
		cond = cond.And("title_icontains", condArr["title"])
	}
	if condArr["keywords"] != "" {
		cond = cond.Or("keywords__icontains", condArr["title"])
	}
	if condArr["status"] != "" {
		cond = cond.And("status", condArr["title"])
	}
	num, _ := qs.SetCond(cond).Count()

	return num
}

// GetArticleById retrieves Article by Id. Returns error if
// Id doesn't exist
func GetArticleById(id int) (v *Article, err error) {
	o := orm.NewOrm()
	v = &Article{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//修改数据

func UpdateArticle(id int, data Article) error {
	model := orm.NewOrm()
	model.Using("default")
	//查询数据
	ret := Article{Id: id}
	ret.Title = data.Title
	ret.Uri = data.Uri
	ret.Keywords = data.Keywords
	ret.Summary = data.Summary
	ret.Content = data.Content
	ret.Author = data.Author
	ret.Status = data.Status
	ret.Created = time.Now().Unix()
	_, err := model.Update(&ret)
	return err
}

//新增
func AddArticle(updArt Article) (int64, error) {
	model := orm.NewOrm()
	model.Using("default")
	art := new(Article)

	art.Title = updArt.Title
	art.Uri = updArt.Uri
	art.Keywords = updArt.Keywords
	art.Summary = updArt.Summary
	art.Content = updArt.Content
	art.Author = updArt.Author
	art.Created = time.Now().Unix()
	art.Viewnum = 1
	art.Status = updArt.Status

	id, err := model.Insert(art)
	return id, err
}

//=====================================================分割线==========================
// AddArticle insert a new Article into database and returns
// last inserted Id on success.

// GetAllArticle retrieves all Article matches certain condition. Returns empty list if
// no records exist
func GetAllArticle(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Article))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Article
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateArticle updates Article by Id and returns error if
// the record to be updated doesn't exist
func UpdateArticleById(m *Article) (err error) {
	o := orm.NewOrm()
	v := Article{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteArticle deletes Article by Id and returns error if
// the record to be deleted doesn't exist
func DeleteArticle(id int) (err error) {
	o := orm.NewOrm()
	v := Article{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Article{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
