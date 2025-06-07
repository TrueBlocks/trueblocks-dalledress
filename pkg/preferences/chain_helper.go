package preferences

// GetPreferredChainName returns the name of the first configured chain from user preferences
// or "mainnet" as a default if no chains are configured.
func GetPreferredChainName() string {
	userPrefs, err := GetUserPreferences()
	if err != nil || len(userPrefs.Chains) == 0 {
		return "mainnet" // Default to mainnet if no chains or error
	}

	// Return the chain name from the first configured chain
	return userPrefs.Chains[0].Chain
}
