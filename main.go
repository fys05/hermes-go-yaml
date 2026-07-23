// hermes-go-yaml — a YAML file format validator.
//
// Usage:
//
//	yaml-validator file1.yaml file2.yaml     # CLI: validate files
//	cat config.yaml | yaml-validator         # CLI: validate stdin
//	yaml-validator --serve :8080             # Web UI mode
//
// In web mode, visit http://host:8080/ for an interactive YAML validator.
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fys05/hermes-go-yaml/validator"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--serve" {
		addr := ":8080"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		serve(addr)
		return
	}

	os.Exit(run())
}

func run() int {
	if len(os.Args) > 1 {
		return validateFiles(os.Args[1:])
	}
	return validateStdin()
}

func validateFiles(paths []string) int {
	allValid := true
	for _, path := range paths {
		r := validator.ValidateFile(path)
		if r.Valid {
			fmt.Printf("✔ %s — valid YAML\n", r.Path)
		} else {
			allValid = false
			if r.Line > 0 {
				fmt.Printf("✘ %s:%d:%d — %s\n", r.Path, r.Line, r.Column, r.Error)
			} else {
				fmt.Printf("✘ %s — %s\n", r.Path, r.Error)
			}
		}
	}
	if allValid {
		return 0
	}
	return 1
}

func validateStdin() int {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
		return 1
	}
	r := validator.ValidateData(data)
	if r.Valid {
		fmt.Println("✔ stdin — valid YAML")
		return 0
	}
	if r.Line > 0 {
		fmt.Printf("✘ stdin:%d:%d — %s\n", r.Line, r.Column, r.Error)
	} else {
		fmt.Printf("✘ stdin — %s\n", r.Error)
	}
	return 1
}

// ---- Web UI ----

const pageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>YAML Validator</title>
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #0d1117; color: #c9d1d9; min-height: 100vh; display: flex; flex-direction: column; }
  header { background: #161b22; border-bottom: 1px solid #30363d; padding: 12px 20px; display: flex; align-items: center; gap: 10px; }
  header h1 { font-size: 18px; color: #58a6ff; }
  header .badge { font-size: 12px; background: #238636; color: #fff; padding: 2px 8px; border-radius: 10px; }
  main { flex: 1; display: flex; padding: 20px; gap: 16px; max-width: 1400px; margin: 0 auto; width: 100%; }
  .panel { flex: 1; display: flex; flex-direction: column; min-width: 0; }
  .panel h2 { font-size: 14px; margin-bottom: 8px; color: #8b949e; text-transform: uppercase; letter-spacing: 0.5px; }
  textarea { flex: 1; background: #0d1117; color: #c9d1d9; border: 1px solid #30363d; border-radius: 6px; padding: 12px; font-family: 'SF Mono', 'Fira Code', monospace; font-size: 13px; line-height: 1.5; resize: none; outline: none; }
  textarea:focus { border-color: #58a6ff; box-shadow: 0 0 0 3px rgba(88,166,255,0.15); }
  #result { flex: 1; background: #0d1117; border: 1px solid #30363d; border-radius: 6px; padding: 12px; font-family: 'SF Mono', 'Fira Code', monospace; font-size: 13px; overflow: auto; white-space: pre-wrap; }
  #result.valid { border-color: #238636; }
  #result.invalid { border-color: #da3633; }
  .valid-text { color: #3fb950; }
  .invalid-text { color: #f85149; }
  .actions { display: flex; gap: 8px; margin-bottom: 12px; }
  button { background: #238636; color: #fff; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; font-size: 13px; font-weight: 500; }
  button:hover { background: #2ea043; }
  button:disabled { opacity: 0.5; cursor: default; }
  button.secondary { background: #21262d; border: 1px solid #30363d; }
  button.secondary:hover { background: #30363d; }
  .status { font-size: 12px; color: #8b949e; margin-left: auto; align-self: center; }
  .error-line { color: #f85149; }
  footer { text-align: center; padding: 12px; font-size: 12px; color: #484f58; border-top: 1px solid #30363d; }
</style>
</head>
<body>
<header>
  <h1>🔍 YAML Validator</h1>
  <span class="badge">hermes-go-yaml</span>
  <span class="status" id="status">Ready</span>
</header>
<main>
  <div class="panel">
    <h2>📝 Input</h2>
    <div class="actions">
      <button onclick="validate()">Validate</button>
      <button class="secondary" onclick="clearAll()">Clear</button>
      <button class="secondary" onclick="loadSample()">Sample</button>
    </div>
    <textarea id="input" placeholder="Paste your YAML here..." spellcheck="false"># Example
name: hello-world
version: "1.0"
items:
  - one
  - two</textarea>
  </div>
  <div class="panel">
    <h2>📋 Result</h2>
    <div id="result" class="valid">
      <span class="valid-text">✔ Valid YAML</span>
    </div>
  </div>
</main>
<footer>hermes-go-yaml · GitHub Actions CI/CD · ghcr.io/fys05/hermes-go-yaml</footer>
<script>
const input = document.getElementById('input');
const result = document.getElementById('result');
const status = document.getElementById('status');
let timer;

input.addEventListener('input', () => {
  clearTimeout(timer);
  timer = setTimeout(validate, 500);
});

async function validate() {
  const yaml = input.value.trim();
  if (!yaml) { result.className = ''; result.innerHTML = 'Waiting for input...'; return; }

  status.textContent = 'Validating...';
  try {
    const resp = await fetch('/validate', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({yaml: yaml})
    });
    const data = await resp.json();
    if (data.valid) {
      result.className = 'valid';
      result.innerHTML = '<span class="valid-text">✔ Valid YAML</span>';
      status.textContent = '✓ Valid';
    } else {
      result.className = 'invalid';
      let msg = data.error;
      if (data.line > 0) {
        msg = '<span class="error-line">Line ' + data.line + ', Col ' + data.column + ':</span> ' + data.error;
      }
      result.innerHTML = '<span class="invalid-text">✘ ' + msg + '</span>';
      status.textContent = '✗ Invalid';
    }
  } catch (e) {
    result.className = 'invalid';
    result.innerHTML = '<span class="invalid-text">Request failed: ' + e.message + '</span>';
    status.textContent = 'Error';
  }
}

function clearAll() {
  input.value = '';
  result.className = '';
  result.innerHTML = 'Waiting for input...';
  status.textContent = 'Ready';
}

function loadSample() {
  input.value = '# Paste your YAML here...\nname: hello-world\nversion: "1.0"';
  validate();
}

// Validate on load
validate();
</script>
</body>
</html>`

var pageTemplate = template.Must(template.New("page").Parse(pageHTML))

type validateRequest struct {
	YAML string `json:"yaml"`
}

type validateResponse struct {
	Valid  bool   `json:"valid"`
	Error  string `json:"error,omitempty"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
}

func serve(addr string) {
	mux := http.NewServeMux()

	// Web UI
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		pageTemplate.Execute(w, nil)
	})

	// API: validate YAML
	mux.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		var req validateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		v := validator.ValidateData([]byte(req.YAML))
		resp := validateResponse{
			Valid:  v.Valid,
			Error:  v.Error,
			Line:   v.Line,
			Column: v.Column,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	fmt.Printf("🌐 YAML Validator Web UI: http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
