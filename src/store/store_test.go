package store

import (
	// "testing"
	"utils"
	// "models"
	// "time"
	"flag"
	// "github.com/golang/glog"
	"github.com/jackc/pgx"
)

var _s *StoreManager
var _conn *pgx.Conn

func init() {

	flag.Set("v", "1")
	flag.Set("stderrthreshold", "ERROR")

	utils.IsTesting = true
	utils.LoadConfig("../../config/config.json")
	_s = NewStore()
	
	_conn, _ = setupConnectionPGX(utils.Cfg.StorageSettings)
}