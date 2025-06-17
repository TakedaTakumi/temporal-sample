package app

type PaymentDetails struct {
    SourceAccount string
    TargetAccount string
    Amount        int
    ReferenceID   string
}

const MoneyTransferTaskQueueName = "TRANSFER_MONEY_TASK_QUEUE"
