package common

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Orders struct {
	Id            int64  `orm:"column(id);auto"`
	AppId         string `orm:"column(appid);size(128);null"`
	ReturnCode    string `orm:"column(return_code);"`
	ReturnMsg     string `orm:"column(return_msg)"`
	MchId         string `orm:"column(mch_id)"`
	Nonce         string `orm:"column(nonce_str)"`
	Sign          string `orm:"column(sign)"`
	ResultCode    string `orm:"column(result_code)"`
	ErrCode       string `orm:"column(err_code)"`
	ErrCodeDes    string `orm:"column(err_code_des)"`
	OpenId        string `orm:"column(openid)"`
	IsSubscribe   string `orm:"column(is_subscribe)"`
	TradeType     string `orm:"column(trade_type)"`
	BankType      string `orm:"column(bank_type)"`
	TotalFee      uint32 `orm:"column(total_fee)"`
	CashFee       uint32 `orm:"column(cash_fee)"`
	TransactionId string `orm:"column(transaction_id)"`
	OutTradeNo    string `orm:"column(out_trade_no)"`
	TimeEnd       string `orm:"column(time_end)"`
	UnixTime      int64  `orm:"column(unix_time)"`
}

func (t *Orders) TableName() string {
	return "payplat_orders"
}

func init() {
	orm.RegisterModel(new(Orders))
}

func AddOrder(m *Orders) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetOrder(m *Orders, col ...string) error {
	o := orm.NewOrm()
	return o.Read(m, col...)
}

func GetOrderById(id int64) (v *Orders, err error) {
	o := orm.NewOrm()
	v = &Orders{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func UpdateOrderById(m *Orders, cols ...string) (err error) {
	o := orm.NewOrm()
	v := Orders{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, cols...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteOrder(id int64) (err error) {
	o := orm.NewOrm()
	v := Orders{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Orders{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetOrders(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Orders))
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

	var l []Orders
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
