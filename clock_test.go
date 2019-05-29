package clock_test

import (
	"github.com/sabina544/awesomeProject/clock"
	"os"
	"sync"
	"testing"
	"time"
)

func TestClock(t *testing.T) {

	var text = &clock.Field{SecondText: "tick", MinuteText: "tock", HourText: "bong"}
	wg := &sync.WaitGroup{}
	secOut := make(chan string)
	wg.Add(1)

	// create a after channel to end the execution
	endCh := time.After(time.Hour * time.Duration(3))

	// ticker for hour
	hourC := time.NewTicker(time.Hour * time.Duration(1))

	// ticker for minute
	minC := time.NewTicker(time.Minute * time.Duration(1))

	//  ticker for second
	secC := time.NewTicker(time.Second * time.Duration(1))
	clockConfig := clock.ClConfig{
		EndCh: endCh,
		Txt:   text,
		Wg:    wg,
		HourC: hourC,
		MinC:  minC,
		SecC:  secC,
		Ch:    secOut,
	}
	go clockConfig.Clock()
	secTexts := []string{}
	wg.Add(1)

	go func() {

		for {
			select {
			case <-secOut:
				secTexts = append(secTexts, <-secOut)

			}

			select {
			case <-secC.C:
				if len(secTexts) > 1 {
					t.Fatalf("lenth of received messages is higher  expected %d got %v ", 1, len(secTexts))
				}
				t.Logf("sucess for 1 sec expect length of messages %d got %d", 1, len(secTexts))
				wg.Done()
				return

			}
		}
	}()

	wg.Add(1)

	go func() {

		for {
			select {
			case <-secOut:
				secTexts = append(secTexts, <-secOut)

			}

			select {

			case <-minC.C:

				if len(secTexts) > 60 {
					t.Fatalf("lenth of received messages is higher  expected %d got %v ", 60, len(secTexts))
				}
				t.Logf("sucess for 1 min expect length of messages %d got %d", 60, len(secTexts))
				wg.Done()
				return

			}
		}
	}()

	wg.Add(1)

	go func() {

		for {
			select {
			case <-secOut:
				secTexts = append(secTexts, <-secOut)

			}

			select {
			case <-hourC.C:

				if len(secTexts) > 3600 {
					t.Fatalf("lenth of received messages is higher  expected %d got %v ", 3600, len(secTexts))
				}
				t.Logf("sucess for 1 hour expect length of messages %d got %d", 3600, len(secTexts))
				wg.Done()
				return

			}
		}
	}()

	wg.Wait()

}

func TestStartInputListner(t *testing.T) {
	filename := "tf"
	go clock.StartInputListner(filename, &clock.Field{})
	_, err := os.Open(filename)

	if err != nil {
		t.Fatalf("expeted file with name %s got %v", filename, err)
	}

}

func TestField_UpdateInput(t *testing.T) {
	nf := clock.Field{SecondText: "tick", MinuteText: "tock", HourText: "bong"}
	for _, v := range []struct {
		InputText string
		OldText   string
	}{
		{InputText: "--sec pick", OldText: "tick"},
		{InputText: "--min pock", OldText: "tock"},
		{InputText: "--hour tong", OldText: "bong"},
	} {
		nf.UpdateInput(v.InputText)

		switch v.InputText {

		case "--sec pick":
			if v.OldText == nf.SecondText {
				t.Fatalf("old and new value matches %s %s", v.OldText, nf.SecondText)
			}

		case "--min pock":
			if v.OldText == nf.MinuteText {
				t.Fatalf("old and new value matches %s %s", v.OldText, nf.SecondText)
			}

		case "--hour tong":

			if v.OldText == nf.HourText {
				t.Fatalf("old and new value matches %s %s", v.OldText, nf.SecondText)
			}

		default:
			t.Fatalf("input doesnt match")

		}

	}

}
