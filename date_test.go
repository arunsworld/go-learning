package learning

import (
	"testing"
	"time"
)

func TestDateParsing(t *testing.T) {
	dateStr := "Wed Jan 9 13:17:22 IST 2019"

	tt, err := time.Parse("Mon Jan 2 15:04:05 MST 2006", dateStr)
	if err != nil {
		t.Error("could not parse time: " + err.Error())
		return
	}

	if tt.Day() != 9 {
		t.Errorf("Expected Day to be 9, got: %d", tt.Day())
	}

	if tt.Month() != time.January {
		t.Errorf("Expected Day to be January, got: %s", tt.Month().String())
	}

	if tt.Year() != 2019 {
		t.Errorf("Expected Year to be 2019, got: %d", tt.Year())
	}

	if tt.Hour() != 13 {
		t.Errorf("Expected Hour to be 1, got: %d", tt.Hour())
	}

	if tt.Weekday() != time.Wednesday {
		t.Errorf("Expected Day to be Wednesday, got: %s", tt.Weekday().String())
	}

	tz, offset := tt.Zone()
	if tz != "IST" {
		t.Errorf("Expected timezone to be IST, got: %s", tz)
	}
	if offset != (5.5)*3600 {
		t.Errorf("Expected offset to be 5.5*3600 seconds, got: %d", offset)
	}

}

// Get timezones: zipinfo /usr/local/go/lib/time/zoneinfo.zip
func TestParseDateInOneLocationAndConvertToAnother(t *testing.T) {
	dateStr := "Wed Jan 9 13:17:22 2019"

	sourceLoc := mustGetLocation(t, "US/Pacific")
	targetLoc := mustGetLocation(t, "Asia/Kolkata")

	tt, err := time.ParseInLocation("Mon Jan 2 15:04:05 2006", dateStr, sourceLoc)
	if err != nil {
		t.Fatal("could not parse time: " + err.Error())
	}

	transformedTT := tt.In(targetLoc)
	transformedTTStr := transformedTT.Format("2006-01-02 15:04:05 -0700 MST")

	expected := "2019-01-10 02:47:22 +0530 IST"
	if transformedTTStr != expected {
		t.Fatalf("Expecting %s, got %s.", expected, transformedTTStr)
	}

}

type fatalistic interface {
	Fatal(...interface{})
}

func mustGetLocation(t fatalistic, name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatal(err)
	}
	return loc
}

func TestStrangeDateDueToDaylightSavings(t *testing.T) {
	// This is an invalid time. Because at 2AM on Mar 10 the clocks jump forward to 3AM.
	dateStr := "Sun Mar 10 02:00:00 2019"

	loc := mustGetLocation(t, "US/Pacific")
	tt, err := time.ParseInLocation("Mon Jan 2 15:04:05 2006", dateStr, loc)
	if err != nil {
		t.Fatal("could not parse time: " + err.Error())
	}

	transformedTT := tt.In(loc)
	transformedTTStr := transformedTT.Format("2006-01-02 15:04:05 -0700 MST")
	// Since this is an invalid date time is rolled back to the nearest correct hour (1 AM)
	// Not sure I'm happy with this behavior (I expect 3AM PDT) but I will reserve opinion and let the test pass
	expected := "2019-03-10 01:00:00 -0800 PST"
	if transformedTTStr != expected {
		t.Fatalf("Expecting %s, got %s.", expected, transformedTTStr)
	}

	transformedTTStr = tt.UTC().Format("2006-01-02 15:04:05 -0700 MST")
	expected = "2019-03-10 09:00:00 +0000 UTC"
	if transformedTTStr != expected {
		t.Fatalf("Expecting %s, got %s.", expected, transformedTTStr)
	}

	afterOneHourTT := tt.Add(time.Hour)
	afterOneHourTTStr := afterOneHourTT.Format("2006-01-02 15:04:05 -0700 MST")
	// After an hour from 1AM we get 3AM... PDT
	expected = "2019-03-10 03:00:00 -0700 PDT"
	if afterOneHourTTStr != expected {
		t.Fatalf("Expecting %s, got %s.", expected, afterOneHourTTStr)
	}
}

func TestMarshalling(t *testing.T) {
	loc := mustGetLocation(t, "US/Pacific")
	d, err := time.ParseInLocation("2006-01-02 15:04", "2019-03-10 03:00", loc)
	if err != nil {
		t.Fatal(err)
	}

	m, err := d.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	expected := "2019-03-10T03:00:00-07:00"
	if string(m) != expected {
		t.Fatalf("Expecting %s, got %s.", expected, string(m))
	}
}

func BenchmarkParsingInLocation(b *testing.B) {
	loc := mustGetLocation(b, "US/Pacific")
	for n := 0; n < b.N; n++ {
		d, err := time.ParseInLocation("2006-01-02 15:04", "2019-03-10 03:00", loc)
		if err != nil {
			b.Fatal(err)
		}
		_ = d
	}
}

func BenchmarkParsing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		d, err := time.Parse("2006-01-02 15:04", "2019-03-10 03:00")
		if err != nil {
			b.Fatal(err)
		}
		_ = d
	}
}

func BenchmarkFormatting(b *testing.B) {
	d, _ := time.Parse("2006-01-02 15:04", "2019-03-10 03:00")
	for n := 0; n < b.N; n++ {
		r := d.Format("2006-01-02 15:04")
		_ = r
	}
}

func BenchmarkMarshallingBinary(b *testing.B) {
	d, _ := time.Parse("2006-01-02 15:04", "2019-03-10 03:00")
	for n := 0; n < b.N; n++ {
		r, err := d.MarshalBinary()
		if err != nil {
			b.Fatal(err)
		}
		_ = r
	}
}

func BenchmarkMarshallingJSON(b *testing.B) {
	d, _ := time.Parse("2006-01-02 15:04", "2019-03-10 03:00")
	for n := 0; n < b.N; n++ {
		r, err := d.MarshalText()
		if err != nil {
			b.Fatal(err)
		}
		_ = r
	}
}
