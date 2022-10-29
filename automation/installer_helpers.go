package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	msiDll                  = syscall.NewLazyDLL("msi.dll")
	msiDllOpenDatabaseW     = msiDll.NewProc("MsiOpenDatabaseW")
	msiDllDatabaseOpenViewW = msiDll.NewProc("MsiDatabaseOpenViewW")
	msiDllViewExecute       = msiDll.NewProc("MsiViewExecute")
	msiDllViewFetch         = msiDll.NewProc("MsiViewFetch")
	msiRecordGetString      = msiDll.NewProc("MsiRecordGetStringW")
	msiDllViewClose         = msiDll.NewProc("MsiViewClose")
	msiDllCloseHandle       = msiDll.NewProc("MsiCloseHandle")
)

func getMsixPackageFamilyName(msixPath string) (string, error) {
	encodingTable := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	err := extractFileFromZip(msixPath, "AppxManifest.xml")
	if err != nil {
		return "", fmt.Errorf("error extracting AppxManifest.xml from msix: %v\n", err)
	}
	appxManifest, _ := os.ReadFile("AppxManifest.xml")
	identityName := regexp.MustCompile(`(?m)<Identity.*?Name="(.+?)"`).FindStringSubmatch(string(appxManifest))[1]
	identityPublisher := regexp.MustCompile(`(?m)<Identity.*?Publisher="(.+?)"`).FindStringSubmatch(string(appxManifest))[1]
	utf16Bytes := utf16.Encode([]rune(identityPublisher))
	var binaryBuffer = new(bytes.Buffer)
	binary.Write(binaryBuffer, binary.LittleEndian, utf16Bytes)
	publisherUnicodeSha256 := sha256.Sum256(binaryBuffer.Bytes())
	var sha256HashBytesInBinary, result string
	for _, char := range publisherUnicodeSha256[:8] {
		sha256HashBytesInBinary += fmt.Sprintf("%08b", char)
	}
	sha256HashBytesInBinary += strings.Repeat("0", 65-len(sha256HashBytesInBinary))
	for i := 0; i < len(sha256HashBytesInBinary); i += 5 {
		index, _ := strconv.ParseInt(sha256HashBytesInBinary[i : i+5], 2, 64)
		result += string(encodingTable[index])
	}
	return identityName + "_" + result, nil
}

func getMsixSignatureHash(msixPath string) (string, error) {
	err := extractFileFromZip(msixPath, "AppxSignature.p7x")
	if err != nil {
		return "", fmt.Errorf("error extracting AppxSignature.p7x from msix: %v\n", err)
	}
	sha256Hash, err := getFileSha256Hash("AppxSignature.p7x")
	if err != nil {
		return "", fmt.Errorf("error getting sha256 hash of AppxSignature.p7x: %v\n", err)
	}
	return sha256Hash, nil
}

func getVersionFromInstaller(installerUrl string) string {
	// build filename from url
	parsedUrl, _ := url.Parse(installerUrl)
	urlSegments := strings.Split(parsedUrl.Path, "/")
	fileName := urlSegments[len(urlSegments)-1]
	fileExtension := strings.Split(fileName, ".")[1]

	workDir, _ := os.Getwd()
	tempFile, _ := os.CreateTemp(workDir, "installer-*."+fileExtension)
	defer os.Remove(tempFile.Name())

	// download the installer using http client
	resp, err := http.Get(installerUrl)
	if err != nil {
		log.Fatalf("error downloading installer: %v\n", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		log.Fatalf("error writing installer to temp file: %v\n", err)
	}

	if fileExtension == "msi" {
		return getMsiDbProperty("ProductVersion", tempFile.Name())
	} else {
		cmdToRun := fmt.Sprintf("(Get-Item -Path '%s').VersionInfo.ProductVersion", tempFile.Name())
		cmdStdout, _ := exec.Command("pwsh", "-Command", cmdToRun).Output()
		return strings.TrimSpace(string(cmdStdout))
	}
}

func getMsiDbProperty(property, msiPath string) string {
	// open the msi database
	var msiHandle syscall.Handle
	msiPathPtr, _ := syscall.UTF16PtrFromString(msiPath)
	ret, _, _ := msiDllOpenDatabaseW.Call(uintptr(unsafe.Pointer(msiPathPtr)), uintptr(0), uintptr(unsafe.Pointer(&msiHandle)))
	if ret != 0 {
		log.Fatalf("error opening msi database: %d\n", ret)
	}

	// create a query to get the property
	var query string = "SELECT Value FROM Property WHERE Property = '" + property + "'"
	var viewHandle syscall.Handle
	ret, _, _ = msiDllDatabaseOpenViewW.Call(uintptr(msiHandle), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(query))), uintptr(unsafe.Pointer(&viewHandle)))
	if ret != 0 {
		log.Fatalf("error opening view: %d\n", ret)
	}

	// execute the query
	ret, _, _ = msiDllViewExecute.Call(uintptr(viewHandle), uintptr(0))
	if ret != 0 {
		log.Fatalf("error executing view: %d\n", ret)
	}

	// fetch the result
	var recordHandle syscall.Handle
	ret, _, _ = msiDllViewFetch.Call(uintptr(viewHandle), uintptr(unsafe.Pointer(&recordHandle)))
	if ret != 0 {
		log.Fatalf("error fetching view: %d\n", ret)
	}

	// get the size of the value
	var valueSize uint32 = 0
	ret, _, _ = msiRecordGetString.Call(uintptr(recordHandle), uintptr(1), uintptr(0), uintptr(unsafe.Pointer(&valueSize)))
	if ret != 0 {
		log.Fatalf("error getting string size: %d\n", ret)
	}

	// allocate a buffer for the value
	valueSize++ // add 1 for null terminator
	valueBuffer := make([]uint16, valueSize)

	// get the value
	ret, _, _ = msiRecordGetString.Call(uintptr(recordHandle), uintptr(1), uintptr(unsafe.Pointer(&valueBuffer[0])), uintptr(unsafe.Pointer(&valueSize)))
	if ret != 0 {
		log.Fatalf("error getting string data: %d\n", ret)
	}

	// close the view
	ret, _, _ = msiDllViewClose.Call(uintptr(viewHandle))
	if ret != 0 {
		log.Fatalf("error closing view: %d\n", ret)
	}

	// close the msi database
	ret, _, _ = msiDllCloseHandle.Call(uintptr(msiHandle))
	if ret != 0 {
		log.Fatalf("error closing msi database: %d\n", ret)
	}

	// return the value
	return syscall.UTF16ToString(valueBuffer)
}
