package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Album struct {
	Id       int    `orm:"column(id);auto"`
	Title    string `orm:"column(title);size(255)" description:"文章标题"`
	Picture  string `orm:"column(picture);size(255);null" description:"Picture"`
	Keywords string `orm:"column(keywords);size(2550);null" description:"关键词"`
	Summary  string `orm:"column(summary);size(255);null"`
	Created  int    `orm:"column(created);null" description:"发布时间"`
	Viewnum  int    `orm:"column(viewnum);null" description:"阅读次数"`
	Status   int8   `orm:"column(status);null" description:"状态: 0草稿，1已发布"`
}

func (t *Album) TableName() string {
	return "album"
}

func init() {
	orm.RegisterModel(new(Album))
}

// AddAlbum insert a new Album into database and returns
// last inserted Id on success.
func AddAlbum(m *Album) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAlbumById retrieves Album by Id. Returns error if
// Id doesn't exist
func GetAlbumById(id int) (v *Album, err error) {
	o := orm.NewOrm()
	v = &Album{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAlbum retrieves all Album matches certain condition. Returns empty list if
// no records exist
func GetAllAlbum(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Album))
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

	var l []Album
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

// UpdateAlbum updates Album by Id and returns error if
// the record to be updated doesn't exist
func UpdateAlbumById(m *Album) (err error) {
	o := orm.NewOrm()
	v := Album{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAlbum deletes Album by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAlbum(id int) (err error) {
	o := orm.NewOrm()
	v := Album{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Album{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
