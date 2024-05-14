package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"yadro-test/cmd/computerClub"
	"yadro-test/cmd/manager"
	"yadro-test/cmd/util"
)

const hoursAndMinutes = "15:04"

const (
	inClientArrived      = 1
	inClientTakeComputer = 2
	inClientWait         = 3
	inClientLeave        = 4
)

type Event struct {
	eventTime   time.Time
	eventType   int
	client      string
	tableNumber int
}

func main() {
	args := os.Args

	if len(args) != 2 {
		log.Fatalf("Неправильное количество аргументов")
	}

	fileName := args[1]

	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		log.Fatalf(err.Error())
	}

	lineNumber := 1

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	text := scanner.Text()
	computersNumber, err := strconv.Atoi(text)
	if err != nil {
		log.Fatalf("Ошибка в строке номер %v", lineNumber)
	}
	lineNumber++

	scanner.Scan()
	text = scanner.Text()
	split := strings.Split(text, " ")
	openTime, err := time.Parse(hoursAndMinutes, split[0])
	if err != nil {
		log.Fatalf("Ошибка в строке номер %v", lineNumber)
	}

	closeTime, err := time.Parse(hoursAndMinutes, split[1])
	if err != nil {
		log.Fatalf("Ошибка в строке номер %v", lineNumber)
	}
	lineNumber++

	scanner.Scan()
	text = scanner.Text()
	costForHour, err := strconv.Atoi(text)

	if err != nil {
		log.Fatalf("Ошибка в строке номер %v", lineNumber)
	}
	lineNumber++

	events := make([]Event, 0)

	paymentManager := *manager.NewPaymentManager(computersNumber, costForHour)
	compClub := computerClub.NewComputerClub(openTime, closeTime, computersNumber, paymentManager)

	fmt.Println(util.GetTimeStr(openTime))

	for scanner.Scan() {
		text = scanner.Text()
		split := strings.Split(text, " ")

		if len(split) == 0 {
			log.Fatalf("Ошибка в строке номер %v", lineNumber)
		}

		eventTime, err := time.Parse(hoursAndMinutes, split[0])
		if err != nil {
			log.Fatalf("Ошибка в строке номер %v", lineNumber)
		}

		event, err := strconv.Atoi(split[1])
		if err != nil {
			log.Fatalf("Ошибка в строке номер %v", lineNumber)
		}

		client := split[2]

		if len(split) == 3 {
			events = append(events, Event{
				eventTime:   eventTime,
				eventType:   event,
				client:      client,
				tableNumber: -1,
			})

		} else if len(split) == 4 {
			tableNumber, err := strconv.Atoi(split[3])
			if err != nil {
				log.Fatalf("Ошибка в строке номер %v", lineNumber)
			}

			events = append(events, Event{
				eventTime:   eventTime,
				eventType:   event,
				client:      client,
				tableNumber: tableNumber,
			})
		} else {
			log.Fatalf("Ошибка в строке номер %v", lineNumber)
		}

		lineNumber++
	}

	for _, event := range events {
		if event.eventType != inClientTakeComputer {
			fmt.Println(util.GetTimeStr(event.eventTime), event.eventType, event.client)

		} else {
			fmt.Println(util.GetTimeStr(event.eventTime), event.eventType, event.client, event.tableNumber)
		}

		switch event.eventType {
		case inClientArrived:
			compClub.HandleClientArrived(event.eventTime, event.client)
		case inClientTakeComputer:
			compClub.HandleClientTakingSeat(event.eventTime, event.client, event.tableNumber)
		case inClientWait:
			compClub.HandleClientWaiting(event.eventTime, event.client)
		case inClientLeave:
			compClub.HandleClientLeaving(event.eventTime, event.client)
		default:
		}
	}

	compClub.PrintCurrentClients(closeTime)
	fmt.Println(util.GetTimeStr(closeTime))

	compClub.CalculateAllTablesAtEnd()
	paymentManager.PrintTotalInfo()
}
