package api

// RequestBody структура для данных запроса транзакции
type RequestBody struct {
	WalletId      string `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount        int64  `json:"amount"`
}

// ResponseBody структура для ответа с балансом
type ResponseBody struct {
	WalletID string `json:"walletId"`
	Balance  int64  `json:"balance"`
}
