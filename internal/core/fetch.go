package core

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spec-kit/task-kit/internal/util"
)

// DownloadZipFromURLWithProgress downloads a zip and reports progress via callback.
func DownloadZipFromURLWithProgress(rawURL, token string, skipTLS, debug bool, onProgress func(downloaded, total int64)) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("无效的 URL: %w", err)
	}

	// 支持本地文件: file:// 协议，便于本地回归测试
	if strings.EqualFold(u.Scheme, "file") {
		localPath := u.Path
		// 处理 Windows 路径: /C:/... -> C:/...
		if len(localPath) >= 4 && localPath[0] == '/' && localPath[2] == ':' {
			localPath = localPath[1:]
		}
		if p, err := url.PathUnescape(localPath); err == nil {
			localPath = p
		}
		st, err := os.Stat(localPath)
		if err != nil {
			return "", err
		}
		r, err := os.Open(localPath)
		if err != nil {
			return "", err
		}
		defer r.Close()

		f, err := os.CreateTemp("", "task-kit-*.zip")
		if err != nil {
			return "", err
		}
		defer f.Close()

		total := st.Size()
		if onProgress == nil || total <= 0 {
			if _, err := io.Copy(f, r); err != nil {
				return "", err
			}
			return f.Name(), nil
		}
		if _, err := copyWithProgress(f, r, total, onProgress); err != nil {
			return "", err
		}
		return f.Name(), nil
	}

	client := NewHTTPClient(skipTLS)
	// Default accept for assets
	headers := map[string]string{"Accept": "application/octet-stream"}
	if token != "" && strings.EqualFold(u.Host, "github.com") {
		headers["Authorization"] = "Bearer " + token
	}

	if debug {
		util.Debugf("download zip from %s\n", rawURL)
	}

	resp, err := client.R().SetHeaders(headers).SetDoNotParseResponse(true).Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.RawBody().Close()
	if resp.IsError() {
		return "", fmt.Errorf("下载失败: %s", resp.Status())
	}

	var total int64 = -1
	if cl := resp.Header().Get("Content-Length"); cl != "" {
		if v, err := strconv.ParseInt(cl, 10, 64); err == nil {
			total = v
		}
	}

	f, err := os.CreateTemp("", "task-kit-*.zip")
	if err != nil {
		return "", err
	}
	defer f.Close()

	if onProgress == nil || total <= 0 {
		if _, err := io.Copy(f, resp.RawBody()); err != nil {
			return "", err
		}
		return f.Name(), nil
	}

	if _, err := copyWithProgress(f, resp.RawBody(), total, onProgress); err != nil {
		return "", err
	}
	return f.Name(), nil
}
