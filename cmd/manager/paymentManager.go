package manager

import (
	"fmt"
	"time"
	"yadro-test/cmd/util"
)

type TableManager struct {
	tableTotalInfo []tableTotalInfo
	costForHour    int
}

func NewPaymentManager(computersNumber, costForHour int) *TableManager {
	return &TableManager{
		tableTotalInfo: make([]tableTotalInfo, computersNumber),
		costForHour:    costForHour,
	}
}

func (m *TableManager) CalculateTableCost(startTime time.Time, endTime time.Time, tableNumber int) {
	sub := endTime.Sub(startTime)

	hours := int(sub.Hours())
	minutes := int(sub.Minutes())

	if hours >= 0 && minutes != 0 {
		hours++
	}

	m.tableTotalInfo[tableNumber].TotalSum += hours * m.costForHour
	m.tableTotalInfo[tableNumber].TotalTime = m.tableTotalInfo[tableNumber].TotalTime.Add(sub)
}

func (m *TableManager) PrintTotalInfo() {
	for i := range m.tableTotalInfo {
		tableTotalInfo := m.tableTotalInfo[i]

		fmt.Println(i+1, tableTotalInfo.TotalSum, util.GetTimeStr(tableTotalInfo.TotalTime))
	}
}

type tableTotalInfo struct {
	TotalSum  int
	TotalTime time.Time
}
