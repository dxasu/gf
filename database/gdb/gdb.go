// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/dxasu/gf.

// Package gdb provides ORM features for popular relationship databases.
package gdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dxasu/gf/errors/gerror"
	"github.com/dxasu/gf/os/gcmd"

	"github.com/dxasu/gf/container/gvar"
	"github.com/dxasu/gf/internal/intlog"

	"github.com/dxasu/gf/os/glog"

	"github.com/dxasu/gf/container/gmap"
	"github.com/dxasu/gf/container/gtype"
	"github.com/dxasu/gf/os/gcache"
	"github.com/dxasu/gf/util/grand"
)

// DB defines the interfaces for ORM operations.
type DB interface {
	// ===========================================================================
	// Model creation.
	// ===========================================================================

	// The DB interface is designed not only for
	// relational databases but also for NoSQL databases in the future. The name
	// "Table" is not proper for that purpose any more.
	// Also see  Core.Table.
	// Deprecated, use Model instead.
	Table(table ...string) *Model

	// Model creates and returns a new ORM model from given schema.
	// The parameter `table` can be more than one table names, and also alias name, like:
	// 1. Model names:
	//    Model("user")
	//    Model("user u")
	//    Model("user, user_detail")
	//    Model("user u, user_detail ud")
	// 2. Model name with alias: Model("user", "u")
	// Also see Core.Model.
	Model(table ...string) *Model

	// Schema creates and returns a schema.
	// Also see Core.Schema.
	Schema(schema string) *Schema

	// With creates and returns an ORM model based on meta data of given object.
	// Also see Core.With.
	With(object interface{}) *Model

	// Open creates a raw connection object for database with given node configuration.
	// Note that it is not recommended using the this function manually.
	// Also see DriverMysql.Open.
	Open(config *ConfigNode) (*sql.DB, error)

	// Ctx is a chaining function, which creates and returns a new DB that is a shallow copy
	// of current DB object and with given context in it.
	// Note that this returned DB object can be used only once, so do not assign it to
	// a global or package variable for long using.
	// Also see Core.Ctx.
	Ctx(ctx context.Context) DB

	// ===========================================================================
	// Query APIs.
	// ===========================================================================

	Query(sql string, args ...interface{}) (*sql.Rows, error) // See Core.Query.
	Exec(sql string, args ...interface{}) (sql.Result, error) // See Core.Exec.
	Prepare(sql string, execOnMaster ...bool) (*Stmt, error)  // See Core.Prepare.

	// ===========================================================================
	// Common APIs for CURD.
	// ===========================================================================

	Insert(table string, data interface{}, batch ...int) (sql.Result, error)       // See Core.Insert.
	InsertIgnore(table string, data interface{}, batch ...int) (sql.Result, error) // See Core.InsertIgnore.
	Replace(table string, data interface{}, batch ...int) (sql.Result, error)      // See Core.Replace.
	Save(table string, data interface{}, batch ...int) (sql.Result, error)         // See Core.Save.

	BatchInsert(table string, list interface{}, batch ...int) (sql.Result, error)  // See Core.BatchInsert.
	BatchReplace(table string, list interface{}, batch ...int) (sql.Result, error) // See Core.BatchReplace.
	BatchSave(table string, list interface{}, batch ...int) (sql.Result, error)    // See Core.BatchSave.

	Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) // See Core.Update.
	Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error)                   // See Core.Delete.

	// ===========================================================================
	// Internal APIs for CURD, which can be overwrote for custom CURD implements.
	// ===========================================================================

	DoQuery(link Link, sql string, args ...interface{}) (rows *sql.Rows, err error)                                           // See Core.DoQuery.
	DoGetAll(link Link, sql string, args ...interface{}) (result Result, err error)                                           // See Core.DoGetAll.
	DoExec(link Link, sql string, args ...interface{}) (result sql.Result, err error)                                         // See Core.DoExec.
	DoPrepare(link Link, sql string) (*Stmt, error)                                                                           // See Core.DoPrepare.
	DoInsert(link Link, table string, data interface{}, option int, batch ...int) (result sql.Result, err error)              // See Core.DoInsert.
	DoBatchInsert(link Link, table string, list interface{}, option int, batch ...int) (result sql.Result, err error)         // See Core.DoBatchInsert.
	DoUpdate(link Link, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) // See Core.DoUpdate.
	DoDelete(link Link, table string, condition string, args ...interface{}) (result sql.Result, err error)                   // See Core.DoDelete.

	// ===========================================================================
	// Query APIs for convenience purpose.
	// ===========================================================================

	GetAll(sql string, args ...interface{}) (Result, error)                        // See Core.GetAll.
	GetOne(sql string, args ...interface{}) (Record, error)                        // See Core.GetOne.
	GetValue(sql string, args ...interface{}) (Value, error)                       // See Core.GetValue.
	GetArray(sql string, args ...interface{}) ([]Value, error)                     // See Core.GetArray.
	GetCount(sql string, args ...interface{}) (int, error)                         // See Core.GetCount.
	GetStruct(objPointer interface{}, sql string, args ...interface{}) error       // See Core.GetStruct.
	GetStructs(objPointerSlice interface{}, sql string, args ...interface{}) error // See Core.GetStructs.
	GetScan(objPointer interface{}, sql string, args ...interface{}) error         // See Core.GetScan.

	// ===========================================================================
	// Master/Slave specification support.
	// ===========================================================================

	Master() (*sql.DB, error) // See Core.Master.
	Slave() (*sql.DB, error)  // See Core.Slave.

	// ===========================================================================
	// Ping-Pong.
	// ===========================================================================

	PingMaster() error // See Core.PingMaster.
	PingSlave() error  // See Core.PingSlave.

	// ===========================================================================
	// Transaction.
	// ===========================================================================

	Begin() (*TX, error)                          // See Core.Begin.
	Transaction(f func(tx *TX) error) (err error) // See Core.Transaction.

	// ===========================================================================
	// Configuration methods.
	// ===========================================================================

	GetCache() *gcache.Cache            // See Core.GetCache.
	SetDebug(debug bool)                // See Core.SetDebug.
	GetDebug() bool                     // See Core.GetDebug.
	SetSchema(schema string)            // See Core.SetSchema.
	GetSchema() string                  // See Core.GetSchema.
	GetPrefix() string                  // See Core.GetPrefix.
	GetGroup() string                   // See Core.GetGroup.
	SetDryRun(enabled bool)             // See Core.SetDryRun.
	GetDryRun() bool                    // See Core.GetDryRun.
	SetLogger(logger *glog.Logger)      // See Core.SetLogger.
	GetLogger() *glog.Logger            // See Core.GetLogger.
	GetConfig() *ConfigNode             // See Core.GetConfig.
	SetMaxIdleConnCount(n int)          // See Core.SetMaxIdleConnCount.
	SetMaxOpenConnCount(n int)          // See Core.SetMaxOpenConnCount.
	SetMaxConnLifeTime(d time.Duration) // See Core.SetMaxConnLifeTime.

	// ===========================================================================
	// Utility methods.
	// ===========================================================================

	GetCtx() context.Context                                                               // See Core.GetCtx.
	GetChars() (charLeft string, charRight string)                                         // See Core.GetChars.
	GetMaster(schema ...string) (*sql.DB, error)                                           // See Core.GetMaster.
	GetSlave(schema ...string) (*sql.DB, error)                                            // See Core.GetSlave.
	QuoteWord(s string) string                                                             // See Core.QuoteWord.
	QuoteString(s string) string                                                           // See Core.QuoteString.
	QuotePrefixTableName(table string) string                                              // See Core.QuotePrefixTableName.
	Tables(schema ...string) (tables []string, err error)                                  // See Core.Tables.
	TableFields(link Link, table string, schema ...string) (map[string]*TableField, error) // See Core.TableFields.
	HasTable(name string) (bool, error)                                                    // See Core.HasTable.
	FilteredLinkInfo() string                                                              // See Core.FilteredLinkInfo.

	// HandleSqlBeforeCommit is a hook function, which deals with the sql string before
	// it's committed to underlying driver. The parameter `link` specifies the current
	// database connection operation object. You can modify the sql string `sql` and its
	// arguments `args` as you wish before they're committed to driver.
	// Also see Core.HandleSqlBeforeCommit.
	HandleSqlBeforeCommit(link Link, sql string, args []interface{}) (string, []interface{})

	// ===========================================================================
	// Internal methods, for internal usage purpose, you do not need consider it.
	// ===========================================================================

	mappingAndFilterData(schema, table string, data map[string]interface{}, filter bool) (map[string]interface{}, error) // See Core.mappingAndFilterData.
	convertFieldValueToLocalValue(fieldValue interface{}, fieldType string) interface{}                                  // See Core.convertFieldValueToLocalValue.
	convertRowsToResult(rows *sql.Rows) (Result, error)                                                                  // See Core.convertRowsToResult.
}

