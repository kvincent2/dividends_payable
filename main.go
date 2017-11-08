package main

// go run main.go [access key]

import (
  "fmt"
  quickbooks "github.com/jinmatt/go-quickbooks.v2"
  "io/ioutil"
  "os"
  "strconv"
  "time"
)

func main() {
  //declare date variable for use later in the function; get current date
  year, month, _ := time.Now().Date()

  //Read last month's ending Series R investment balance from individual txt files.
  //Panic if file doesn't exist or doesn't contain a number value.
  emily_total_div := getLastMonthsTotal("emilyP", int(month-1), checkIfLastYear(int(month), year))
  luke_total_div := getLastMonthsTotal("lukeJ", int(month-1), checkIfLastYear(int(month), year))
  jordan_total_div := getLastMonthsTotal("jordanS", int(month-1), checkIfLastYear(int(month), year))
  ray_total_div := getLastMonthsTotal("rayS", int(month-1), checkIfLastYear(int(month), year))

  // Calculate current month dividend expense and store value for use in JE json object.
  dividendPercent := .10
  emily_current_month := ((emily_total_div * dividendPercent) / float64(365) * float64(daysIn(month, year)))
  luke_current_month := ((luke_total_div * dividendPercent) / float64(365) * float64(daysIn(month, year)))
  jordan_current_month := ((jordan_total_div * dividendPercent) / float64(365) * float64(daysIn(month, year)))
  ray_current_month := ((ray_total_div * dividendPercent) / float64(365) * float64(daysIn(month, year)))

  total_current_exp := emily_current_month + luke_current_month + jordan_current_month + ray_current_month

  //Prepare JE lines from stored values.
  journalEntryLines := []quickbooks.Line{}

  emily_JE_Line := createJournalEntryLine("0", fmt.Sprintf("To record Emily P's %s dividends payable.", month), emily_current_month, "Debit", "41", "Dividend Expense")
  luke_JE_Line := createJournalEntryLine("1", fmt.Sprintf("To record Luke J's %s dividends payable.", month), luke_current_month, "Debit", "41", "Dividend Expense")
  jordan_JE_Line := createJournalEntryLine("2", fmt.Sprintf("To record Jordan S's %s dividends payable.", month), jordan_current_month, "Debit", "41", "Dividend Expense")
  ray_JE_Line := createJournalEntryLine("3", fmt.Sprintf("To record Ray S's %s dividends payable.", month), ray_current_month, "Debit", "41", "Dividend Expense")

  journalEntryLines = append(journalEntryLines, emily_JE_Line)
  journalEntryLines = append(journalEntryLines, luke_JE_Line)
  journalEntryLines = append(journalEntryLines, jordan_JE_Line)
  journalEntryLines = append(journalEntryLines, ray_JE_Line)

  journalEntryLines = append(journalEntryLines, createJournalEntryLine("4", fmt.Sprintf("To record %s dividends payable.", month), total_current_exp, "Credit", "213", "Dividends Payable"))

  //create quickbooks client
  //expects QBO_realmID_production environment variable to be set
  //expects to be passed access key from quickbooks developer dashboard
  quickbooksClient := quickbooks.NewClient(os.Getenv("QBO_realmID_production"), os.Args[1], false)

  journalEntry := quickbooks.Journalentry{
    TxnDate: fmt.Sprintf("%d-%d-%d", year, int(month), daysIn(month, year)),
    Line:    journalEntryLines,
  }

  JournalentryObject, err := quickbooksClient.CreateJE(journalEntry)
  fmt.Println(JournalentryObject, err)

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

func createJournalEntryLine(lineID string, description string, amount float64, postingType string, accountNum string, accountName string) quickbooks.Line {
  return quickbooks.Line{
    LineID:      lineID,
    Description: description,
    Amount:      amount,
    DetailType:  "JournalEntryLineDetail",
    JournalEntryLineDetail: &quickbooks.JournalEntryLineDetail{
      PostingType: postingType,
      AccountRef: quickbooks.JournalEntryRef{
        Value: accountNum,
        Name:  accountName,
      },
    },
  }
}

func getLastMonthsTotal(investor string, month int, year int) float64 {
  fileContents, err := ioutil.ReadFile(fmt.Sprintf("./investorDivPayableYTD/%s_%d_%d.txt", investor, month, year))
  if err != nil {
    panic(err)
  }
  total, err := strconv.ParseFloat(string(fileContents), 64)
  if err != nil {
    panic(err)
  }
  return total
}

func checkIfLastYear(month int, year int) int {
  if month == 1 {
    return year - 1
  } else {
    return year
  }
}
