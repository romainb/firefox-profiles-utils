package profiles

import (
    "os"
    "path"
    "path/filepath"
    "regexp"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

var (
    profileFolders  = []string{"profile.1", "profile.2", "profile.3", "not_a_profile"}
    profileFolders2 = []string{"profile.4", "not_a_profile"}
    testTmpDir      = path.Join(os.TempDir(), "profiles")
    testTmpDir2     = path.Join(os.TempDir(), "profiles2")
)

// FirefoxProfilesSuite déclare la suite de tests.
type FirefoxProfilesSuite struct {
    suite.Suite
}

// TestFirefoxProfilesSuite exécute la suite de tests.
func TestFirefoxProfilesSuite(t *testing.T) {
    suite.Run(t, new(FirefoxProfilesSuite))
}

// SetupSuite est exécuté une fois au lancement de la suite de test, avant les tests
func (s *FirefoxProfilesSuite) SetupSuite() {}

// TearDownSuite s'exécute après la suite
func (s *FirefoxProfilesSuite) TearDownSuite() {}

// SetupTest est lancé avant chaque test.
func (s *FirefoxProfilesSuite) SetupTest() {
    _ = os.Mkdir(testTmpDir, 0755)
    _ = os.Mkdir(testTmpDir2, 0755)
    for _, folder := range profileFolders {
        _ = os.Mkdir(path.Join(testTmpDir, folder), 0755)
    }
    for _, folder := range profileFolders2 {
        _ = os.Mkdir(path.Join(testTmpDir2, folder), 0755)
    }
}

// TearDownSuite s'exécute après la suite
func (s *FirefoxProfilesSuite) TearDownTest() {
    // Suppression du dossier temporaire
    err := os.RemoveAll(testTmpDir)
    if err != nil {
        s.Error(err)
    }
}

// TestFirefoxProfilesDefaultImpl_GetProfilesList vérifie que la fonction remonte
// bien des dossiers correspondant à des profils à partir d'un dossier temporaire
// créé à l'occasion du test.
func (s *FirefoxProfilesSuite) TestFirefoxProfilesDefaultImpl_GetProfilesList() {
    ffxp := NewWithCustomPath([]string{testTmpDir, testTmpDir2})
    foundProfiles, _ := ffxp.GetProfilesList()
    assert.Contains(s.T(), foundProfiles, filepath.Join(testTmpDir, "profile.1"), "profile.1 should be returned")
    assert.Contains(s.T(), foundProfiles, filepath.Join(testTmpDir, "profile.2"), "profile.2 should be returned")
    assert.Contains(s.T(), foundProfiles, filepath.Join(testTmpDir, "profile.3"), "profile.3 should be returned")
    assert.Contains(s.T(), foundProfiles, filepath.Join(testTmpDir2, "profile.4"), "profile.3 should be returned")
    assert.NotContains(s.T(), foundProfiles, "not_a_profile", "not_a_profiles should not be returned")
}

// On vérifie que le profil retourné correspond à la regex passée en paramètre
func (s *FirefoxProfilesSuite) TestGetProfilesDirMatching() {
    ffxp := NewWithCustomPath([]string{testTmpDir})
    foundProfiles, err := ffxp.GetProfilesMatching(regexp.MustCompile(`^not_a.*$`))
    assert.Nil(s.T(), err)
    assert.Contains(s.T(), foundProfiles, filepath.Join(testTmpDir, "not_a_profile"), "not_a_profile should be returned")
}

// Si le répertoire de profils n'existe pas on retourne une collection vide et pas d'erreur
func (s *FirefoxProfilesSuite) TestGetProfilesList_whenFolderDoesNotExists() {
    _ = os.RemoveAll(testTmpDir)
    ffxp := NewWithCustomPath([]string{testTmpDir})
    foundProfiles, err := ffxp.GetProfilesList()
    assert.Emptyf(s.T(), foundProfiles, "foundProfiles should be empty")
    assert.Nil(s.T(), err, "there should not be an error")
}

// Si le répertoire de profils existe mais est vide, on retourne une collection vide et pas d'erreur
func (s *FirefoxProfilesSuite) TestGetProfilesList_whenFolderExistsButIsEmpty() {
    _ = os.RemoveAll(testTmpDir)
    _ = os.Mkdir(testTmpDir, 0755)
    ffxp := NewWithCustomPath([]string{testTmpDir})
    foundProfiles, err := ffxp.GetProfilesList()
    assert.Emptyf(s.T(), foundProfiles, "foundProfiles should be empty")
    assert.Nil(s.T(), err, "there should not be an error")
}

// Si les deuxs répertoires de profils existent mais sont vides, on retourne une collection vide et pas d'erreur
func (s *FirefoxProfilesSuite) TestGetProfilesList_whenFoldersExistButAreEmpty() {
    _ = os.RemoveAll(testTmpDir)
    _ = os.RemoveAll(testTmpDir2)
    _ = os.Mkdir(testTmpDir, 0755)
    _ = os.Mkdir(testTmpDir2, 0755)
    ffxp := NewWithCustomPath([]string{testTmpDir, testTmpDir2})
    foundProfiles, err := ffxp.GetProfilesList()
    assert.Emptyf(s.T(), foundProfiles, "foundProfiles should be empty")
    assert.Nil(s.T(), err, "there should not be an error")
}