// Core is the base struct for database management.
type Core struct {
	db     DB              // DB interface object.
	ctx    context.Context // Context for chaining operation only.
	group  string          // Configuration group name.
	debug  *gtype.Bool     // Enable debug mode for the database, which can be changed in runtime.
	cache  *gcache.Cache   // Cache manager, SQL result cache only.
	schema *gtype.String   // Custom schema for this object.
	logger *glog.Logger    // Logger.
	config *ConfigNode     // Current config node.
}

// Driver is the interface for integrating sql drivers into package gdb.
type Driver interface {
	// New creates and returns a database object for specified database server.
	New(core *Core, node *ConfigNode) (DB, error)
}

// Sql is the sql recording struct.
type Sql struct {
	Sql    string        // SQL string(may contain reserved char '?').
	Type   string        // SQL operation type.
	Args   []interface{} // Arguments for this sql.
	Format string        // Formatted sql which contains arguments in the sql.
	Error  error         // Execution result.
	Start  int64         // Start execution timestamp in milliseconds.
	End    int64         // End execution timestamp in milliseconds.
	Group  string        // Group is the group name of the configuration that the sql is executed from.
}

// TableField is the struct for table field.
type TableField struct {
	Index   int         // For ordering purpose as map is unordered.
	Name    string      // Field name.
	Type    string      // Field type.
	Null    bool        // Field can be null or not.
	Key     string      // The index information(empty if it's not a index).
	Default interface{} // Default value for the field.
	Extra   string      // Extra information.
	Comment string      // Comment.
}

