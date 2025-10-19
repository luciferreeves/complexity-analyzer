let editor;
let currentLanguage = 'javascript';
let chart = null;
let loadingInterval = null;
let loadingMessageIndex = 0;

const loadingMessages = [
    "Parsing code structure...",
    "Analyzing algorithm patterns...",
    "Detecting time complexity...",
    "Identifying nested loops...",
    "Checking for recursion...",
    "Examining data structures...",
    "Loading execution environment...",
    "Compiling test cases...",
    "Generating test data...",
    "Running performance benchmarks...",
    "Measuring execution time...",
    "Testing with small inputs...",
    "Testing with large inputs...",
    "Analyzing growth patterns...",
    "Calculating complexity metrics...",
    "Fitting complexity curves...",
    "Validating results...",
    "Finalizing analysis..."
];

const irregularIntervals = [2400, 3800, 4500, 1000, 2900, 4300, 1100, 3400, 950, 2250, 850, 3200, 2050, 1200, 1400, 1800, 900, 1300];

require.config({
    paths: {
        vs: 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.44.0/min/vs'
    }
});

require(['vs/editor/editor.main'], function () {
    editor = monaco.editor.create(document.getElementById('editor'), {
        value: '// Write your algorithm here\nfunction algorithm(arr) {\n    // Your code\n    return arr;\n}',
        language: 'javascript',
        theme: 'vs-dark',
        fontSize: 14,
        minimap: { enabled: false },
        scrollBeyondLastLine: false,
        automaticLayout: true,
        lineNumbers: 'on',
        roundedSelection: false,
        scrollbar: {
            vertical: 'visible',
            horizontal: 'visible'
        }
    });

    setupEventListeners();
});

function setupEventListeners() {
    // Language selector
    document.getElementById('languageSelect').addEventListener('change', (e) => {
        currentLanguage = e.target.value;
        monaco.editor.setModelLanguage(editor.getModel(), currentLanguage);

        // default code template
        const templates = {
            javascript: '// Write your algorithm here\nfunction algorithm(arr) {\n    // Your code\n    return arr;\n}',
            python: '# Write your algorithm here\ndef algorithm(arr):\n    # Your code\n    return arr',
            go: '// Write your algorithm here\nfunc algorithm(arr []int) []int {\n    // Your code\n    return arr\n}'
        };
        editor.setValue(templates[currentLanguage]);
    });

    // Analyze button
    document.getElementById('analyzeBtn').addEventListener('click', analyzeCode);
}

function startDynamicLoading() {
    loadingMessageIndex = 0;
    const loadingText = document.querySelector('#loadingState p');

    function updateMessage() {
        if (loadingMessageIndex < loadingMessages.length) {
            loadingText.textContent = loadingMessages[loadingMessageIndex];

            if (loadingMessageIndex < loadingMessages.length - 1) {
                const nextInterval = irregularIntervals[loadingMessageIndex];
                loadingMessageIndex++;
                loadingInterval = setTimeout(updateMessage, nextInterval);
            } else {
                loadingMessageIndex++;
            }
        }
    }

    updateMessage();
}

function stopDynamicLoading() {
    if (loadingInterval) {
        clearTimeout(loadingInterval);
        loadingInterval = null;
    }
}

async function analyzeCode() {
    const code = editor.getValue();

    if (!code.trim()) {
        alert('Please enter code to analyze');
        return;
    }

    showLoading();
    startDynamicLoading();

    const analyzeBtn = document.getElementById('analyzeBtn');
    analyzeBtn.disabled = true;

    try {
        const response = await fetch('/api/analyze', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                code: code,
                language: currentLanguage
            })
        });

        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.error || 'Analysis failed');
        }

        displayResults(result);
    } catch (error) {
        displayError(error.message);
    } finally {
        stopDynamicLoading();
        analyzeBtn.disabled = false;
    }
}

function showLoading() {
    document.getElementById('resultsContent').style.display = 'none';
    document.getElementById('loadingState').classList.add('active');
}

function displayResults(result) {
    stopDynamicLoading();
    document.getElementById('loadingState').classList.remove('active');
    document.getElementById('resultsContent').style.display = 'block';

    const resultsHTML = `
        <div class="result-section">
            <h3>Detected Complexity</h3>
            <div class="complexity-badge">${result.complexity}</div>
            <p style="margin-top: 8px; font-size: 12px; color: #858585;">
                Confidence: ${result.confidence.toFixed(1)}%
            </p>
        </div>
        
        <div class="result-section">
            <h3>Static Analysis</h3>
            <ul class="info-list">
                ${result.staticAnalysis.map(item => `<li>${item}</li>`).join('')}
            </ul>
        </div>
        
        <div class="result-section">
            <h3>Performance Graph</h3>
            <div class="chart-container">
                <canvas id="performanceChart"></canvas>
            </div>
        </div>
        
        <div class="result-section">
            <h3>Execution Timings</h3>
            <ul class="info-list">
                ${result.performanceData.map(t =>
        `<li>n = ${t.size}: ${t.time.toFixed(3)}ms</li>`
    ).join('')}
            </ul>
        </div>
    `;

    document.getElementById('resultsContent').innerHTML = resultsHTML;

    renderChart(result.performanceData, result.complexity);
}

function renderChart(data, complexity) {
    if (typeof Chart === 'undefined') {
        console.error('Chart.js is not available');
        return;
    }

    const ctx = document.getElementById('performanceChart');

    // Destroy previous chart if exists
    if (chart) {
        chart.destroy();
    }

    chart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: data.map(d => d.size),
            datasets: [{
                label: `Execution Time (${complexity})`,
                data: data.map(d => d.time),
                borderColor: '#667eea',
                backgroundColor: 'rgba(102, 126, 234, 0.1)',
                borderWidth: 2,
                pointRadius: 4,
                pointBackgroundColor: '#667eea',
                tension: 0.4
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    labels: { color: '#d4d4d4' }
                },
                tooltip: {
                    backgroundColor: '#252526',
                    titleColor: '#fff',
                    bodyColor: '#d4d4d4',
                    borderColor: '#3e3e42',
                    borderWidth: 1
                }
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: 'Input Size (n)',
                        color: '#d4d4d4'
                    },
                    ticks: { color: '#858585' },
                    grid: { color: '#3e3e42' }
                },
                y: {
                    title: {
                        display: true,
                        text: 'Time (ms)',
                        color: '#d4d4d4'
                    },
                    ticks: { color: '#858585' },
                    grid: { color: '#3e3e42' }
                }
            }
        }
    });
}

function displayError(message) {
    stopDynamicLoading();
    document.getElementById('loadingState').classList.remove('active');
    document.getElementById('resultsContent').style.display = 'block';
    document.getElementById('resultsContent').innerHTML = `
        <div class="result-section">
            <h3 style="color: #f48771;">Error</h3>
            <div style="background: #1e1e1e; padding: 12px; border-radius: 4px; color: #f48771;">
                ${message}
            </div>
        </div>
    `;
}

// Initialize app when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initApp);
} else {
    initApp();
}

function initApp() {
    if (typeof Chart === 'undefined') {
        console.error('Chart.js failed to load');
        alert('Failed to load required libraries. Please refresh the page.');
    }
}