package scaner

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	computerclub "computerClub"
)

const (
	placeBusy            = "PlaceIsBusy"
	clientUnknown        = "ClientUnknown"
	notOpenYet           = "NotOpenYet"
	youShallNotPass      = "YouShallNotPass"
	iCanWaitNoLonger     = "ICanWaitNoLonger!"
	incorrectTableNumber = "IncorrectTableNumber"
	timeParsingFormat    = "15:04"
)

type ScanConsole struct {
	FileName            string
	compClub            computerclub.Club
	numberStrings       int
	dayBalance          map[int]computerclub.DayBalance
	mapActiveCliens     map[string]computerclub.ActiveClientCard
	mapFreedomTables    map[int]bool
	broker              []string
	numberFreedomTables int
}

func NewScanConsole(FileName string) *ScanConsole {
	return &ScanConsole{
		FileName:            FileName,
		compClub:            computerclub.Club{},
		numberStrings:       0,
		dayBalance:          map[int]computerclub.DayBalance{},
		mapActiveCliens:     map[string]computerclub.ActiveClientCard{},
		mapFreedomTables:    map[int]bool{},
		broker:              []string{},
		numberFreedomTables: 0,
	}
}

func (s *ScanConsole) Read() {
	file, err := os.Open(s.FileName)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		// return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(line)
		s.match(line, s.numberStrings)
		s.numberStrings++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при сканировании файла:", err)
	}
}

func (s *ScanConsole) match(st string, n int) {
	switch {
	case n == 0:
		s.firstString(st)
	case n == 1:
		s.secondString(st)
	case n == 2:
		s.thirdString(st)
	case n >= 3:
		fmt.Println(st)
		s.fourthString(st)
	}

}

func (s *ScanConsole) firstString(st string) {
	res, _ := strconv.Atoi(st)
	s.compClub.NumberOfTables = res
	s.numberFreedomTables = res
	s.broker = make([]string, 0)
	for i := 1; i < res+1; i++ {
		s.mapFreedomTables[i] = true
		s.dayBalance[i] = computerclub.DayBalance{}
	}
}
func (s *ScanConsole) secondString(st string) {
	sl := strings.Split(st, " ")
	startTime, err := time.Parse(timeParsingFormat, sl[0])
	if err != nil {
		fmt.Println("Ошибка при парсинге времени:", err)
	}
	endTime, err := time.Parse(timeParsingFormat, sl[1])
	if err != nil {
		fmt.Println("Ошибка при парсинге времени:", err)
	}
	s.compClub.StartTime = startTime
	s.compClub.EndTime = endTime
	fmt.Printf("%02d:%02d\n", s.compClub.StartTime.Hour(), s.compClub.StartTime.Minute())
}
func (s *ScanConsole) thirdString(st string) {
	res, _ := strconv.Atoi(st)
	s.compClub.Coast = res
}
func (s *ScanConsole) fourthString(st string) {
	sl := strings.Split(st, " ")
	if len(sl) == 4 {
		s.clientSit(sl)
	} else if len(sl) == 3 {
		s.clientNotSit(sl)
	}
}
func (s *ScanConsole) logicIfSecondType(numTable int, startTime time.Time, t string) bool {
	res := true
	if numTable > s.compClub.NumberOfTables || numTable <= 0 {
		fmt.Println(t, 13, incorrectTableNumber)
		res = false
	}
	if ok := s.mapFreedomTables[numTable]; !ok {
		fmt.Println(t, 13, placeBusy)
		res = false
	}
	if !res {
		return res
	}
	res = s.validateTime(startTime, t)
	return res
}
func (s *ScanConsole) validateTime(startTime time.Time, t string) bool {
	res := true
	if startTime.Before(s.compClub.StartTime) || startTime.After(s.compClub.EndTime) {
		fmt.Println(t, 13, notOpenYet)
		res = false
	}
	return res
}
func (s *ScanConsole) logicIfSecondName(name, t string) bool {
	res := true
	for _, value := range name {
		if !(value >= 48 && value <= 57) && !(value >= 65 && value <= 90) && !(value >= 97 && value <= 122) {
			fmt.Println(t, 13, clientUnknown)
			res = false
			break
		}
	}
	return res
}
func (s *ScanConsole) addClientToBroker(st, t string) {
	pr := s.broker
	inBrock := true
	for _, value := range pr {
		if value == st {
			inBrock = false
			break
		}
	}
	if inBrock {
		if len(pr) < s.compClub.NumberOfTables {
			pr = append(pr, st)
		} else {
			fmt.Println(t, 11, st)
		}
	}
	s.broker = pr
}
func ceilHour(t time.Duration) int {
	hour := t.Hours()
	return int(math.Ceil(hour))
}

