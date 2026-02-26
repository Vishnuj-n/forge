package remote

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// DownloadRepoZip downloads the given URL into a temporary file and
// returns the path to the downloaded zip file. Caller must remove the file.
func DownloadRepoZip(url string) (string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download repository: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to download repository: http %d", resp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "forge-templates-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	return tmp.Name(), nil
}

// detectZipPrefix returns the top-level prefix present in paths inside the zip,
// e.g. "forge-templates-main/". If none found, empty string is returned.
func detectZipPrefix(r *zip.Reader) string {
	for _, f := range r.File {
		name := f.Name
		if idx := strings.Index(name, "/"); idx > 0 {
			return name[:idx+1]
		}
	}
	return ""
}

// ListTopLevelTemplates lists the top-level directories (template names)
// contained in the zip archive. It does not extract files.
func ListTopLevelTemplates(zipPath string) ([]string, error) {
	z, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip: %w", err)
	}
	defer z.Close()

	prefix := detectZipPrefix(&z.Reader)
	names := map[string]struct{}{}

	for _, f := range z.File {
		name := f.Name
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		rest := strings.TrimPrefix(name, prefix)
		parts := strings.SplitN(rest, "/", 2)
		if len(parts) > 0 && parts[0] != "" {
			names[parts[0]] = struct{}{}
		}
	}

	out := make([]string, 0, len(names))
	for n := range names {
		out = append(out, n)
	}
	return out, nil
}

// extractPrefixedFiles extracts all files under zipPrefix/templateName/ into destDir.
// It validates presence of template.yaml inside the extracted content.
func extractPrefixedFiles(z *zip.Reader, zipPrefix, templateName, destDir string) error {
	wantedPrefix := zipPrefix + templateName + "/"
	foundAny := false
	foundYAML := false

	// Remove existing destination and recreate
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("failed to remove existing template dir: %w", err)
	}

	for _, f := range z.File {
		if !strings.HasPrefix(f.Name, wantedPrefix) {
			continue
		}
		foundAny = true
		rel := strings.TrimPrefix(f.Name, wantedPrefix)
		if rel == "" {
			// This is the template root directory entry
			continue
		}
		targetPath := filepath.Join(destDir, filepath.FromSlash(rel))

		if strings.HasSuffix(f.Name, "/") {
			// Directory
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create dir %s: %w", targetPath, err)
			}
			continue
		}

		// Ensure parent dir exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent dir %s: %w", filepath.Dir(targetPath), err)
		}

		// Helper to ensure files are closed promptly
		err := func() error {
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open zip entry: %w", err)
			}
			defer rc.Close()

			out, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}
			defer out.Close()

			if _, err := io.Copy(out, rc); err != nil {
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
			return nil
		}()

		if err != nil {
			return err
		}

		if filepath.Base(targetPath) == "template.yaml" {
			foundYAML = true
		}
	}

	if !foundAny {
		return errors.New("template not found in archive")
	}
	if !foundYAML {
		// Clean up partial extraction
		os.RemoveAll(destDir)
		return errors.New("invalid template: missing template.yaml")
	}
	return nil
}

// InstallSingleTemplate extracts the requested template from the zip archive into destParentDir/templateName
func InstallSingleTemplate(zipPath, templateName, destParentDir string) error {
	zf, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer zf.Close()

	prefix := detectZipPrefix(&zf.Reader)
	if prefix == "" {
		return errors.New("unable to detect zip prefix")
	}

	destDir := filepath.Join(destParentDir, templateName)
	if err := extractPrefixedFiles(&zf.Reader, prefix, templateName, destDir); err != nil {
		return err
	}
	return nil
}

// InstallAllTemplates extracts top-level directories that contain template.yaml
// into destParentDir. It returns a slice of installed template names and an error.
func InstallAllTemplates(zipPath, destParentDir string) ([]string, error) {
	zf, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip: %w", err)
	}
	defer zf.Close()

	prefix := detectZipPrefix(&zf.Reader)
	if prefix == "" {
		return nil, errors.New("unable to detect zip prefix")
	}

	// Collect candidate top-level names
	candidates := map[string]struct{}{}
	for _, f := range zf.File {
		name := f.Name
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		rest := strings.TrimPrefix(name, prefix)
		parts := strings.SplitN(rest, "/", 2)
		if len(parts) > 0 && parts[0] != "" {
			candidates[parts[0]] = struct{}{}
		}
	}

	installed := make([]string, 0, len(candidates))
	for name := range candidates {
		destDir := filepath.Join(destParentDir, name)
		err := extractPrefixedFiles(&zf.Reader, prefix, name, destDir)
		if err != nil {
			// skip invalid templates, but continue
			continue
		}
		installed = append(installed, name)
	}

	return installed, nil
}
