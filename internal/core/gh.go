package core

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/spec-kit/task-kit/internal/util"
)

type DownloadOptions struct {
	Repo          string // owner/repo 或 https://github.com/owner/repo
	Branch        string // ref
	Token         string
	SkipTLSVerify bool
	Debug         bool
}

// GitHub release schema (partial)
type GHAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type GHRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []GHAsset `json:"assets"`
}

var reOwnerRepo = regexp.MustCompile(`^([A-Za-z0-9_.-]+)/([A-Za-z0-9_.-]+)$`)

func ResolveGitHubToken(flag string) string {
	if flag != "" {
		return flag
	}
	if v := os.Getenv("GH_TOKEN"); v != "" {
		return v
	}
	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		return v
	}
	return ""
}

// parseRepo returns owner, repo
func parseRepo(repo string) (string, string, error) {
	if m := reOwnerRepo.FindStringSubmatch(repo); len(m) == 3 {
		return m[1], m[2], nil
	}
	u, err := url.Parse(repo)
	if err != nil {
		return "", "", fmt.Errorf("无法解析 repo: %w", err)
	}
	if u.Host != "github.com" {
		return "", "", fmt.Errorf("仅支持 github.com 源: %s", u.Host)
	}
	p := strings.TrimPrefix(path.Clean(u.Path), "/")
	p = strings.TrimSuffix(p, ".git")
	parts := strings.Split(p, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("repo 路径不完整: %s", p)
	}
	return parts[0], parts[1], nil
}

// DownloadTemplateFromGitHubWithProgress downloads a zipball and reports progress via callback
func DownloadTemplateFromGitHubWithProgress(opt DownloadOptions, onProgress func(downloaded, total int64)) (string, error) {
	owner, repo, err := parseRepo(opt.Repo)
	if err != nil {
		return "", err
	}
	ref := opt.Branch
	if ref == "" {
		ref = "main"
	}
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/%s", owner, repo, ref)

	client := NewHTTPClient(opt.SkipTLSVerify)
	headers := GetGitHubAuthHeaders(opt.Token)
	if opt.Debug {
		util.Debugf("download %s\n", endpoint)
	}
	resp, err := client.R().
		SetHeaders(headers).
		SetDoNotParseResponse(true).
		Get(endpoint)
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

	file, err := os.CreateTemp("", "task-kit-*.zip")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if onProgress == nil || total <= 0 {
		if _, err := io.Copy(file, resp.RawBody()); err != nil {
			return "", err
		}
		return file.Name(), nil
	}

	if _, err := copyWithProgress(file, resp.RawBody(), total, onProgress); err != nil {
		return "", err
	}
	return file.Name(), nil
}

// GetLatestReleaseAssetURL fetches latest release and selects asset matching ai+script (zip)
// Returns (assetName, browserDownloadURL)
func GetLatestReleaseAssetURL(repo, ai, script, token string, skipTLS, debug bool) (string, string, error) {
	owner, r, err := parseRepo(repo)
	if err != nil {
		return "", "", err
	}
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, r)

	client := NewHTTPClient(skipTLS)
	var release GHRelease
	if debug {
		util.Debugf("fetch release %s\n", endpoint)
	}
	resp, err := client.R().
		SetHeaders(GetGitHubAuthHeaders(token)).
		SetResult(&release).
		Get(endpoint)
	if err != nil {
		return "", "", err
	}
	if resp.IsError() {
		return "", "", fmt.Errorf("获取 release 失败: %s", resp.Status())
	}

	ai = strings.ToLower(ai)
	script = strings.ToLower(script)
	var scriptTokens []string
	switch script {
	case "powershell", "pwsh", "ps":
		scriptTokens = []string{"ps", "powershell", "pwsh"}
	case "bash", "sh":
		scriptTokens = []string{"bash", "sh"}
	case "zsh":
		scriptTokens = []string{"zsh"}
	default:
		scriptTokens = []string{script}
	}

	// prefer exact ai+script matches
	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if !strings.HasSuffix(name, ".zip") {
			continue
		}
		if ai != "" && !strings.Contains(name, "-"+ai+"-") {
			continue
		}
		for _, st := range scriptTokens {
			if strings.Contains(name, "-"+st+"-") {
				if debug {
					util.Debugf("selected asset: %s (%s)\n", asset.Name, release.TagName)
				}
				return asset.Name, asset.BrowserDownloadURL, nil
			}
		}
	}

	// fallback: first zip asset
	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.HasSuffix(name, ".zip") {
			if debug {
				util.Debugf("fallback asset: %s (%s)\n", asset.Name, release.TagName)
			}
			return asset.Name, asset.BrowserDownloadURL, nil
		}
	}

	return "", "", fmt.Errorf("未找到匹配的 zip 资产 (ai=%s, script=%s)", ai, script)
}
