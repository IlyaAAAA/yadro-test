package computerClub

import (
	"fmt"
	"log"
	"sort"
	"time"
	"yadro-test/cmd/manager"
	"yadro-test/cmd/queue"
	"yadro-test/cmd/util"
)

const (
	outClientLeave        = 11
	outClientTakeComputer = 12
	outError              = 13
)

const stub = 1

type TableInfo struct {
	TakeTime time.Time
	Client   string
}

type ComputerClub struct {
	openTime        time.Time
	closeTime       time.Time
	computersNumber int
	enteredMap      map[string]int
	tables          []TableInfo
	clientMap       map[string]int
	freeTables      int
	clientQueue     queue.Queue
	paymentManager  manager.TableManager
}

func NewComputerClub(openTime, closeTime time.Time, computersNumber int, paymentManager manager.TableManager) *ComputerClub {
	return &ComputerClub{
		openTime:        openTime,
		closeTime:       closeTime,
		computersNumber: computersNumber,
		enteredMap:      make(map[string]int),
		tables:          make([]TableInfo, computersNumber),
		clientMap:       make(map[string]int),
		freeTables:      computersNumber,
		clientQueue:     queue.New(),
		paymentManager:  paymentManager,
	}
}

func (c *ComputerClub) HandleClientArrived(eventTime time.Time, client string) { //+
	if _, ok := c.enteredMap[client]; ok {
		fmt.Println("YouShallNotPass")
	} else if eventTime.Before(c.openTime) || eventTime.After(c.closeTime) {
		fmt.Println(util.GetTimeStr(eventTime), outError, "NotOpenYet")
	} else {
		c.enteredMap[client] = stub
	}
}

func (c *ComputerClub) HandleClientTakingSeat(eventTime time.Time, client string, tableToSeat int) {
	if !c.isClientInClub(client) {
		fmt.Println("ClientUnknown")

		return
	}

	if c.isTableBusy(tableToSeat - 1) {
		fmt.Println(util.GetTimeStr(eventTime), outError, "PlaceIsBusy")

		return
	}

	if clientAlreadySatTable, isOk := c.clientMap[client]; isOk { // когда пересаживается
		c.tables[clientAlreadySatTable] = TableInfo{}

		c.clientMap[client] = tableToSeat - 1
		c.tables[tableToSeat-1] = TableInfo{TakeTime: eventTime, Client: client}
	} else { //когда никакой не занимает стол и это стол свободный (ток вошел и садится)
		c.clientMap[client] = tableToSeat - 1
		c.tables[tableToSeat-1] = TableInfo{TakeTime: eventTime, Client: client}

		c.freeTables--
	}
}

func (c *ComputerClub) HandleClientWaiting(eventTime time.Time, client string) {
	if c.freeTables < 0 {
		log.Fatalf("Неправильное состояние")
	}

	if !c.isClientInClub(client) { // логичнее, чтобы и тут проверка на unknown была
		fmt.Println("ClientUnknown")

		return
	}

	if c.freeTables != 0 {
		fmt.Println(util.GetTimeStr(eventTime), outError, "ICanWaitNoLonger!")

		return
	}

	if c.clientQueue.Len() > c.computersNumber {
		c.PrintCurrentClients(eventTime)

		return
	}

	c.clientQueue.Add(client)
}

func (c *ComputerClub) HandleClientLeaving(eventTime time.Time, client string) {
	if !c.isClientInClub(client) {
		fmt.Println("ClientUnknown")

		return
	}

	if tableNumber, clientOk := c.clientMap[client]; clientOk { // сидит за столом
		tableInfoToCalculate := c.tables[tableNumber]

		c.tables[tableNumber] = TableInfo{}
		delete(c.clientMap, client)
		delete(c.enteredMap, client)

		if clientFromQueue, err := c.clientQueue.Poll(); err == nil {
			c.tables[tableNumber] = TableInfo{TakeTime: eventTime, Client: clientFromQueue}
			c.clientMap[clientFromQueue] = tableNumber

			fmt.Println(util.GetTimeStr(eventTime), outClientTakeComputer, clientFromQueue, tableNumber+1)
		} else {
			c.freeTables++
		}

		c.paymentManager.CalculateTableCost(tableInfoToCalculate.TakeTime, eventTime, tableNumber)
	}
}

func (c *ComputerClub) CalculateAllTablesAtEnd() {
	for i := range c.tables {
		tableInfo := c.tables[i]

		if tableInfo.Client != "" {
			c.paymentManager.CalculateTableCost(tableInfo.TakeTime, c.closeTime, i)
		}
	}
}

func (c *ComputerClub) isTableBusy(tableNumber int) bool {
	el := c.tables[tableNumber]

	return el.Client != ""
}

func (c *ComputerClub) isClientInClub(client string) bool {
	_, ok := c.enteredMap[client]

	return ok
}

func (c *ComputerClub) PrintCurrentClients(time time.Time) {
	tmp := make([]TableInfo, len(c.tables))
	copy(tmp, c.tables)
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Client < tmp[j].Client
	})

	for i := range tmp {
		client := tmp[i].Client

		if client != "" {
			fmt.Println(util.GetTimeStr(time), outClientLeave, client)
		}
	}
}
