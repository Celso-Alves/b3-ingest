package ingestion

import (
	"archive/zip"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// DownloadAndUnzipLast7Workdays downloads and unzips the last 7 workdays' files to destDir. It can be cancelled via ctx.
func DownloadAndUnzipLast7Workdays(ctx context.Context, destDir string, logf func(string, ...interface{})) error {
	const baseURL = "https://arquivos.b3.com.br/rapinegocios/tickercsv/"
	dates := lastNWorkdays(7)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	client := &http.Client{}
	for _, d := range dates {
		select {
		case <-ctx.Done():
			logf("Download cancelled by user.")
			return ctx.Err()
		default:
		}
		url := baseURL + d.Format("2006-01-02")
		zipPath := filepath.Join(destDir, d.Format("2006-01-02")+".zip")
		logf("Downloading %s...", url)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			logf("Failed to create request for %s: %v", url, err)
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			logf("Failed to download %s: %v", url, err)
			continue
		}
		if resp.StatusCode != 200 {
			logf("No file for %s (HTTP %d)", d.Format("2006-01-02"), resp.StatusCode)
			resp.Body.Close()
			continue
		}
		f, err := os.Create(zipPath)
		if err != nil {
			resp.Body.Close()
			logf("Failed to create file: %v", err)
			continue
		}
		_, err = io.Copy(f, resp.Body)
		resp.Body.Close()
		f.Close()
		if err != nil {
			logf("Failed to save zip: %v", err)
			continue
		}
		if err := unzip(zipPath, destDir, logf); err != nil {
			logf("Failed to unzip %s: %v", zipPath, err)
			continue
		}
		os.Remove(zipPath)
		logf("Downloaded and extracted %s", d.Format("2006-01-02"))
	}
	return nil
}

// lastNWorkdays returns the last n workdays, excluding today.
func lastNWorkdays(n int) []time.Time {
	var days []time.Time
	now := time.Now()
	// Exclude today
	d := now.AddDate(0, 0, -1)
	for len(days) < n {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			days = append(days, d)
		}
		d = d.AddDate(0, 0, -1)
	}
	// reverse to oldest first
	for i, j := 0, len(days)-1; i < j; i, j = i+1, j-1 {
		days[i], days[j] = days[j], days[i]
	}
	return days
}

func unzip(src, dest string, logf func(string, ...interface{})) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
		logf("Extracted %s", fpath)
	}
	return nil
}
