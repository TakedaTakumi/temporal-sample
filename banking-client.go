package app

// このコードは仮想の銀行サービスのクライアントをシミュレートします。
// 出金と入金の両方をサポートし、各リクエストごとに
// 疑似ランダムな取引IDを生成します。
//
// ヒント: これらの関数を修正して遅延やエラーを導入することで、
// 障害やタイムアウトの実験ができます。
import (
	"errors"
	"math/rand"
)

type account struct {
	AccountNumber string
	Balance       int64
}

type bank struct {
	Accounts []account
}

func (b bank) findAccount(accountNumber string) (account, error) {

	for _, v := range b.Accounts {
		if v.AccountNumber == accountNumber {
			return v, nil
		}
	}

	return account{}, errors.New("account not found")
}

// 残高不足時に発生するエラー
// InsufficientFundsError is raised when the account doesn't have enough money.
type InsufficientFundsError struct{}

func (m *InsufficientFundsError) Error() string {
	return "Insufficient Funds"
}

// アカウント番号が無効な場合に発生するエラー
// InvalidAccountError is raised when the account number is invalid
type InvalidAccountError struct{}

func (m *InvalidAccountError) Error() string {
	return "Account number supplied is invalid"
}

// モックバンク
// our mock bank
var mockBank = &bank{
	Accounts: []account{
		{AccountNumber: "85-150", Balance: 2000},
		{AccountNumber: "43-812", Balance: 0},
	},
}

// BankingServiceは銀行APIとのやりとりを模倣します。出金と入金をサポートします。
// BankingService mocks interaction with a bank API. It supports withdrawals
// and deposits
type BankingService struct {
	// ホスト名はより現実的にするためのものです。このコードは実際にはネットワーク通信を行いません。
	// the hostname is to make it more realistic. This code does not
	// actually make any network calls.
	Hostname string
}

// Withdrawは銀行からの出金をシミュレートします。
// アカウント番号（string）、金額（int）、リファレンスID（string）を受け取ります。
// リファレンスIDは冪等性のためのトランザクショントラッキング用です。
// 成功時はトランザクションIDを返します。
// 金額やアカウント番号に応じて様々なエラーを返します。
func (client BankingService) Withdraw(accountNumber string, amount int, referenceID string) (string, error) {
	acct, err := mockBank.findAccount(accountNumber)

	if err != nil {
		return "", &InvalidAccountError{}
	}

	if amount > int(acct.Balance) {
		return "", &InsufficientFundsError{}
	}

	return generateTransactionID("W", 10), nil
}

// Depositは銀行への入金をシミュレートします。
// アカウント番号（string）、金額（int）、リファレンスID（string）を受け取ります。
// リファレンスIDは冪等性のためのトランザクショントラッキング用です。
// 成功時はトランザクションIDを返します。
// アカウントが無効な場合はInvalidAccountErrorを返します。
func (client BankingService) Deposit(accountNumber string, amount int, referenceID string) (string, error) {

	_, err := mockBank.findAccount(accountNumber)
	if err != nil {
		return "", &InvalidAccountError{}
	}

	return generateTransactionID("D", 10), nil
}

// DepositThatFailsは未知のエラーをシミュレートします。
func (client BankingService) DepositThatFails(accountNumber string, amount int, referenceID string) (string, error) {
	return "", errors.New("This deposit has failed.")
}

// トランザクションIDを生成します
func generateTransactionID(prefix string, length int) string {
	randChars := make([]byte, length)
	for i := range randChars {
		allowedChars := "0123456789"
		randChars[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return prefix + string(randChars)
}
