package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

func compareBeforeAndAfterArpEntries(original, new []ArpEntry) []ArpEntry {
	var changed []ArpEntry

	for _, new_entry := range new {
		for _, original_entry := range original {
			if !(new_entry.DisplayName == original_entry.DisplayName &&
				new_entry.DisplayVersion == original_entry.DisplayVersion &&
				new_entry.Publisher == original_entry.Publisher &&
				new_entry.ProductCode == original_entry.ProductCode) {
				changed = append(changed, new_entry)
				break
			}
		}
	}

	return changed
}

func getAppsAndFeaturesEntries() ([]ArpEntry, error) {
	registry_paths := []string{
		"Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\",
		"Software\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\",
	}
	path_roots := []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER}
	var entries []ArpEntry

	for _, path_root := range path_roots {
		for _, registry_path := range registry_paths {
			key, err := registry.OpenKey(path_root, registry_path, registry.ENUMERATE_SUB_KEYS) // open the registry key
			if err != registry.ErrNotExist && err != nil {
				return nil, fmt.Errorf("error opening registry key: %v\n", err)
			}

			// get all subkeys and iterate through them
			subkeys, _ := key.ReadSubKeyNames(-1) // -1 means read all subkeys
			for _, subkey := range subkeys {
				product_code := subkey // product code is the subkey name
				subkey, err := registry.OpenKey(path_root, registry_path+subkey, registry.QUERY_VALUE)
				if err != nil {
					return nil, fmt.Errorf("error opening subkey: %v\n", err)
				}

				// get the values
				display_name, _, _ := subkey.GetStringValue("DisplayName")
				display_version, _, _ := subkey.GetStringValue("DisplayVersion")
				publisher, _, _ := subkey.GetStringValue("Publisher")
				system_component, _, _ := subkey.GetIntegerValue("SystemComponent")

				// add the entry to the list if display name is *not* empty, and it's *not* a system component
				if display_name != "" && system_component != 1 {
					entries = append(entries, ArpEntry{display_name, display_version, publisher, product_code})
				}
			}
		}
	}

	return entries, nil
}
