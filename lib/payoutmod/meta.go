package payoutmod

const (
	// From upper levels 1 --> 7
	DailyProfitAffiliateRate = (33 + 16 + 16 + 10 + 10 + 6 + 6) / 100.0

	DailyProfitLeaderRewardTotalPercentage = 0 +
		16 + // Senior Partner
		8 + // Regional Partner
		8 + // Global Partner
		16 // Governer
	DailyProfitLeaderRewardRate = DailyProfitLeaderRewardTotalPercentage / 100.0
)
