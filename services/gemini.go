package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AnalysisResult struct {
	Complexity      string             `json:"complexity"`
	Confidence      float64            `json:"confidence"`
	StaticAnalysis  []string           `json:"staticAnalysis"`
	TestCode        string             `json:"testCode"`
	PerformanceData []PerformancePoint `json:"performanceData"`
}

type PerformancePoint struct {
	Size int     `json:"size"`
	Time float64 `json:"time"`
}

func AnalyzeWithGemini(code string, language string) (*AnalysisResult, error) {
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SetTemperature(0.1)

	prompt := buildAnalysisPrompt(code, language)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	var result AnalysisResult
	responseText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	responseText = extractJSON(responseText)

	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	perfData, err := runPerformanceTest(result.TestCode)
	if err != nil {
		log.Printf("ERROR: %v", err)
		result.PerformanceData = []PerformancePoint{}
	} else {
		result.PerformanceData = perfData
	}

	return &result, nil
}

func buildAnalysisPrompt(code string, language string) string {
	return fmt.Sprintf(`You are an expert Big O complexity analyzer. Analyze this %s code.

CODE:
%s

Return ONLY valid JSON (no markdown, no code blocks):
{
  "complexity": "O(n·n!)",
  "confidence": 95.0,
  "staticAnalysis": [
    "Specific detail about what this algorithm does",
    "Another specific observation", 
    "A third unique detail"
  ],
  "testCode": "package main\n\nimport (\"fmt\"; \"time\")\n\nfunc algorithm() {}\n\nfunc main() {\n\tfmt.Println(\"1000,0.123\")\n}"
}

COMPLEXITY NOTATION:
- Use · for multiplication: O(n·2ⁿ), O(n·n!), O(n²·log n)
- Common: O(1), O(log n), O(√n), O(n), O(n log n), O(n²), O(n³)
- Exponential: O(2ⁿ), O(3ⁿ), O(n·2ⁿ)
- Factorial: O(n!), O(n·n!)
- When copying/building each result: multiply by result size (e.g., building n items each of size n → O(n²))

STATIC ANALYSIS - Describe WHAT this specific code does:
❌ Generic: "Has a loop", "Uses recursion", "Sorted array required"
✅ Specific: "Halves search space by comparing middle element with target each iteration"
✅ Specific: "Recursively generates all arrangements by trying each element in first position"
✅ Specific: "Maintains left/right pointers that converge while tracking maximum heights"

TEST CODE - Generate COMPLETE runnable Go program:

1. Implement the algorithm with appropriate signature (analyze the code to determine params)

2. Choose test sizes based on expected complexity:
   - Fast (O(1) to O(n log n)): [1000, 5000, 10000, 50000, 100000]
   - Quadratic (O(n²)): [100, 500, 1000, 2000, 3000]
   - Cubic (O(n³)): [20, 40, 60, 80, 100]
   - Exponential (O(2ⁿ)): [10, 12, 14, 16, 18, 20]
   - Factorial (O(n!)): [7, 8, 9, 10]

3. For EACH size, run algorithm repeatedly until ~500ms elapsed:
   - Start iterations = 1
   - Run algorithm iterations times, measure total time
   - If time < 500ms, double iterations and repeat
   - Once ≥500ms, calculate average time per iteration
   - Cap iterations at 10,000,000 to prevent infinite loops

4. Output format: size,time_in_ms (one per line)
   Example: 1000,0.000123

5. Example template:
package main
import ("fmt"; "time")

func algorithm(/* your params */) /* return type */ {
    /* implementation */
}

func main() {
    sizes := []int{/* appropriate sizes */}
    targetTime := 500.0 // milliseconds
    
    for _, size := range sizes {
        // Generate test data appropriate for this algorithm
        
        iterations := 1
        for {
            start := time.Now()
            for i := 0; i < iterations; i++ {
                algorithm(/* call with test data */)
            }
            elapsed := time.Since(start).Milliseconds()
            
            if float64(elapsed) >= targetTime {
                avgMs := float64(elapsed) / float64(iterations)
                fmt.Printf("%%d,%%.9f\n", size, avgMs)
                break
            }
            
            iterations *= 2
            if iterations > 10000000 {
                avgMs := float64(elapsed) / float64(iterations/2)
                fmt.Printf("%%d,%%.9f\n", size, avgMs)
                break
            }
        }
    }
}

Return ONLY JSON.`, language, code)
}

func extractJSON(text string) string {
	if strings.Contains(text, "```") {
		start := strings.Index(text, "\n")
		if start != -1 {
			text = text[start+1:]
		}
		end := strings.LastIndex(text, "```")
		if end != -1 {
			text = text[:end]
		}
	}
	return strings.TrimSpace(text)
}

func runPerformanceTest(testCode string) ([]PerformancePoint, error) {
	tmpDir, err := os.MkdirTemp("", "complexity-test-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := tmpDir + "/main.go"
	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write test file: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", testFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("failed to run test: %w\nOutput: %s", err, string(output))
	}

	var perfData []PerformancePoint
	lines := strings.SplitSeq(string(output), "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var size int
		var timeMs float64

		if _, err := fmt.Sscanf(line, "%d,%f", &size, &timeMs); err == nil {
			perfData = append(perfData, PerformancePoint{
				Size: size,
				Time: timeMs,
			})
		}
	}

	if len(perfData) == 0 {
		return nil, fmt.Errorf("no performance data parsed")
	}

	return perfData, nil
}
