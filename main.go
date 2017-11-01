package main

// go run main.go [access key]

import (
  "fmt"
  quickbooks "github.com/jinmatt/go-quickbooks.v2"
  "io/ioutil"
  "math"
  "os"
  "strconv"
  "time"
)

func main() {
  //declare date variable for use later in the function; get current date
  year, month, _ := time.Now().Date()

  //Read last month's ending Series R investment balance from individual txt files.
  //Panic if file doesn't exist or doesn't contain a number value.
  emily_total_div_raw, err := ioutil.ReadFile(fmt.Sprintf("./investorDivPayableYTD/emilyP_%d_%d.txt", int(month-1), year))
  if err != nil {
    panic(err)
  }
  emily_total_div, err := strconv.ParseFloat(string(emily_total_div_raw), 64)
  if err != nil {
    panic(err)
  }

  luke_total_div_raw, err := ioutil.ReadFile(fmt.Sprintf("./investorDivPayableYTD/lukeJ_%d_%d.txt", int(month-1), year))
  if err != nil {
    panic(err)
  }
  luke_total_div, err := strconv.ParseFloat(string(luke_total_div_raw), 64)
  if err != nil {
    panic(err)
  }

  jordan_total_div_raw, err := ioutil.ReadFile(fmt.Sprintf("./investorDivPayableYTD/jordanS_%d_%d.txt", int(month-1), year))
  if err != nil {
    panic(err)
  }
  jordan_total_div, err := strconv.ParseFloat(string(jordan_total_div_raw), 64)
  if err != nil {
    panic(err)
  }

  ray_total_div_raw, err := ioutil.ReadFile(fmt.Sprintf("./investorDivPayableYTD/rayS_%d_%d.txt", int(month-1), year))
  if err != nil {
    panic(err)
  }
  ray_total_div, err := strconv.ParseFloat(string(ray_total_div_raw), 64)
  if err != nil {
    panic(err)
  }

  // Calculate current month dividend expense and store value for use in JE json object.
  dividendPercent := .10
  emily_current_month := ((emily_total_div * dividendPercent) / float64(365) * float64(daysIn(month, int(year))))
  luke_current_month := ((luke_total_div * dividendPercent) / float64(365) * float64(daysIn(month, int(year))))
  jordan_current_month := ((jordan_total_div * dividendPercent) / float64(365) * float64(daysIn(month, int(year))))
  ray_current_month := ((ray_total_div * dividendPercent) / float64(365) * float64(daysIn(month, int(year))))

  total_current_exp := emily_current_month + luke_current_month + jordan_current_month + ray_current_month

  //Prepare JE lines from stored values.
  journalEntryLines := []quickbooks.Line{}
  emily_JE_Line := quickbooks.Line{
    LineID:      "0",
    Description: fmt.Sprintf("To record Emily P's %s dividends payable.", month),
    Amount:      emily_current_month,
    DetailType:  "JournalEntryLineDetail",
    JournalEntryLineDetail: &quickbooks.JournalEntryLineDetail{
      PostingType: "Debit",
      AccountRef: quickbooks.JournalEntryRef{
        Value: "41",
        Name:  "Dividend Expense",
      },
    },
  }

  luke_JE_Line := quickbooks.Line{
    LineID:      "1",
    Description: fmt.Sprintf("To record Luke J's %s dividends payable.", month),
    Amount:      luke_current_month,
    DetailType:  "JournalEntryLineDetail",
    JournalEntryLineDetail: &quickbooks.JournalEntryLineDetail{
      PostingType: "Debit",
      AccountRef: quickbooks.JournalEntryRef{
        Value: "41",
        Name:  "Dividend Expense",
      },
    },
  }

  jordan_JE_Line := quickbooks.Line{
    LineID:      "2",
    Description: fmt.Sprintf("To record Jordan S's %s dividends payable.", month),
    Amount:      jordan_current_month,
    DetailType:  "JournalEntryLineDetail",
    JournalEntryLineDetail: &quickbooks.JournalEntryLineDetail{
      PostingType: "Debit",
      AccountRef: quickbooks.JournalEntryRef{
        Value: "41",
        Name:  "Dividend Expense",
      },
    },
  }

  ray_JE_Line := quickbooks.Line{
    LineID:      "3",
    Description: fmt.Sprintf("To record Ray S's %s dividends payable.", month),
    Amount:      ray_current_month,
    DetailType:  "JournalEntryLineDetail",
    JournalEntryLineDetail: &quickbooks.JournalEntryLineDetail{
      PostingType: "Debit",
      AccountRef: quickbooks.JournalEntryRef{
        Value: "41",
        Name:  "Dividend Expense",
      },
    },
  }

  journalEntryLines = append(journalEntryLines, emily_JE_Line)
  journalEntryLines = append(journalEntryLines, luke_JE_Line)
  journalEntryLines = append(journalEntryLines, jordan_JE_Line)
  journalEntryLines = append(journalEntryLines, ray_JE_Line)

  //Create credit and append to array
  totalExpenseLine := quickbooks.Line{
    LineID:      "4",
    Description: fmt.Sprintf("To record %s dividends payable.", month),
    Amount:      total_current_exp,
    DetailType:  "JournalEntryLineDetail",
    JournalEntryLineDetail: &quickbooks.JournalEntryLineDetail{
      PostingType: "Credit",
      AccountRef: quickbooks.JournalEntryRef{
        Value: "213",
        Name:  "Dividends Payable",
      },
    },
  }

  journalEntryLines = append(journalEntryLines, totalExpenseLine)

  //create quickbooks client
  //expects QBO_realmID_production environment variable to be set
  //expects to be passed access key from quickbooks developer dashboard
  quickbooksClient := quickbooks.NewClient(os.Getenv("QBO_realmID_production"), os.Args[3], false)

  journalEntry := quickbooks.Journalentry{
    Line: journalEntryLines,
  }

  JournalentryObject, err := quickbooksClient.CreateJE(journalEntry)

  //add current month's div expense to running total to get ending balance for the month.
  emily_ending_div := emily_total_div + emily_current_month
  luke_ending_div := luke_total_div + luke_current_month
  jordan_ending_div := jordan_total_div + jordan_current_month
  ray_ending_div := ray_total_div + ray_current_month

  //Write ending balance to txt file for use next month.
  ioutil.WriteFile(fmt.Sprintf("./investorDivPayableYTD/emilyP_%d_%d.txt", int(month), year), []byte(strconv.FormatFloat(emily_ending_div, 'f', 6, 64)), 0644)
  ioutil.WriteFile(fmt.Sprintf("./investorDivPayableYTD/lukeJ_%d_%d.txt", int(month), year), []byte(strconv.FormatFloat(luke_ending_div, 'f', 6, 64)), 0644)
  ioutil.WriteFile(fmt.Sprintf("./investorDivPayableYTD/jordanS_%d_%d.txt", int(month), year), []byte(strconv.FormatFloat(jordan_ending_div, 'f', 6, 64)), 0644)
  ioutil.WriteFile(fmt.Sprintf("./investorDivPayableYTD/rayS_%d_%d.txt", int(month), year), []byte(strconv.FormatFloat(ray_ending_div, 'f', 6, 64)), 0644)

}

//How many days in each month?
func daysIn(m time.Month, year int) int {
  // This is equivalent to time.daysIn(m, year).
  return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()

  return 0
}
