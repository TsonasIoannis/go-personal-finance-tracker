package controllers

type messageResponse struct {
	Message string `json:"message"`
}

type authResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    userResponse `json:"user"`
}

type transactionPageResponse struct {
	Data       []transactionResponse `json:"data"`
	Pagination paginationResponse    `json:"pagination"`
}

type budgetPageResponse struct {
	Data       []budgetResponse   `json:"data"`
	Pagination paginationResponse `json:"pagination"`
}
