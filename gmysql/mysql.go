package gmysql

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

var m_db *gorm.DB
var m_dbs map[string]*gorm.DB
var m_multi bool
var m_config IConfig

func Init(configMgr IConfig) {
	m_config = configMgr
	if m_db == nil && len(m_dbs) == 0 {
		m_multi = false
		if m_config.Bool("go.data.mysql.multidb") {
			m_multi = true
			m_dbs = make(map[string]*gorm.DB)
			dbNames := strings.Split(m_config.String("go.data.mysql.dbNames"), ",")
			for _, dbName := range dbNames {
				if dbName != "" && m_config.String("go.data.mysql."+dbName) != "" {
					conn, err := gorm.Open(mysql.Open(m_config.String("go.data.mysql."+dbName)), &gorm.Config{})
					if err != nil {
						log.Fatalln(dbName + " mysql connection error:" + err.Error())
					}
					m_dbs[dbName] = conn
				}
			}
		} else {
			m_db, _ = gorm.Open(mysql.Open(m_config.String("go.data.mysql.uri")), &gorm.Config{})
		}
		if m_config.Bool("go.data.mysql.debug") {
			if m_multi {
				for k, _ := range m_dbs {
					m_dbs[k] = m_dbs[k].Debug()
				}
			} else {
				m_db = m_db.Debug()
			}
		}
		if m_config.Int("go.data.mysql.pool.max") > 1 {
			max := m_config.Int("go.data.mysql_pool.max")
			if max < 10 {
				max = 10
			}
			idle := m_config.Int("go.data.mysql.pool.total")
			if idle == 0 || idle < max {
				idle = 5 * max
			}
			idleTimeout := m_config.Int("go.data.mysql.pool.timeout")
			if idleTimeout == 0 {
				idleTimeout = 60
			}
			lifetime := m_config.Int("go.data.mysql.pool.life")
			if lifetime == 0 {
				lifetime = 60
			}
			if m_multi {
				sqldb, _ := m_db.DB()
				sqldb.SetConnMaxIdleTime(time.Duration(idleTimeout) * time.Second)
				sqldb.SetMaxIdleConns(idle)
				sqldb.SetMaxOpenConns(max)
				sqldb.SetConnMaxLifetime(time.Duration(lifetime) * time.Minute)
			} else {
				for k, _ := range m_dbs {
					sqldb, _ := m_dbs[k].DB()
					sqldb.SetConnMaxIdleTime(time.Duration(idleTimeout) * time.Second)
					sqldb.SetMaxIdleConns(idle)
					sqldb.SetMaxOpenConns(max)
					sqldb.SetConnMaxLifetime(time.Duration(lifetime) * time.Minute)
				}
			}
		}
	}
}

func dbClose(db *gorm.DB) {
	sqldb, _ := db.DB()
	err := sqldb.Close()
	if err != nil {
		fmt.Printf("db close fail" + err.Error())
	}
	return
}

func dbReopen(dbName ...string) *gorm.DB {
	if len(dbName) == 0 {
		conn, err := gorm.Open(mysql.Open(m_config.String("go.data.mysql.uri")), &gorm.Config{})
		if err != nil {
			log.Fatalln(" mysql connection error:" + err.Error())
		}
		return conn
	}
	conn, err := gorm.Open(mysql.Open(m_config.String("go.data.mysql."+dbName[0])), &gorm.Config{})
	if err != nil {
		log.Fatalln(dbName[0] + " mysql connection error:" + err.Error())
	}
	return conn
}

func GetConnection(dbName ...string) (*gorm.DB, error) {
	if len(dbName) == 0 && m_multi {
		return nil, errors.New("multi get connection must specify a database name")
	}
	if len(dbName) > 1 {
		return nil, errors.New("Multidb can only get one connection")
	}
	if !m_multi {
		if m_db == nil {
			log.Fatal("mysql datasource not init")
		}
		sqldb, _ := m_db.DB()
		err := sqldb.Ping()
		if err != nil {
			dbClose(m_db)
			m_db = dbReopen()
			if m_db == nil {
				return nil, errors.New("mySQL connection error")
			}
		}
		return m_db, nil
	}

	if m_dbs[dbName[0]] == nil {
		return nil, errors.New(dbName[0] + " mysql datasource not init")
	}
	sqldb, _ := m_dbs[dbName[0]].DB()
	err := sqldb.Ping()
	if err != nil {
		dbClose(m_dbs[dbName[0]])
		m_dbs[dbName[0]] = dbReopen(dbName[0])
		if m_dbs[dbName[0]] == nil {
			return nil, errors.New("mySQL connection error")
		}
	}

	return m_dbs[dbName[0]], nil
}