func (s *ScanConsole) clientSit(sl []string) {
	numTable, _ := strconv.Atoi(sl[3])
	startTime, _ := time.Parse(timeParsingFormat, sl[0])
	if _, ok := s.mapActiveCliens[sl[2]]; !ok {
		//ещё нет
		first := s.logicIfSecondType(numTable, startTime, sl[0])
		second := s.logicIfSecondName(sl[2], sl[0])
		if first && second {
			s.mapActiveCliens[sl[2]] = computerclub.ActiveClientCard{
				StartTime:   startTime,
				MiddleTime:  startTime,
				TableNumber: numTable,
			}
			s.mapFreedomTables[numTable] = false
			s.numberFreedomTables--
			inBrock := false
			indexBrock := -1
			for index, value := range s.broker {
				if value == sl[2] {
					inBrock = true
					indexBrock = index
					break
				}
			}
			if inBrock {
				s.broker = append(s.broker[:indexBrock], s.broker[indexBrock+1:]...)
			}
		} else if !first && second {
			// если не удалось сесть, то добавляем в очередь если ещё не в очереди
			s.addClientToBroker(sl[2], sl[0])
		}
	} else {
		// если уже есть
		if s.logicIfSecondType(numTable, startTime, sl[0]) {
			pr := s.mapActiveCliens[sl[2]]
			s.mapFreedomTables[pr.TableNumber] = true
			s.mapFreedomTables[numTable] = false

			del := startTime.Sub(pr.MiddleTime)
			balance := s.dayBalance[pr.TableNumber]
			balance.Time = s.dayBalance[pr.TableNumber].Time.Add(del)
			balance.MoneyTime += ceilHour(del)
			balance.Money += s.compClub.Coast * ceilHour(del)
			s.dayBalance[pr.TableNumber] = balance

			pr.TableNumber = numTable
			pr.MiddleTime = startTime
			s.mapActiveCliens[sl[2]] = pr
		}
	}
}
func (s *ScanConsole) clientNotSit(sl []string) {
	typeS, _ := strconv.Atoi(sl[1])
	switch typeS {
	case 1:
		s.firstEventID(sl)
	case 3:
		s.thirdEventID(sl)
	case 4:
		s.fourthEventID(sl)
	}
}

func (s *ScanConsole) firstEventID(sl []string) {
	startTime, _ := time.Parse(timeParsingFormat, sl[0])
	if _, ok := s.mapActiveCliens[sl[2]]; ok {
		fmt.Println(sl[0], 13, youShallNotPass)
	}
	s.validateTime(startTime, sl[0])
}
func (s *ScanConsole) thirdEventID(sl []string) {
	if s.numberFreedomTables >= 1 {
		fmt.Println(sl[0], 13, iCanWaitNoLonger)
	}
	startTime, _ := time.Parse(timeParsingFormat, sl[0])
	t := s.validateTime(startTime, sl[0])
	nam := s.logicIfSecondName(sl[2], sl[0])
	if t && nam {
		s.addClientToBroker(sl[2], sl[0])
	}
}
func (s *ScanConsole) fourthEventID(sl []string) {
	endTime, _ := time.Parse(timeParsingFormat, sl[0])
	t := s.validateTime(endTime, sl[0])
	if _, ok := s.mapActiveCliens[sl[2]]; ok && t {
		pr := s.mapActiveCliens[sl[2]]
		pr.EntTime = endTime

		del := endTime.Sub(pr.MiddleTime)

		balance := s.dayBalance[pr.TableNumber]
		balance.Time = s.dayBalance[pr.TableNumber].Time.Add(del)
		balance.MoneyTime += ceilHour(del)
		balance.Money += s.compClub.Coast * ceilHour(del)
		s.dayBalance[pr.TableNumber] = balance

		if len(s.broker) > 0 {
			s.mapActiveCliens[s.broker[0]] = computerclub.ActiveClientCard{
				StartTime:   endTime,
				MiddleTime:  endTime,
				TableNumber: pr.TableNumber,
			}
			fmt.Println(sl[0], 12, s.broker[0], pr.TableNumber)
			s.broker = s.broker[1:]
		} else {
			s.mapFreedomTables[pr.TableNumber] = true
			s.numberFreedomTables++
		}
		delete(s.mapActiveCliens, sl[2])
	}
}

func (s *ScanConsole) Close() {
	if s.numberFreedomTables > 0 {
		for i := range s.mapActiveCliens {
			pr := s.mapActiveCliens[i]
			del := s.compClub.EndTime.Sub(pr.MiddleTime)

			balance := s.dayBalance[pr.TableNumber]
			balance.Time = s.dayBalance[pr.TableNumber].Time.Add(del)
			balance.MoneyTime += ceilHour(del)
			balance.Money += s.compClub.Coast * ceilHour(del)
			s.dayBalance[pr.TableNumber] = balance

			fmt.Println(
				fmt.Sprintf("%02d:%02d", s.compClub.EndTime.Hour(), s.compClub.EndTime.Minute()),
				11, i)
		}
	}
	fmt.Printf("%02d:%02d\n", s.compClub.EndTime.Hour(), s.compClub.EndTime.Minute())
	for i, j := range s.dayBalance {
		fmt.Println(i, j.Money, fmt.Sprintf("%02d:%02d", j.Time.Hour(), j.Time.Minute()))
	}
}
