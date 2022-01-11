package profiles

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "runtime"
    "strings"

    "github.com/shirou/gopsutil/v3/process"
)

// FirefoxProfiles interfaces Firefox profiles related operations.
type FirefoxProfiles interface {

    // IsProfileUsed returns true if the specified profile is currently used.
    IsProfileUsed(profileName string) (bool, error)

    // GetProfilesList returns the list of existing folders in Firefox profiles dir
    // that match the profiles name pattern.
    GetProfilesList() ([]string, error)

    // GetProfilesPath returns the Firefox profiles path.
    GetProfilesPath() string

    // GetProfilesMatching returns the list of existing folders in Firefox profiles dir
    // that match the regex parameter.
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

func (ffxp FirefoxProfilesDefaultImpl) GetProfilesList() ([]string, error) {
    regex := regexp.MustCompile(`^.*\..*$`)
    return getDirsMatchingRegex(ffxp.ProfilesPath, regex)
}

func (ffxp FirefoxProfilesDefaultImpl) GetProfilesPath() string {
    return ffxp.ProfilesPath
}

func (ffxp FirefoxProfilesDefaultImpl) GetProfilesMatching(regex *regexp.Regexp) ([]string, error) {
    return getDirsMatchingRegex(ffxp.ProfilesPath, regex)
}

// getDirsMatchingRegex returns the list of folders located at root, which match the regex
func getDirsMatchingRegex(root string, regex *regexp.Regexp) ([]string, error) {
    result := make([]string, 0)
    files, err := os.ReadDir(root)
    if err != nil {
        return result, errors.New("impossible d'accéder au répertoire " + root)
    }
    for _, f := range files {
        matched := regex.MatchString(f.Name())
        if matched && f.IsDir() {
            result = append(result, f.Name())
        }
    }
    return result, nil
}