// Link is a common database function wrapper interface.
type Link interface {
	Query(sql string, args ...interface{}) (*sql.Rows, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Prepare(sql string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, sql string) (*sql.Stmt, error)
}

// Counter  is the type for update count.
type Counter struct {
	Field string
	Value float64
}

type (
	Raw    string                   // Raw is a raw sql that will not be treated as argument but as a direct sql part.
	Value  = *gvar.Var              // Value is the field value type.
	Record map[string]Value         // Record is the row record of the table.
	Result []Record                 // Result is the row record array.
	Map    = map[string]interface{} // Map is alias of map[string]interface{}, which is the most common usage map type.
	List   = []Map                  // List is type of map array.
)

const (
	insertOptionDefault     = 0
	insertOptionReplace     = 1
	insertOptionSave        = 2
	insertOptionIgnore      = 3
	defaultBatchNumber      = 10               // Per count for batch insert/replace/save.
	defaultMaxIdleConnCount = 10               // Max idle connection count in pool.
	defaultMaxOpenConnCount = 100              // Max open connection count in pool.
	defaultMaxConnLifeTime  = 30 * time.Second // Max life time for per connection in pool in seconds.
	ctxTimeoutTypeExec      = iota
	ctxTimeoutTypeQuery
	ctxTimeoutTypePrepare
)

var (
	// ErrNoRows is alias of sql.ErrNoRows.
	ErrNoRows = sql.ErrNoRows

	// instances is the management map for instances.
	instances = gmap.NewStrAnyMap(true)

	// driverMap manages all custom registered driver.
	driverMap = map[string]Driver{
		"mysql":  &DriverMysql{},
		"mssql":  &DriverMssql{},
		"pgsql":  &DriverPgsql{},
		"oracle": &DriverOracle{},
		"sqlite": &DriverSqlite{},
	}

	// lastOperatorRegPattern is the regular expression pattern for a string
	// which has operator at its tail.
	lastOperatorRegPattern = `[<>=]+\s*$`

	// regularFieldNameRegPattern is the regular expression pattern for a string
	// which is a regular field name of table.
	regularFieldNameRegPattern = `^[\w\.\-\_]+$`

	// regularFieldNameWithoutDotRegPattern is similar to regularFieldNameRegPattern but not allows '.'.
	// Note that, although some databases allow char '.' in the field name, but it here does not allow '.'
	// in the field name as it conflicts with "db.table.field" pattern in SOME situations.
	regularFieldNameWithoutDotRegPattern = `^[\w\-\_]+$`

	// internalCache is the memory cache for internal usage.
	internalCache = gcache.New()

	// tableFieldsMap caches the table information retrived from database.
	tableFieldsMap = gmap.New(true)

	// allDryRun sets dry-run feature for all database connections.
	// It is commonly used for command options for convenience.
	allDryRun = false
)

func init() {
	// allDryRun is initialized from environment or command options.
	allDryRun = gcmd.GetOptWithEnv("gf.gdb.dryrun", false).Bool()
}

// Register registers custom database driver to gdb.
func Register(name string, driver Driver) error {
	driverMap[name] = driver
	return nil
}

// New creates and returns an ORM object with global configurations.
// The parameter `name` specifies the configuration group name,
// which is DefaultGroupName in default.
func New(group ...string) (db DB, err error) {
	groupName := configs.group
	if len(group) > 0 && group[0] != "" {
		groupName = group[0]
	}
	configs.RLock()
	defer configs.RUnlock()

	if len(configs.config) < 1 {
		return nil, gerror.New("empty database configuration")
	}
	if _, ok := configs.config[groupName]; ok {
		if node, err := getConfigNodeByGroup(groupName, true); err == nil {
			c := &Core{
				group:  groupName,
				debug:  gtype.NewBool(),
				cache:  gcache.New(),
				schema: gtype.NewString(),
				logger: glog.New(),
				config: node,
			}
			if v, ok := driverMap[node.Type]; ok {
				c.db, err = v.New(c, node)
				if err != nil {
					return nil, err
				}
				return c.db, nil
			} else {
				return nil, gerror.New(fmt.Sprintf(`unsupported database type "%s"`, node.Type))
			}
		} else {
			return nil, err
		}
	} else {
		return nil, gerror.New(fmt.Sprintf(`database configuration node "%s" is not found`, groupName))
	}
}

// Instance returns an instance for DB operations.
// The parameter `name` specifies the configuration group name,
// which is DefaultGroupName in default.
func Instance(name ...string) (db DB, err error) {
	group := configs.group
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	v := instances.GetOrSetFuncLock(group, func() interface{} {
		db, err = New(group)
		return db
	})
	if v != nil {
		return v.(DB), nil
	}
	return
}

// getConfigNodeByGroup calculates and returns a configuration node of given group. It
// calculates the value internally using weight algorithm for load balance.
//
// The parameter `master` specifies whether retrieving a master node, or else a slave node
// if master-slave configured.
func getConfigNodeByGroup(group string, master bool) (*ConfigNode, error) {
	if list, ok := configs.config[group]; ok {
		// Separates master and slave configuration nodes array.
		masterList := make(ConfigGroup, 0)
		slaveList := make(ConfigGroup, 0)
		for i := 0; i < len(list); i++ {
			if list[i].Role == "slave" {
				slaveList = append(slaveList, list[i])
			} else {
				masterList = append(masterList, list[i])
			}
		}
		if len(masterList) < 1 {
			return nil, gerror.New("at least one master node configuration's need to make sense")
		}
		if len(slaveList) < 1 {
			slaveList = masterList
		}
		if master {
			return getConfigNodeByWeight(masterList), nil
		} else {
			return getConfigNodeByWeight(slaveList), nil
		}
	} else {
		return nil, gerror.New(fmt.Sprintf("empty database configuration for item name '%s'", group))
	}
}

// getConfigNodeByWeight calculates the configuration weights and randomly returns a node.
//
// Calculation algorithm brief:
// 1. If we have 2 nodes, and their weights are both 1, then the weight range is [0, 199];
// 2. Node1 weight range is [0, 99], and node2 weight range is [100, 199], ratio is 1:1;
// 3. If the random number is 99, it then chooses and returns node1;
func getConfigNodeByWeight(cg ConfigGroup) *ConfigNode {
	if len(cg) < 2 {
		return &cg[0]
	}
	var total int
	for i := 0; i < len(cg); i++ {
		total += cg[i].Weight * 100
	}
	// If total is 0 means all of the nodes have no weight attribute configured.
	// It then defaults each node's weight attribute to 1.
	if total == 0 {
		for i := 0; i < len(cg); i++ {
			cg[i].Weight = 1
			total += cg[i].Weight * 100
		}
	}
	// Exclude the right border value.
	r := grand.N(0, total-1)
	min := 0
	max := 0
	for i := 0; i < len(cg); i++ {
		max = min + cg[i].Weight*100
		//fmt.Printf("r: %d, min: %d, max: %d\n", r, min, max)
		if r >= min && r < max {
			return &cg[i]
		} else {
			min = max
		}
	}
	return nil
}

// getSqlDb retrieves and returns a underlying database connection object.
// The parameter `master` specifies whether retrieves master node connection if
// master-slave nodes are configured.
func (c *Core) getSqlDb(master bool, schema ...string) (sqlDb *sql.DB, err error) {
	// Load balance.
	node, err := getConfigNodeByGroup(c.group, master)
	if err != nil {
		return nil, err
	}
	// Default value checks.
	if node.Charset == "" {
		node.Charset = "utf8"
	}
	// Changes the schema.
	nodeSchema := c.schema.Val()
	if len(schema) > 0 && schema[0] != "" {
		nodeSchema = schema[0]
	}
	if nodeSchema != "" {
		// Value copy.
		n := *node
		n.Name = nodeSchema
		node = &n
	}
	// Cache the underlying connection pool object by node.
	v, _ := internalCache.GetOrSetFuncLock(node.String(), func() (interface{}, error) {
		intlog.Printf(
			`open new connection, master:%#v, config:%#v, node:%#v`,
			master, c.config, node,
		)
		defer func() {
			if err != nil {
				intlog.Printf(`open new connection failed: %v, %#v`, err, node)
			} else {
				intlog.Printf(
					`open new connection success, master:%#v, config:%#v, node:%#v`,
					master, c.config, node,
				)
			}
		}()

		sqlDb, err = c.db.Open(node)
		if err != nil {
			return nil, err
		}

		if c.config.MaxIdleConnCount > 0 {
			sqlDb.SetMaxIdleConns(c.config.MaxIdleConnCount)
		} else {
			sqlDb.SetMaxIdleConns(defaultMaxIdleConnCount)
		}
		if c.config.MaxOpenConnCount > 0 {
			sqlDb.SetMaxOpenConns(c.config.MaxOpenConnCount)
		} else {
			sqlDb.SetMaxOpenConns(defaultMaxOpenConnCount)
		}
		if c.config.MaxConnLifeTime > 0 {
			// Automatically checks whether MaxConnLifetime is configured using string like: "30s", "60s", etc.
			// Or else it is configured just using number, which means value in seconds.
			if c.config.MaxConnLifeTime > time.Second {
				sqlDb.SetConnMaxLifetime(c.config.MaxConnLifeTime)
			} else {
				sqlDb.SetConnMaxLifetime(c.config.MaxConnLifeTime * time.Second)
			}
		} else {
			sqlDb.SetConnMaxLifetime(defaultMaxConnLifeTime)
		}
		return sqlDb, nil
	}, 0)
	if v != nil && sqlDb == nil {
		sqlDb = v.(*sql.DB)
	}
	if node.Debug {
		c.db.SetDebug(node.Debug)
	}
	if node.DryRun {
		c.db.SetDryRun(node.DryRun)
	}
	return
}