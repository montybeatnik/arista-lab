package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ======= Data models (same as before) =======

type InspectResult map[string][]ContainerInfo

type ContainerInfo struct {
	LabName     string `json:"lab_name"`
	LabPath     string `json:"labPath"`
	AbsLabPath  string `json:"absLabPath"`
	Name        string `json:"name"`
	ContainerID string `json:"container_id"`
	Image       string `json:"image"`
	Kind        string `json:"kind"`
	State       string `json:"state"`
	Status      string `json:"status"`
	IPv4        string `json:"ipv4_address"`
	IPv6        string `json:"ipv6_address"`
	Owner       string `json:"owner"`
}

// ======= Config =======

type serverCfg struct {
	Listen  string
	BaseDir string // lab files must live under here
}

func (c serverCfg) sanitizeLabPath(p string) (string, error) {
	if p == "" {
		return "", errors.New("lab file required")
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(c.BaseDir, p)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	baseAbs, _ := filepath.Abs(c.BaseDir)
	if abs != baseAbs && !strings.HasPrefix(abs, baseAbs+string(os.PathSeparator)) {
		return "", errors.New("lab file must be under basedir")
	}
	info, err := os.Stat(abs)
	if err != nil || info.IsDir() {
		return "", errors.New("lab file not found")
	}
	return abs, nil
}

// ======= Handlers =======

type pageData struct {
	BaseDir string
}

//go:embed web/templates/*.tmpl
var tplFS embed.FS

//go:embed web/static/*
var staticFS embed.FS

// parse templates (explicit order to be clear)
func makeTemplate() *template.Template {
	return template.Must(template.ParseFS(
		tplFS,
		"web/templates/layout.tmpl",
		"web/templates/index.tmpl",
	))
}

func indexHandler(cfg serverCfg, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		if err := t.ExecuteTemplate(&buf, "layout", pageData{BaseDir: cfg.BaseDir}); err != nil {
			// nothing has been written to the client yet → safe to send an error
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = buf.WriteTo(w)
	}
}

type inspectReq struct {
	Lab        string `json:"lab"`
	UseSudo    bool   `json:"sudo"`
	TimeoutSec int    `json:"timeoutSec"`
}

type inspectResp struct {
	OK      bool            `json:"ok"`
	Error   string          `json:"error,omitempty"`
	LabKey  string          `json:"labKey,omitempty"`
	Nodes   []ContainerInfo `json:"nodes,omitempty"`
	RawJSON json.RawMessage `json:"rawJson,omitempty"`
}

func runInspect(ctx context.Context, labPath string, useSudo bool) ([]byte, error) {
	args := []string{"containerlab", "inspect", "-t", labPath, "--format", "json"}
	if useSudo {
		args = append([]string{"sudo", "-n"}, args...)
	}
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

func inspectHandler(cfg serverCfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST; if you clicked and nothing happened before, it’s
		// often script not loaded or method mismatch. We enforce POST here.
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req inspectReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, inspectResp{OK: false, Error: "bad JSON: " + err.Error()})
			return
		}

		labAbs, err := cfg.sanitizeLabPath(req.Lab)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, inspectResp{OK: false, Error: err.Error()})
			return
		}

		timeout := time.Duration(req.TimeoutSec) * time.Second
		if timeout <= 0 || timeout > 60*time.Second {
			timeout = 15 * time.Second
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		out, err := runInspect(ctx, labAbs, req.UseSudo)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, inspectResp{OK: false, Error: "inspect failed: " + err.Error()})
			return
		}

		var parsed InspectResult
		if err := json.Unmarshal(out, &parsed); err != nil {
			writeJSON(w, http.StatusInternalServerError, inspectResp{OK: false, Error: "parse failed: " + err.Error()})
			return
		}

		var key string
		var nodes []ContainerInfo
		for k, v := range parsed {
			key, nodes = k, v
			break
		}

		writeJSON(w, http.StatusOK, inspectResp{
			OK:      true,
			LabKey:  key,
			Nodes:   nodes,
			RawJSON: out,
		})
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	cfg := serverCfg{
		Listen:  ":8080",
		BaseDir: "/home/ubuntu/lab",
	}

	// Templates
	t := makeTemplate()

	// Static files (embed FS subdir)
	static, _ := fs.Sub(staticFS, "web/static")
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))

	// Pages & API
	mux.HandleFunc("/", indexHandler(cfg, t))
	mux.HandleFunc("/inspect", inspectHandler(cfg))

	srv := &http.Server{
		Addr:              cfg.Listen,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	println("listening on", cfg.Listen, "basedir:", cfg.BaseDir)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
