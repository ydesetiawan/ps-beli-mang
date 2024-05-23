package helper

func ValidateAllowedSortBy(sortBy string, allowedSortBy []string) bool {
	for _, v := range allowedSortBy {
		if v == sortBy {
			return true
		}
	}
	return false
}

func AddKeywordSuffix(sortBy string, needKeywordSuffix []string) string {
	for _, v := range needKeywordSuffix {
		if v == sortBy {
			sortBy += ".keyword"
		}
	}
	return sortBy
}
