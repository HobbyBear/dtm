/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package test

import (
	"os"
	"testing"
	"time"

	"github.com/dtm-labs/dtm/dtmcli"
	"github.com/dtm-labs/dtm/dtmcli/logger"
	"github.com/dtm-labs/dtm/dtmsvr"
	"github.com/dtm-labs/dtm/dtmsvr/config"
	"github.com/dtm-labs/dtm/test/busi"
)

func exitIf(code int) {
	if code != 0 {
		os.Exit(code)
	}
}

func TestMain(m *testing.M) {
	config.MustLoadConfig("")
	logger.InitLog("debug")
	dtmcli.SetCurrentDBType(busi.BusiConf.Driver)
	dtmsvr.TransProcessedTestChan = make(chan string, 1)
	dtmsvr.NowForwardDuration = 0 * time.Second
	dtmsvr.CronForwardDuration = 180 * time.Second
	conf.UpdateBranchSync = 1

	tenv := os.Getenv("TEST_STORE")
	if tenv == "boltdb" {
		conf.Store.Driver = "boltdb"
	} else if tenv == "mysql" {
		conf.Store.Driver = "mysql"
		conf.Store.Host = "localhost"
		conf.Store.Port = 3306
		conf.Store.User = "root"
		conf.Store.Password = ""
	} else {
		conf.Store.Driver = "redis"
		conf.Store.Host = "localhost"
		conf.Store.User = ""
		conf.Store.Password = ""
		conf.Store.Port = 6379
	}
	dtmsvr.PopulateDB(false)
	go dtmsvr.StartSvr()

	busi.PopulateDB(false)
	_ = busi.Startup()
	exitIf(m.Run())

}
