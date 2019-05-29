package clock

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Input interface {
	InitInput(string,<- chan time.Time)
	UpdateInput(string)
}

//Field for text
type Field struct {
	SecondText string
	MinuteText string
	HourText   string
	sync.Mutex
}


//UpdateInput check input values from file and updates the print text
func (f *Field) UpdateInput(input string) {

	str := strings.Split(input, " ")

	if len(str) > 2 || len(str) <= 1 {
		return
	}
	iv := strings.TrimSpace(str[1])
	switch str[0] {
	case "--sec":
		if iv != "" && iv != f.SecondText {
			f.Lock()
			f.SecondText = iv
			f.Unlock()

		}
		return

	case "--min":
		if iv != "" && iv != f.MinuteText {
			f.Lock()
			f.MinuteText = strings.TrimSpace(str[1])
			f.Unlock()
		}

		return

	case "--hour":
		if iv != "" && iv != f.HourText {
			f.Lock()
			f.HourText = iv
			f.Unlock()
		}

		return
	}
}

// InitInput initiates an input after 10 min
func (f *Field) InitInput(fname string,startInputTime <- chan time.Time) {
	for {
		select {
		case <-startInputTime:
			StartInputListner(fname, f)
		}
	}

}

// StartInputLintner is init input after 10 min
func StartInputListner(fname string, it Input) error {

	for {
		//mutex for file access
		mu := sync.Mutex{}
		f, err := os.Open(fname)

		if err != nil {
			f, err = os.Create(fname)
			if err != nil {
				return err
			}

		}

		mu.Lock()
		fmt.Fprintln(f, "--sec")
		fmt.Fprintln(f, "--min")
		fmt.Fprintln(f, "--hour")
		mu.Unlock()
		rd := bufio.NewReader(f)
		for {
			line, err := rd.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}

			}
			it.UpdateInput(line)
		}

		f.Close()
		time.Sleep(100 * time.Millisecond)


	}
	return nil

}

//ClConfig is config for clock
type ClConfig struct {
	Txt *Field
	Wg *sync.WaitGroup
	Ch chan string
	EndCh <-chan time.Time
	MinC *time.Ticker
	SecC *time.Ticker
	HourC *time.Ticker
}

// Clock starts a clock for hour ,minute and sec
func(cl *ClConfig) Clock() {

	sec := 0
	min := 0


	for {
		select {

		case <-cl.HourC.C:

			cl.Ch <- fmt.Sprintf("%s\n", cl.Txt.HourText)

		case <-cl.MinC.C:
			min = min + 1
			// print only when no one else is printing
			if min%60 != 0 {
				cl.Ch <- fmt.Sprintf("%s\n", cl.Txt.MinuteText)

			}
		case <-cl.SecC.C:
			sec = sec + 1
			// print only when no one else is printing
			if sec%60 != 0 {
				cl.Ch <- fmt.Sprintf("%s\n", cl.Txt.SecondText)
			}

		case <-cl.EndCh:
			fmt.Fprintf(os.Stdout, "%s\n", "clock exiting ")
			cl.Wg.Done()
			return

		}
	}
}
