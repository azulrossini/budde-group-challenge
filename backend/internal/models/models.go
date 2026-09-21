package models

type Bill struct {
	ID          int64
	Description string
	TotalCents  int64
}

type Share struct {
	ID                    int64
	BillID                int64
	Name                  string
	PercentageBasisPoints int32
}

type BillShares struct {
	Bill   Bill
	Shares []Share
}
