/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dtmsvr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	CronExpiredTrans(1)
	sleepCronTime()
}

func TestSetNextCron(t *testing.T) {
	tg := TransGlobal{}
	tg.RetryInterval = 15
	tg.setNextCron(cronReset)
	assert.Equal(t, int64(15), tg.NextCronInterval)
	tg.RetryInterval = 0
	tg.setNextCron(cronReset)
	assert.Equal(t, config.RetryInterval, tg.NextCronInterval)
	tg.setNextCron(cronBackoff)
	assert.Equal(t, config.RetryInterval*2, tg.NextCronInterval)
}
