package dto

type CreatePlanRequest struct {
	Goal string `json:"goal" binding:"required"`
}

type PlanResponse struct {
	ID    string      `json:"id"`
	Goal  string      `json:"goal"`
	Tasks interface{} `json:"tasks"`
}
