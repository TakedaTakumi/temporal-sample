package app

import (
    "context"
    "fmt"
    "log"
)

// 引き出し
func Withdraw(ctx context.Context, data PaymentDetails) (string, error) {
    log.Printf("Withdrawing $%d from account %s.\n\n",
        data.Amount,
        data.SourceAccount,
    )

    referenceID := fmt.Sprintf("%s-withdrawal", data.ReferenceID)
    bank := BankingService{"bank-api.example.com"}
    confirmation, err := bank.Withdraw(data.SourceAccount, data.Amount, referenceID)
    return confirmation, err
}

// 入金
func Deposit(ctx context.Context, data PaymentDetails) (string, error) {
    log.Printf("Depositing $%d into account %s.\n\n",
        data.Amount,
        data.TargetAccount,
    )

    referenceID := fmt.Sprintf("%s-deposit", data.ReferenceID)
    bank := BankingService{"bank-api.example.com"}
    // 未知の失敗をシミュレートするには、次の行のコメントを外し、その次の行をコメントアウトしてください
    // confirmation, err := bank.DepositThatFails(data.TargetAccount, data.Amount, referenceID)
    confirmation, err := bank.Deposit(data.TargetAccount, data.Amount, referenceID)
    return confirmation, err
}

// 払い戻し
func Refund(ctx context.Context, data PaymentDetails) (string, error) {
    log.Printf("Refunding $%v back into account %v.\n\n",
        data.Amount,
        data.SourceAccount,
    )

    referenceID := fmt.Sprintf("%s-refund", data.ReferenceID)
    bank := BankingService{"bank-api.example.com"}
    confirmation, err := bank.Deposit(data.SourceAccount, data.Amount, referenceID)
    return confirmation, err
}
