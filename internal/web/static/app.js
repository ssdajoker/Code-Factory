// Code Factory Web UI

document.addEventListener('DOMContentLoaded', () => {
    loadStatus();
    loadFiles();
    setupModeButtons();
});

async function loadStatus() {
    try {
        const resp = await fetch('/api/status');
        const data = await resp.json();
        document.getElementById('status').innerHTML = `
            <p>Status: <span class="success">${data.status}</span></p>
            <p>Version: ${data.version}</p>
        `;
    } catch (err) {
        console.error('Failed to load status:', err);
    }
}

async function loadFiles() {
    try {
        const [contracts, reports] = await Promise.all([
            fetch('/api/contracts').then(r => r.json()),
            fetch('/api/reports').then(r => r.json())
        ]);

        const contractsList = document.getElementById('contracts-list');
        contractsList.innerHTML = (contracts.contracts || []).map(f => 
            `<li>${f}</li>`
        ).join('') || '<li>No contracts yet</li>';

        const reportsList = document.getElementById('reports-list');
        reportsList.innerHTML = (reports.reports || []).map(f => 
            `<li>${f}</li>`
        ).join('') || '<li>No reports yet</li>';
    } catch (err) {
        console.error('Failed to load files:', err);
    }
}

function setupModeButtons() {
    document.querySelectorAll('.mode-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.mode-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            showModeForm(btn.dataset.mode);
        });
    });
}

function showModeForm(mode) {
    const content = document.getElementById('content');
    
    switch(mode) {
        case 'intake':
            content.innerHTML = `
                <h2>📝 INTAKE - Capture Vision</h2>
                <form id="intake-form">
                    <div class="form-group">
                        <label>Project Name</label>
                        <input type="text" name="project_name" required>
                    </div>
                    <div class="form-group">
                        <label>Description</label>
                        <textarea name="description" required></textarea>
                    </div>
                    <div class="form-group">
                        <label>Target Users</label>
                        <textarea name="target_users"></textarea>
                    </div>
                    <div class="form-group">
                        <label>Core Features (one per line)</label>
                        <textarea name="core_features"></textarea>
                    </div>
                    <div class="form-group">
                        <label>Technical Constraints</label>
                        <textarea name="technical_constraints"></textarea>
                    </div>
                    <div class="form-group">
                        <label>Success Criteria</label>
                        <textarea name="success_criteria"></textarea>
                    </div>
                    <button type="submit" class="btn">Generate Spec</button>
                </form>
                <div id="result"></div>
            `;
            document.getElementById('intake-form').addEventListener('submit', handleIntake);
            break;

        case 'review':
            content.innerHTML = `
                <h2>🔍 REVIEW - Check Code</h2>
                <form id="review-form">
                    <div class="form-group">
                        <label>Spec File Path</label>
                        <input type="text" name="spec_file" placeholder="contracts/vision_spec.md" required>
                    </div>
                    <div class="form-group">
                        <label>Code Paths (comma-separated)</label>
                        <input type="text" name="code_paths" placeholder="internal/,cmd/" required>
                    </div>
                    <button type="submit" class="btn">Run Review</button>
                </form>
                <div id="result"></div>
            `;
            document.getElementById('review-form').addEventListener('submit', handleReview);
            break;

        case 'rescue':
            content.innerHTML = `
                <h2>🆘 RESCUE - Reverse Engineer</h2>
                <form id="rescue-form">
                    <div class="form-group">
                        <label>Codebase Path</label>
                        <input type="text" name="codebase_path" placeholder="." required>
                    </div>
                    <button type="submit" class="btn">Scan Codebase</button>
                </form>
                <div id="result"></div>
            `;
            document.getElementById('rescue-form').addEventListener('submit', handleRescue);
            break;

        case 'change-order':
            content.innerHTML = `
                <h2>📋 CHANGE ORDER - Track Drift</h2>
                <form id="change-order-form">
                    <div class="form-group">
                        <label>Spec File Path</label>
                        <input type="text" name="spec_file" placeholder="contracts/vision_spec.md" required>
                    </div>
                    <div class="form-group">
                        <label>Codebase Path</label>
                        <input type="text" name="codebase_path" placeholder="." required>
                    </div>
                    <button type="submit" class="btn">Detect Drift</button>
                </form>
                <div id="result"></div>
            `;
            document.getElementById('change-order-form').addEventListener('submit', handleChangeOrder);
            break;
    }
}

