package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"reflect"
	"strings"
	"time"
)

type Comment struct {
	Id        int    `orm:"column(id);auto"`
	ArticleId int    `orm:"column(article_id);null"`
	Nickname  string `orm:"column(nickname);size(15);null"`
	Uri       string `orm:"column(uri);size(255);null"`
	Content   string `orm:"column(content);null"`
	Created   int64  `orm:"column(created);null"`
	Status    int8   `orm:"column(status);null" description:"0屏蔽，1正常"`
}

func (t *Comment) TableName() string {
	return "comment"
}

func init() {
	orm.RegisterModel(new(Comment))
}

func ListComment(where map[string]string, page int, offset int) (num int64, err error, result []Comment) {
	model := orm.NewOrm()
	query := model.QueryTable("comment")
	data := orm.NewCondition()
	if where["article_id"] != "" {
		data = data.And("article_id", where["article_id"])
	}
	if where["status"] != "" {
		data = data.And("status", where["status"])
	}
	query = query.SetCond(data)
	start := (page - 1) * offset
	var comments []Comment
	num, err1 := query.Limit(offset, start).All(&comments)
	return num, err1, comments

}

func CommentCount(where map[string]string) int64 {
	model := orm.NewOrm()
	query := model.QueryTable("comment")
	data := orm.NewCondition()
	if where["article_id"] != "" {
		data = data.And("article_id", where["article_id"])
	}
	if where["status"] != "" {
		data = data.And("status", where["status"])
	}
	num, _ := query.SetCond(data).Count()
	return num
}

//新增评论
func AddComment(data Comment) (id int64, err error) {
	model := orm.NewOrm()
	model.Using("detault")
	com := new(Comment)

	com.ArticleId = data.ArticleId
	com.Nickname = data.Nickname
	com.Uri = data.Uri
	com.Content = data.Content
	com.Created = time.Now().Unix()
	com.Status = data.Status

	id, err1 := model.Insert(com)
	return id, err1

}

//==================================分割线=======================
// AddComment insert a new Comment into database and returns
// last inserted Id on success.

// GetCommentById retrieves Comment by Id. Returns error if
// Id doesn't exist
func GetCommentById(id int) (v *Comment, err error) {
	o := orm.NewOrm()
	v = &Comment{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllComment retrieves all Comment matches certain condition. Returns empty list if
// no records exist
func GetAllComment(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Comment))
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

	var l []Comment
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

// UpdateComment updates Comment by Id and returns error if
// the record to be updated doesn't exist
func UpdateCommentById(m *Comment) (err error) {
	o := orm.NewOrm()
	v := Comment{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteComment deletes Comment by Id and returns error if
// the record to be deleted doesn't exist
func DeleteComment(id int) (err error) {
	o := orm.NewOrm()
	v := Comment{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Comment{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
