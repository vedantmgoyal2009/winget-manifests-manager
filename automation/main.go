package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	ctx                                                        = context.Background()
	firestoreClient, _                                         = getFirestoreClient(ctx)
	githubClient                                               = getGithubClient(os.Getenv("GH_TOKEN"), ctx)
	githubAppToken, _                                          = getGithubAppAuthToken()
	githubAppClient                                            = getGithubClient(githubAppToken, ctx)
	workingDir, _                                              = os.Getwd()
	wingetDevPreviousBuildInfo                                 WinGetDevBuildInfo
	failedCheckUpdates, failedManifestUpdates, skippedPackages []string
)

func main() {
	// clone winget-pkgs, copy YamlCreate.ps1, and update git configuration
	_ = runCmd("git config --global user.name 'vedantmgoyal2009[bot]'")
	_ = runCmd("git config --global user.email '110876359+vedantmgoyal2009[bot]@users.noreply.github.com'")
	_ = runCmd("git clone https://github.com/microsoft/winget-pkgs.git --quiet")
	_ = runCmd("git -C winget-pkgs remote rename origin upstream")
	_ = runCmd("git -C winget-pkgs remote add origin https://x-access-token:" + githubAppToken + "@github.com/vedantmgoyal2009/winget-pkgs.git")
	_ = runCmd("git -C winget-pkgs fetch origin --quiet")
	_ = runCmd("git -C winget-pkgs config core.safecrlf false")
	_ = runCmd("Copy-Item -Path .\\YamlCreate.ps1 -Destination .\\winget-pkgs\\Tools\\YamlCreate.ps1 -Force")
	_ = runCmd("git -C winget-pkgs commit --all -m 'Update YamlCreate.ps1 with InputObject functionality'")
	fmt.Println("Cloned winget-pkgs, copied YamlCreate.ps1, and updated git configuration.")

	// block Microsoft Edge updates, install powershell-yaml, import functions, copy YamlCreate.ps1 to the Tools folder, and update git configuration
	// to prevent edge from updating and changing ARP table during ARP metadata validation
	_ = runCmd("New-Item -Path HKLM:\\SOFTWARE\\Microsoft\\EdgeUpdate -Force")
	_ = runCmd("New-ItemProperty -Path HKLM:\\SOFTWARE\\Microsoft\\EdgeUpdate -Name DoNotUpdateToEdgeWithChromium -Value 1 -PropertyType DWord -Force")
	_ = runCmd("Set-Service -Name edgeupdate -Status Stopped -StartupType Disabled")
	_ = runCmd("Set-Service -Name edgeupdatem -Status Stopped -StartupType Disabled")
	_ = runCmd("Install-Module -Name powershell-yaml -Repository PSGallery -Scope CurrentUser -Force")
	fmt.Println("Successfully installed powershell-yaml")

	wingetCliCommitInfo, _, _ := githubAppClient.Repositories.ListCommits(ctx, "microsoft", "winget-cli", nil)
	wingetDevLastBuild, _ := os.ReadFile("wingetdev/build.json")
	json.NewDecoder(io.NopCloser(bytes.NewReader(wingetDevLastBuild))).Decode(&wingetDevPreviousBuildInfo)
	if wingetDevPreviousBuildInfo.Commit.Sha != wingetCliCommitInfo[0].GetSHA() {
		fmt.Println("New commit pushed on microsoft/winget-cli, updating wingetdev...")
		fmt.Println("This will take about ~15 minutes... please wait...")
		_ = runCmd("git clone https://github.com/microsoft/winget-cli.git --quiet")
		_ = runCmd("& 'C:\\Program Files\\Microsoft Visual Studio\\2022\\Enterprise\\VC\\Auxiliary\\Build\vcvarsall.bat' x64")
		_ = runCmd("& 'C:\\Program Files\\Microsoft Visual Studio\\2022\\Enterprise\\MSBuild\\Current\\Bin\\MSBuild.exe' -t:restore -m -p:RestorePackagesConfig=true -p:Configuration=Release -p:Platform=x64 .\\winget-cli\\src\\AppInstallerCLI.sln | Out-File -FilePath ..\\tools\\wingetdev\\log.txt -Append")
		_ = runCmd("& 'C:\\Program Files\\Microsoft Visual Studio\\2022\\Enterprise\\MSBuild\\Current\\Bin\\MSBuild.exe' -m -p:Configuration=Release -p:Platform=x64 .\\winget-cli\\src\\AppInstallerCLI.sln | Out-File -FilePath ..\\tools\\wingetdev\\log.txt -Append")
		_ = runCmd("Copy-Item -Path .\\winget-cli\\src\\x64\\Release\\WindowsPackageManager\\WindowsPackageManager.dll -Destination ..\\tools\\wingetdev\\WindowsPackageManager.dll -Force")
		_ = runCmd("Move-Item -Path .\\winget-cli\\src\\x64\\Release\\AppInstallerCLI\\* -Destination ..\\tools\\wingetdev -Force")
		_ = runCmd("Move-Item -Path ..\tools\\wingetdev\\winget.exe -Destination wingetdev.exe -Force")
		wingetDevBuildJson, _ := os.Create("wingetdev/build.json")
		json.NewEncoder(wingetDevBuildJson).Encode(WinGetDevBuildInfo{
			Commit: WinGetDevBuildInfoCommit{
				Sha:     wingetCliCommitInfo[0].GetSHA(),
				Message: wingetCliCommitInfo[0].Commit.GetMessage(),
				Author:  *wingetCliCommitInfo[0].GetAuthor().Login,
			},
			BuildDateTime: time.Now().Format(time.RFC1123),
		})
	}
	os.Setenv("WINGETDEV", workingDir+"\\wingetdev\\wingetdev.exe")
	// enable installation of local manifests by wingetdev, disabled by default for security purposes
	// see https://github.com/microsoft/winget-cli/pull/1453 for more info
	runCmd("& $env:WINGETDEV settings --enable LocalManifestFiles")

	for _, pkg := range getPackages(firestoreClient, ctx) {
		if pkg.SkipPackage != "false" {
			skippedPackages = append(skippedPackages, pkg.Identifier)
			continue
		}

		for key, value := range pkg.AdditionalInfo {
			os.Setenv(key, value)
		}

		parsedUrl, _ := url.Parse(pkg.Update.Uri)
		var headers http.Header
		if pkg.Update.Headers != nil {
			for key, value := range pkg.Update.Headers {
				headers.Add(key, value)
			}
		}
		httpRequest := &http.Request{
			URL:    parsedUrl,
			Method: pkg.Update.Method,
			Header: headers,
		}

		resp, err := http.DefaultClient.Do(httpRequest)
		if err != nil {
			error := fmt.Sprintf("%s: %s", pkg.Identifier, err)
			failedCheckUpdates = append(failedCheckUpdates, error)
			continue
		}
		defer resp.Body.Close()

		var parsedResp = make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&parsedResp)
		if err != nil {
			// if it's not a valid JSON, then convert the response body to a string
			respBytes, _ := io.ReadAll(resp.Body)
			parsedString := string(respBytes)
		}

		if isNewVersionGreater(updateInfo["Version"], pkg.PreviousVersion) {
			// if the new version is greater than the previous version, then update the package
		}

		for key, _ := range pkg.AdditionalInfo {
			os.Remove(key)
		}
	}
}