async function handleIntake(e) {
    e.preventDefault();
    const form = e.target;
    const result = document.getElementById('result');
    result.innerHTML = '<div class="loading">Generating specification</div>';

    const data = {
        project_name: form.project_name.value,
        description: form.description.value,
        target_users: form.target_users.value,
        core_features: form.core_features.value,
        technical_constraints: form.technical_constraints.value,
        success_criteria: form.success_criteria.value
    };

    try {
        const resp = await fetch('/api/intake', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        const json = await resp.json();
        if (json.success) {
            result.innerHTML = `<p class="success">✓ Saved to ${json.path}</p><div class="result">${json.spec}</div>`;
            loadFiles();
        } else {
            result.innerHTML = `<p class="error">Error: ${json.error}</p>`;
        }
    } catch (err) {
        result.innerHTML = `<p class="error">Error: ${err.message}</p>`;
    }
}

async function handleReview(e) {
    e.preventDefault();
    const form = e.target;
    const result = document.getElementById('result');
    result.innerHTML = '<div class="loading">Analyzing code</div>';

    const data = {
        spec_file: form.spec_file.value,
        code_paths: form.code_paths.value.split(',').map(s => s.trim())
    };

    try {
        const resp = await fetch('/api/review', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        const json = await resp.json();
        if (json.success) {
            result.innerHTML = `
                <p class="success">✓ Compliance Score: ${json.compliance_score}/100</p>
                <p>Saved to ${json.path}</p>
                <div class="result">${json.report}</div>
            `;
            loadFiles();
        } else {
            result.innerHTML = `<p class="error">Error: ${json.error}</p>`;
        }
    } catch (err) {
        result.innerHTML = `<p class="error">Error: ${err.message}</p>`;
    }
}

async function handleRescue(e) {
    e.preventDefault();
    const form = e.target;
    const result = document.getElementById('result');
    result.innerHTML = '<div class="loading">Scanning codebase</div>';

    const data = {
        codebase_path: form.codebase_path.value
    };

    try {
        const resp = await fetch('/api/rescue', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        const json = await resp.json();
        if (json.success) {
            result.innerHTML = `
                <p class="success">✓ Scanned ${json.files_scanned} files</p>
                <p>Spec: ${json.spec_path}</p>
                <p>Report: ${json.report_path}</p>
                <div class="result">${json.spec}</div>
            `;
            loadFiles();
        } else {
            result.innerHTML = `<p class="error">Error: ${json.error}</p>`;
        }
    } catch (err) {
        result.innerHTML = `<p class="error">Error: ${err.message}</p>`;
    }
}

async function handleChangeOrder(e) {
    e.preventDefault();
    const form = e.target;
    const result = document.getElementById('result');
    result.innerHTML = '<div class="loading">Detecting drift</div>';

    const data = {
        spec_file: form.spec_file.value,
        codebase_path: form.codebase_path.value
    };

    try {
        const resp = await fetch('/api/change-order', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        const json = await resp.json();
        if (json.success) {
            result.innerHTML = `
                <p class="success">✓ Saved to ${json.path}</p>
                <div class="result">${json.report}</div>
            `;
            loadFiles();
        } else {
            result.innerHTML = `<p class="error">Error: ${json.error}</p>`;
        }
    } catch (err) {
        result.innerHTML = `<p class="error">Error: ${err.message}</p>`;
    }
}
