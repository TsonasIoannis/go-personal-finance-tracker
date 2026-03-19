package controllers

// messageResponse is referenced by Swagger annotations.
//
//nolint:unused
type messageResponse struct {
	Message string `json:"message"`
}

// authResponse is referenced by Swagger annotations.
//
//nolint:unused
type authResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    userResponse `json:"user"`
}

// transactionPageResponse is referenced by Swagger annotations.
//
//nolint:unused
type transactionPageResponse struct {
	Data       []transactionResponse `json:"data"`
	Pagination paginationResponse    `json:"pagination"`
}

// budgetPageResponse is referenced by Swagger annotations.
//
//nolint:unused
type budgetPageResponse struct {
	Data       []budgetResponse   `json:"data"`
	Pagination paginationResponse `json:"pagination"`
}
