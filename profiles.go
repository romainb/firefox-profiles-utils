package profiles

import (
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "runtime"
    "strings"

    "github.com/shirou/gopsutil/v3/process"
)

// FirefoxProfiles interfaces Firefox profiles related operations.
type FirefoxProfiles interface {
    IsProfileUsed(profileName string) (bool, error)
    GetProfilesList() ([]string, error)
    GetProfilesPath() string
}

// FirefoxProfilesDefaultImpl implements FirefoxProfiles
type FirefoxProfilesDefaultImpl struct {
    // ProfilesPath Firefox profiles directory path.
    ProfilesPath string
}

// NewWithDefaultPath builds a FirefoxProfiles instance with the default path according to the OS.
func NewWithDefaultPath() FirefoxProfiles {
    userHome, _ := os.UserHomeDir()
    switch currentOs := runtime.GOOS; currentOs {
    case "windows":
        return FirefoxProfilesDefaultImpl{filepath.Join(userHome, `AppData\Roaming\Mozilla\Firefox\Profiles`)}
    case "darwin":
        return FirefoxProfilesDefaultImpl{filepath.Join(userHome, "Library/Application Support/Firefox/Profiles")}
    case "linux":
        return FirefoxProfilesDefaultImpl{filepath.Join(userHome, ".mozilla/firefox")}
    default:
        return FirefoxProfilesDefaultImpl{}
    }
}

// NewWithCustomPath builds a FirefoxProfiles instance with a specific profiles path.
func NewWithCustomPath(profilesPath string) FirefoxProfiles {
    return &FirefoxProfilesDefaultImpl{ProfilesPath: profilesPath}
}

// IsProfileUsed returns true if the specified profile is currently used.
func (ffxp FirefoxProfilesDefaultImpl) IsProfileUsed(profileName string) (bool, error) {
    processes, err := process.Processes()
    if err != nil {
        return false, errors.New("impossible d'obtenir la liste des processus en cours")
    }
    for _, p := range processes {
        name, _ := p.Name()
        args, _ := p.Cmdline()
        if name == "firefox" && strings.Contains(args, fmt.Sprintf("-P %s", profileName)) {
            return true, nil
        }
    }
    return false, nil
}

// GetProfilesList returns the list of existing Firefox profiles.
func (ffxp FirefoxProfilesDefaultImpl) GetProfilesList() ([]string, error) {
    result := make([]string, 0)
    files, err := ioutil.ReadDir(ffxp.ProfilesPath)
    if err != nil {
        return result, errors.New("impossible d'accéder au répertoire " + ffxp.ProfilesPath)
    }
    for _, f := range files {
        matched, _ := regexp.Match(`^.*\..*$`, []byte(f.Name()))
        if matched && f.IsDir() {
            result = append(result, f.Name())
        }
    }
    return result, nil
}

// GetProfilesPath returns the profiles path.
func (ffxp FirefoxProfilesDefaultImpl) GetProfilesPath() string {
    return ffxp.ProfilesPath
}
