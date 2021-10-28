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
    GetProfilesMatching(regex *regexp.Regexp) ([]string, error)
}

// FirefoxProfilesDefaultImpl implements FirefoxProfiles.
type FirefoxProfilesDefaultImpl struct {
    // ProfilesPath Firefox profiles directory path.
    ProfilesPath string
}

// NewWithDefaultPath builds a FirefoxProfiles instance with the default path according to the OS.
func NewWithDefaultPath() FirefoxProfiles {
    userHome, _ := os.UserHomeDir()
    switch currentOs := runtime.GOOS; currentOs {
    case "windows":
        return FirefoxProfilesDefaultImpl{filepath.Join(os.Getenv("APPDATA"), "Mozilla", "Firefox", "Profiles")}
    case "darwin":
        return FirefoxProfilesDefaultImpl{
            filepath.Join(userHome, "Library", "Application Support", "Firefox", "Profiles")}
    case "linux":
        return FirefoxProfilesDefaultImpl{filepath.Join(userHome, ".mozilla", "firefox")}
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
        if strings.HasPrefix(name, "firefox") && strings.Contains(args, fmt.Sprintf("-P %s", profileName)) {
            return true, nil
        }
    }
    return false, nil
}

// GetProfilesList returns the list of existing folders in Firefox profiles dir
// that match the profiles name pattern.
func (ffxp FirefoxProfilesDefaultImpl) GetProfilesList() ([]string, error) {
    regex := regexp.MustCompile(`^.*\..*$`)
    return getDirsMatchingRegex(ffxp.ProfilesPath, regex)
}

// GetProfilesPath returns the Firefox profiles path.
func (ffxp FirefoxProfilesDefaultImpl) GetProfilesPath() string {
    return ffxp.ProfilesPath
}

// GetProfilesMatching returns the list of existing folders in Firefox profiles dir
// that match the regex parameter.
func (ffxp FirefoxProfilesDefaultImpl) GetProfilesMatching(regex *regexp.Regexp) ([]string, error) {
    return getDirsMatchingRegex(ffxp.ProfilesPath, regex)
}

// getDirsMatchingRegex returns the list of folders located at root, which match the regex
func getDirsMatchingRegex(root string, regex *regexp.Regexp) ([]string, error) {
    result := make([]string, 0)
    files, err := ioutil.ReadDir(root)
    if err != nil {
        return result, errors.New("impossible d'accéder au répertoire " + root)
    }
    for _, f := range files {
        matched := regex.Match([]byte(f.Name()))
        if matched && f.IsDir() {
            result = append(result, f.Name())
        }
    }
    return result, nil
}
