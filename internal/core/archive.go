package core

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spec-kit/task-kit/internal/util"
)

// ExtractZip extracts zip to destDir
func ExtractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		fp := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fp, f.Mode()); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fp), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		func() {
			defer rc.Close()
			w, err := os.OpenFile(fp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, f.Mode())
			if err != nil {
				panic(err)
			}
			defer w.Close()
			if _, err := io.Copy(w, rc); err != nil {
				panic(err)
			}
		}()
	}
	return nil
}

// FlattenSingleTopDir flattens if the destDir contains exactly one top-level dir
// returns the actual root path containing content
func FlattenSingleTopDir(destDir string) (string, error) {
	entries, err := os.ReadDir(destDir)
	if err != nil {
		return "", err
	}
	if len(entries) == 1 && entries[0].IsDir() {
		src := filepath.Join(destDir, entries[0].Name())
		return src, nil
	}
	return destDir, nil
}

// CopyTree copies all files from src to dst (overwrites existing)
func CopyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		srcf, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcf.Close()
		dstf, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer dstf.Close()
		if _, err := io.Copy(dstf, srcf); err != nil {
			return err
		}
		return nil
	})
}

func FixScriptPermissions(root string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	var changed []string
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		lower := strings.ToLower(filepath.Base(p))
		if strings.HasSuffix(lower, ".sh") || lower == "gradlew" || lower == "mvnw" {
			info, err := os.Stat(p)
			if err != nil {
				return err
			}
			mode := info.Mode() | 0o111
			if err := os.Chmod(p, mode); err != nil {
				return err
			}
			changed = append(changed, p)
		}
		return nil
	})
	if err == nil && len(changed) > 0 {
		util.Debugf("executable set: %v\n", changed)
	}
	return err
}

// ExtractZipWithProgress extracts zip to destDir and reports combined progress in bytes
func ExtractZipWithProgress(zipPath, destDir string, onProgress func(done, total int64)) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Sum uncompressed sizes for all files
	var total int64
	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			total += int64(f.UncompressedSize64)
		}
	}
	var doneSoFar int64

	for _, f := range r.File {
		fp := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fp, f.Mode()); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fp), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		w, err := os.OpenFile(fp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, f.Mode())
		if err != nil {
			_ = rc.Close()
			return err
		}
		fileSize := int64(f.UncompressedSize64)
		var copyErr error
		if onProgress != nil && total > 0 && fileSize > 0 {
			_, copyErr = copyWithProgress(w, rc, fileSize, func(written, _ int64) {
				if onProgress != nil {
					onProgress(doneSoFar+written, total)
				}
			})
		} else {
			_, copyErr = io.Copy(w, rc)
		}
		cerr := rc.Close()
		err2 := w.Close()
		if copyErr != nil {
			return copyErr
		}
		if cerr != nil {
			return cerr
		}
		if err2 != nil {
			return err2
		}
		doneSoFar += fileSize
		if onProgress != nil && total > 0 {
			onProgress(doneSoFar, total)
		}
	}
	return nil
}

// CopyTreeWithProgress copies all files and reports combined progress in bytes
func CopyTreeWithProgress(src, dst string, onProgress func(done, total int64)) error {
	// Calculate total bytes
	var total int64
	err := filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, sErr := os.Stat(p)
		if sErr != nil {
			return sErr
		}
		total += info.Size()
		return nil
	})
	if err != nil {
		return err
	}
	var doneSoFar int64

	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		srcf, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcf.Close()
		info, _ := srcf.Stat()
		fileSize := info.Size()
		dstf, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer dstf.Close()
		var copyErr error
		if onProgress != nil && total > 0 && fileSize > 0 {
			_, copyErr = copyWithProgress(dstf, srcf, fileSize, func(written, _ int64) {
				if onProgress != nil {
					onProgress(doneSoFar+written, total)
				}
			})
		} else {
			_, copyErr = io.Copy(dstf, srcf)
		}
		if copyErr != nil {
			return copyErr
		}
		doneSoFar += fileSize
		if onProgress != nil && total > 0 {
			onProgress(doneSoFar, total)
		}
		return nil
	})
}