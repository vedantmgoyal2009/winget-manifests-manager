package main

type AutomationPackage struct {
	Identifier string
	Update     struct {
		InvokeType string
		Uri        string
		Method     string
		Headers    map[string]string
		Body       string
		UserAgent  string
	}
	PostResponseScript []string
	VersionRegex       string
	PreviousVersion    string
	ManifestFields     map[string]interface{}
	PostUpgradeScript  []string
	AdditionalInfo     map[string]string
	SkipPackage        string
}

type ArpEntry struct {
	DisplayName    string
	DisplayVersion string
	Publisher      string
	ProductCode    string
}

type WinGetDevBuildInfo struct {
	Commit        _WinGetDevBuildInfoCommit `json:"Commit"`
	BuildDateTime string                   `json:"BuildDateTime"`
}

type _WinGetDevBuildInfoCommit struct {
	Sha     string `json:"Sha"`
	Message string `json:"Message"`
	Author  string `json:"Author"`
}
