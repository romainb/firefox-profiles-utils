package profiles

import (
    "os"
    "path"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
)

// TestFirefoxProfilesDefaultImpl_GetProfilesList vérifie que la fonction remonte bien des dossiers correspondant à des profils
// à partir d'un dossier temporaire créé à l'occasion du test.
func TestFirefoxProfilesDefaultImpl_GetProfilesList(t *testing.T) {
    profileFolders := []string{"profile.1", "profile.2", "profile.3", "not_a_profile"}
    testTmpDir := path.Join(os.TempDir(), "profiles")

    // Création du dossier temporaire
    err := os.Mkdir(testTmpDir, 0755)
    if err != nil {
        t.Error(err)
    }

    // Création des répertoires dans le dossier temporaire
    for _, folder := range profileFolders {
        err := os.Mkdir(path.Join(testTmpDir, folder), 0755)
        if err != nil {
            t.Error(err)
        }
    }

    ffxp := NewWithCustomPath(filepath.Join(os.TempDir(), "profiles"))
    foundProfiles, _ := ffxp.GetProfilesList()
    assert.Contains(t, foundProfiles, "profile.1", "profile.1 should be returned")
    assert.Contains(t, foundProfiles, "profile.2", "profile.2 should be returned")
    assert.Contains(t, foundProfiles, "profile.3", "profile.3 should be returned")
    assert.NotContains(t, foundProfiles, "not_a_profile", "profile.3 should be returned")

    // Suppression du dossier temporaire
    err = os.RemoveAll(testTmpDir)
    if err != nil {
        t.Error(err)
    }
}
