package preferences

//-------------------------------------------------------------------------------------------
func GetLastChain() string {
	if appPreferences, err := GetAppPreferences(); err == nil {
		if len(appPreferences.LastChain) > 0 {
			return appPreferences.LastChain
		}
	}
	userPrefs, err := GetUserPreferences()
	if err != nil || len(userPrefs.Chains) == 0 {
		return "mainnet"
	}
	return userPrefs.Chains[0].Chain
}

//-------------------------------------------------------------------------------------------
func GetLastAddress() string {
	if appPreferences, err := GetAppPreferences(); err == nil {
		if len(appPreferences.LastAddress) > 0 {
			return appPreferences.LastAddress
		}
	}
	return "0xf503017d7baf7fbc0fff7492b751025c6a78179b"
}

//-------------------------------------------------------------------------------------------
func GetLastContract() string {
	if appPreferences, err := GetAppPreferences(); err == nil {
		if len(appPreferences.LastContract) > 0 {
			return appPreferences.LastContract
		}
	}
	return "0x52df6e4d9989e7cf4739d687c765e75323a1b14c"
}

//-------------------------------------------------------------------------------------------
// GetLastProject returns the active project from app preferences with caching

