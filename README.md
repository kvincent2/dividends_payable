The purpose of this app is to calculate my company's dividend expense for the month, develop a journal entry, post it to Quickbooks Online using their API, and keeps track of the running balance of dividends payable.

Requirements:
In order to run this app, you need a few things:
1. Go installation
2. A developer.intuit.com account
3. An app on developer.intuit.com and the associated access key
4. This app reads .txt files with the investor's YTD dividends payable amount. These must be created and saved in the investorDivPayableYTD folder of the app. This file only contains the balance with two decimal places. The script casts this to a float64. The file names are indicative of the name of the investor and date.

First Use Instructions:
1. Clone the GitHub repo to your computer and place it in the src directory of your $GOPATH
2. Set your Quickbooks Online realmID as an environment variable "QBO_realmID_production"
3. Use the OAuth2 Playground at developer.intuit.com to request an access key 
4. When you first set this up, create .txt files for each investor for the previous month.
 - userL_month_year.txt

Running the code:
1. Access key must be entered as a parameter to the main function.
2. Use the command `go run main.go [access key]` to run script.
3. If the required .txt files are not available, or if the .txt files contain non-digit values, the script will panic and end.

