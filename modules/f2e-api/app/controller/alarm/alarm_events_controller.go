package alarm

import (
	"errors"
	"fmt"
	"math"
	"strings"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	alm "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/alarm"
	"gopkg.in/gin-gonic/gin.v1"
)

type APIGetAlarmListsInputs struct {
	StartTime     int64  `json:"startTime" form:"startTime"`
	EndTime       int64  `json:"endTime" form:"endTime"`
	Priority      int    `json:"priority" form:"priority"`
	Status        string `json:"status" form:"status"`
	ProcessStatus string `json:"process_status" form:"process_status"`
	Metrics       string `json:"metrics" form:"metrics"`
	//id
	EventId string `json:"event_id" form:"event_id"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
}

func (input APIGetAlarmListsInputs) checkInputsContain() error {
	if input.StartTime == 0 && input.EndTime == 0 {
		if input.EventId == "" {
			return errors.New("startTime, endTime OR event_id, You have to at least pick one on the request.")
		}
	}
	return nil
}

func (s APIGetAlarmListsInputs) collectFilters() string {
	tmp := []string{}
	if s.StartTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", s.StartTime))
	}
	if s.EndTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", s.EndTime))
	}
	if s.Priority != -1 {
		tmp = append(tmp, fmt.Sprintf("priority = %d", s.Priority))
	}
	if s.Status != "" {
		status := ""
		statusTmp := strings.Split(s.Status, ",")
		for indx, n := range statusTmp {
			if indx == 0 {
				status = fmt.Sprintf(" status = '%s' ", n)
			} else {
				status = fmt.Sprintf(" %s OR status = '%s' ", status, n)
			}
		}
		status = fmt.Sprintf("( %s )", status)
		tmp = append(tmp, status)
	}
	if s.ProcessStatus != "" {
		pstatus := ""
		pstatusTmp := strings.Split(s.ProcessStatus, ",")
		for indx, n := range pstatusTmp {
			if indx == 0 {
				pstatus = fmt.Sprintf(" process_status = '%s' ", n)
			} else {
				pstatus = fmt.Sprintf(" %s OR process_status = '%s' ", pstatus, n)
			}
		}
		pstatus = fmt.Sprintf("( %s )", pstatus)
		tmp = append(tmp, pstatus)
	}
	if s.Metrics != "" {
		tmp = append(tmp, fmt.Sprintf("metrics regexp '%s'", s.Metrics))
	}
	if s.EventId != "" {
		tmp = append(tmp, fmt.Sprintf("id = '%s'", s.EventId))
	}
	filterStrTmp := strings.Join(tmp, " AND ")
	if filterStrTmp != "" {
		filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
	}
	return filterStrTmp
}

func AlarmLists(c *gin.Context) {
	var inputs APIGetAlarmListsInputs
	//set default
	inputs.Page = -1
	inputs.Priority = -1
	inputs.Limit = 50
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	filterCollector := inputs.collectFilters()
	//for get correct table name
	f := alm.EventCases{}
	cevens := []alm.EventCases{}
	perparedSql := ""
	//if no specific, will give return first 2000 records
	if inputs.Page == -1 {
		if inputs.Limit >= 2000 || inputs.Limit == 0 {
			inputs.Limit = 2000
		}
		perparedSql = fmt.Sprintf("select * from %s %s order by timestamp DESC limit %d", f.TableName(), filterCollector, inputs.Limit)
		db.Alarm.Raw(perparedSql).Find(&cevens)
		h.JSONR(c, map[string]interface{}{
			"limit":    inputs.Limit,
			"priority": inputs.Priority,
			"data":     cevens,
		})
		return
	} else {
		//set the max limit of each page
		if inputs.Limit >= 50 {
			inputs.Limit = 50
		}
		perparedSql = fmt.Sprintf("select * from %s %s  order by timestamp DESC limit %d,%d", f.TableName(), filterCollector, inputs.Page, inputs.Limit)
		db.Alarm.Raw(perparedSql).Find(&cevens)
		var totalCount int64
		db.Alarm.Raw(fmt.Sprintf("select count(id) from %s %s ", f.TableName(), filterCollector)).Count(&totalCount)
		totalPage := math.Ceil(float64(totalCount) / float64(inputs.Limit))
		h.JSONR(c, map[string]interface{}{
			"total_count":  totalCount,
			"total_page":   totalPage,
			"current_page": inputs.Page,
			"limit":        inputs.Limit,
			"priority":     inputs.Priority,
			"data":         cevens,
		})
		return
	}
}

type APIEventsGetInputs struct {
	StartTime int64 `json:"startTime" form:"startTime"`
	EndTime   int64 `json:"endTime" form:"endTime"`
	Status    int   `json:"status" form:"status" binding:"gte=-1,lte=1"`
	//event_caseId
	EventId string `json:"event_id" form:"event_id" binding:"required"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
}

func (s APIEventsGetInputs) collectFilters() string {
	tmp := []string{}
	filterStrTmp := ""
	if s.StartTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", s.StartTime))
	}
	if s.EndTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", s.EndTime))
	}
	if s.EventId != "" {
		tmp = append(tmp, fmt.Sprintf("event_caseId = '%s'", s.EventId))
	}
	if s.Status == 0 || s.Status == 1 {
		tmp = append(tmp, fmt.Sprintf("status = %d", s.Status))
	}
	if len(tmp) != 0 {
		filterStrTmp = strings.Join(tmp, " AND ")
		filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
	}
	return filterStrTmp
}

func EventsGet(c *gin.Context) {
	var inputs APIEventsGetInputs
	inputs.Status = -1
	inputs.Page = -1
	inputs.Limit = 10
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	filterCollector := inputs.collectFilters()
	//for get correct table name
	f := alm.Events{}
	evens := []alm.Events{}
	perparedSql := fmt.Sprintf("select id, event_caseId, cond, status, timestamp from %s %s order by timestamp DESC limit %d,%d", f.TableName(), filterCollector, inputs.Page, inputs.Limit)
	db.Alarm.Raw(perparedSql).Scan(&evens)
	h.JSONR(c, evens)
}
