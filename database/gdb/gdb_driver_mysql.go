// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/dxasu/gf.

package gdb

import (
	"database/sql"
	"fmt"

	"github.com/dxasu/gf/errors/gerror"
	"github.com/dxasu/gf/internal/intlog"
	"github.com/dxasu/gf/text/gregex"
	"github.com/dxasu/gf/text/gstr"

	_ "github.com/go-sql-driver/mysql"
)

// DriverMysql is the driver for mysql database.
type DriverMysql struct {
	*Core
}

// New creates and returns a database object for mysql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *DriverMysql) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverMysql{
		Core: core,
	}, nil
}

// Open creates and returns a underlying sql.DB object for mysql.
// Note that it converts time.Time argument to local timezone in default.
func (d *DriverMysql) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gregex.ReplaceString(`/([\w\.\-]+)+`, "/"+config.Name, source)
		}
	} else {
		source = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s",
			config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset,
		)
	}
	intlog.Printf("Open: %s", source)
	if db, err := sql.Open("mysql", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// FilteredLinkInfo retrieves and returns filtered `linkInfo` that can be using for
// logging or tracing purpose.
func (d *DriverMysql) FilteredLinkInfo() string {
	linkInfo := d.GetConfig().LinkInfo
	if linkInfo == "" {
		return ""
	}
	s, _ := gregex.ReplaceString(
		`(.+?):(.+)@tcp(.+)`,
		`$1:xxx@tcp$3`,
		linkInfo,
	)
	return s
}

// GetChars returns the security char for this type of database.
func (d *DriverMysql) GetChars() (charLeft string, charRight string) {
	return "`", "`"
}

// HandleSqlBeforeCommit handles the sql before posts it to database.
func (d *DriverMysql) HandleSqlBeforeCommit(link Link, sql string, args []interface{}) (string, []interface{}) {
	return sql, args
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *DriverMysql) Tables(schema ...string) (tables []string, err error) {
	var result Result
	link, err := d.db.GetSlave(schema...)
	if err != nil {
		return nil, err
	}
	result, err = d.db.DoGetAll(link, `SHOW TABLES`)
	if err != nil {
		return
	}
	for _, m := range result {
		for _, v := range m {
			tables = append(tables, v.String())
		}
	}
	return
}

// TableFields retrieves and returns the fields information of specified table of current
// schema.
//
// The parameter `link` is optional, if given nil it automatically retrieves a raw sql connection
// as its link to proceed necessary sql query.
//
// Note that it returns a map containing the field name and its corresponding fields.
// As a map is unsorted, the TableField struct has a "Index" field marks its sequence in
// the fields.
//
// It's using cache feature to enhance the performance, which is never expired util the
// process restarts.
func (d *DriverMysql) TableFields(link Link, table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.New("function TableFields supports only single table operations")
	}
	checkSchema := d.schema.Val()
	if len(schema) > 0 && schema[0] != "" {
		checkSchema = schema[0]
	}
	tableFieldsCacheKey := fmt.Sprintf(
		`mysql_table_fields_%s_%s@group:%s`,
		table, checkSchema, d.GetGroup(),
	)
	v := tableFieldsMap.GetOrSetFuncLock(tableFieldsCacheKey, func() interface{} {
		var (
			result Result
		)
		if link == nil {
			link, err = d.db.GetSlave(checkSchema)
			if err != nil {
				return nil
			}
		}
		result, err = d.db.DoGetAll(
			link,
			fmt.Sprintf(`SHOW FULL COLUMNS FROM %s`, d.db.QuoteWord(table)),
		)
		if err != nil {
			return nil
		}
		fields = make(map[string]*TableField)
		for i, m := range result {
			fields[m["Field"].String()] = &TableField{
				Index:   i,
				Name:    m["Field"].String(),
				Type:    m["Type"].String(),
				Null:    m["Null"].Bool(),
				Key:     m["Key"].String(),
				Default: m["Default"].Val(),
				Extra:   m["Extra"].String(),
				Comment: m["Comment"].String(),
			}
		}
		return fields
	})
	if v != nil {
		fields = v.(map[string]*TableField)
	}
	return
}
