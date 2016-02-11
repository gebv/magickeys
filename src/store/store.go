package store

import (
	"github.com/golang/glog"

	_ "github.com/lib/pq"

	"models"
	"utils"

	"strings"
	"github.com/jackc/pgx"
)

// var db *sql.DB
var db *pgx.Conn

var registredStores = make(map[string]StoreBuilder)

func registrationOfStoreBuilder(name string, builder StoreBuilder) {
	registredStores[name] = builder
}

type StoreBuilder func(*StoreManager) Store

func NewStore() *StoreManager {
	_sm := &StoreManager{}
	var err error

	_sm.db, err = setupConnectionPGX(utils.Cfg.StorageSettings)
	if err != nil {
		panic("error setup connection, err="+err.Error())
	}
	db = _sm.db

	_sm.stores = make(map[string]Store)

	for _name, _builder := range registredStores {
		_sm.stores[_name] = _builder(_sm)
	}

	return _sm
}

type Store interface {
	Name() string
}

type StoreManager struct {
	db   *pgx.Conn

	stores map[string]Store
}

// Get получить Store по имени
func (sm *StoreManager) Get(name string) Store {
	return sm.stores[name]
}

func (sm StoreManager) ErrorLog(prefix string, args ...interface{}) {
	if len(args)%2 == 0 {
		glog.Errorf(prefix+": "+strings.Repeat("%v='%v', ", len(args)/2), args...)
		return
	}

	glog.Errorf(prefix+": "+strings.Repeat("%v, ", len(args)), args...)
}

func setupConnectionPGX(c models.StorageSettings) (*pgx.Conn, error) {
	// config := pgx.ConnPoolConfig{ConnConfig: extractPGXStorageConfig(c), MaxConnections: 20}
	// pool, err := pgx.NewConnPool(config)
	return pgx.Connect(extractPGXStorageConfig(c))
}

type databaseLogger struct {}

func (l databaseLogger) Debug(msg string, ctx ...interface{}) {
	glog.Infof("\tSQL[debug]: msg='%s', %v", msg, ctx)
}

func (l databaseLogger) Info(msg string, ctx ...interface{}) {
	glog.Infof("\tSQL: msg='%s', %v", msg, ctx)
}

func (l databaseLogger) Warn(msg string, ctx ...interface{}) {
	glog.Warningf("\tSQL] msg='%s', %v", msg, ctx)
}

func (l databaseLogger) Error(msg string, ctx ...interface{}) {
	glog.Errorf("\tSQL] msg='%s', %v", msg, ctx)
}

func extractPGXStorageConfig(c models.StorageSettings) pgx.ConnConfig {
	var config pgx.ConnConfig

	config.Host = c.Host
	config.User = c.User
	config.Password = c.Password
	config.Database = c.Database
	config.Logger = databaseLogger{}
	config.LogLevel = pgx.LogLevelError

	return config
}
