package timer

import (
	"fmt"
	"github.com/robfig/cron"
	"strings"
	"testing"
)

//每隔5秒执行一次：*/5 * * * * ?
//每隔1分钟执行一次：0 */1 * * * ?
//每天23点执行一次：0 0 23 * * ?
//每天凌晨1点执行一次：0 0 1 * * ?
//每月1号凌晨1点执行一次：0 0 1 1 * ?
//每月最后一天23点执行一次：0 0 23 L * ?
//每周星期天凌晨1点实行一次：0 0 1 ? * L
//在26分、29分、33分执行一次：0 26,29,33 * * * ?
//每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
func TestCron(t *testing.T) {
	/**
	编码建议：
	CreateCronFunc 内部的 func任务只是一个触发器（类似电路的继电器，使用多电控制强电）
	由func可以出发一个一步网络请求，或者一个函数的执行
	*/
	// 高级写发@every 5s  每5秒运行一次
	_, err := s.CreatCronFunc("@every 5s", func() {
		logger.Info("每5秒执行一次")
	})
	// 匹配规则类似 linux 中 crontab 但是多一位输入 首位表示秒
	cronId, err := s.CreatCronFunc("*/1 * * * * ?", func() {
		logger.Info("每秒执行1次")
	})

	_, err = s.CreatCronFunc("0 */1 * * * *", func() {
		logger.Warn("每分钟执行1次")
		// 删除一个定时任务
		s.DelCronFunc(cronId)
	})
	if err != nil {
		t.Fatal("CreatCronFunc error err=", err)
	} else {

	}

	select {}

}

func TestCronV1(t *testing.T) {
	c := cron.New()
	c.Start()
	err := c.AddFunc("0/5 * * * * ?", func() {
		logger.Info("0/5 * * * * ?")
	})
	if err != nil {
		t.Fatal("cron v1 addfunc error err=", err)
	}

	select {}
}

func TestSlice(t *testing.T) {
	a := []string{"1", "2", "3", "4", "5"}
	fmt.Println(a[len(a)-2:])
	fmt.Println(strings.Join(a, "/"))
}
