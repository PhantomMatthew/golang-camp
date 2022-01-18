package jobs

import (
	"github.com/bamzi/jobrunner"
	"newline.com/newline/src/srv/customer/bll"
	"time"
)

func StartCustomerJob() {

	jobrunner.Start()
	//schedule, err := cron.ParseStandard("0 18 * * ?")
	//jobrunner.Schedule("0 10 * * ?", CustomTask{})
	jobrunner.Every(time.Minute*60, CustomTask{})
}

type CustomTask struct {
}

func (e CustomTask) Run() {
	currentTime := time.Now()
	today := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.Local)
	yesterday := today.AddDate(0, 0, -1)
	bll.NewCustBll().SyncHistoryInfos(nil, "dz_customer", "", 50, yesterday.Format("2006-01-02"), "")
}
