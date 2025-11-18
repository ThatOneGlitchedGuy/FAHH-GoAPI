package schemas

type AdminStatsOut struct {
	TotalUsers   int64 `json:"total_users"`
	TotalProducts int64 `json:"total_products"`
	TotalOrders  int64 `json:"total_orders"`
}
