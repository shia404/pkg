package xgorm

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Config struct {
	Dsn         string   `json:"Dsn"`
	DsnReplicas []string `json:"DsnReplicas,optional"`
}

func NewGorm(c Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 判断是否集群
	if len(c.DsnReplicas) > 0 {
		var replicas []gorm.Dialector
		for _, replica := range c.DsnReplicas {
			replicas = append(replicas, mysql.Open(replica))
		}
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources:           []gorm.Dialector{mysql.Open(c.Dsn)},
			Replicas:          replicas,
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}))
		if err != nil {
			panic(err)
		}
	}

	return db
}

// FirstErr 单条查询错误信息处理
func FirstErr(err error) error {
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return nil
}

// DbErr 数据库的错误，方便做日志使用
func DbErr(err error) error {
	if err != nil {
		return err
	}
	return nil
}
