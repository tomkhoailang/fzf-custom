package algo

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Network struct {
	Weights      [][][]float32 `json:"weights"`
	Biases       [][]float32   `json:"biases"`
	Gammas       [][]float32   `json:"gammas"`
	Betas        [][]float32   `json:"betas"`
	RunningMeans [][]float32   `json:"running_means"`
	RunningVars  [][]float32   `json:"running_vars"`
}

type NeuralOpenWeights struct {
	NN struct {
		Network Network `json:"network"`
	} `json:"nn"`
}

var (
	NeuralNet              *Network
	NeuralWeightsFile      string
	CurrentFile            string
	CurrentFileDir         string
	CurrentFileDepth       int
	CurrentFileTrigrams    map[int]bool
	CurrentFileTrigramsLen int
	OpenBuffers            map[string]bool
	AltBuffer              string
	MruListSize            int = 100
	TransitionScores       map[string]float32
	FrecencyScores         map[string]float32
	NeuralInitOnce         sync.Once
)

func InitNeural() {
	NeuralInitOnce.Do(func() {
		NeuralWeightsFile = os.Getenv("FZF_NEURAL_WEIGHTS_FILE")
		if NeuralWeightsFile == "" {
			return
		}

		// Load neural weights
		if err := LoadNeuralWeights(NeuralWeightsFile); err != nil {
			return
		}

		CurrentFile = normalizePath(os.Getenv("FZF_CURRENT_FILE"))
		if CurrentFile != "" {
			dir := filepath.Dir(CurrentFile)
			if dir != "." {
				CurrentFileDir = dir + "/"
				CurrentFileDepth = countSlashes(CurrentFileDir)
			}
			CurrentFileTrigrams = computeTrigrams(filepath.Base(CurrentFile))
			CurrentFileTrigramsLen = len(CurrentFileTrigrams)
		}

		AltBuffer = normalizePath(os.Getenv("FZF_ALT_BUFFER"))

		// Open buffers
		OpenBuffers = make(map[string]bool)
		for _, b := range strings.Split(os.Getenv("FZF_OPEN_BUFFERS"), ";") {
			if b != "" {
				OpenBuffers[normalizePath(b)] = true
			}
		}

		// MRU Map is already handled by algo.MruMap

		// Transitions
		TransitionScores = make(map[string]float32)
		transEnv := os.Getenv("FZF_TRANSITIONS")
		if transEnv != "" {
			for _, part := range strings.Split(transEnv, ";") {
				if part != "" {
					kv := strings.SplitN(part, ":", 2)
					if len(kv) == 2 {
						if score, err := strconv.ParseFloat(kv[1], 32); err == nil {
							TransitionScores[normalizePath(kv[0])] = float32(score)
						}
					}
				}
			}
		}

		// Frecency
		FrecencyScores = make(map[string]float32)
		frecEnv := os.Getenv("FZF_FRECENCY")
		if frecEnv != "" {
			for _, part := range strings.Split(frecEnv, ";") {
				if part != "" {
					kv := strings.SplitN(part, ":", 2)
					if len(kv) == 2 {
						if score, err := strconv.ParseFloat(kv[1], 32); err == nil {
							FrecencyScores[normalizePath(kv[0])] = float32(score)
						}
					}
				}
			}
		}
	})
}

func LoadNeuralWeights(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var data NeuralOpenWeights
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	NeuralNet = &data.NN.Network
	return nil
}

func normalizePath(p string) string {
	if p == "" {
		return ""
	}
	p = strings.ReplaceAll(p, "\\", "/")
	p = strings.TrimSpace(p)
	return p
}

func countSlashes(s string) int {
	count := 0
	start := 0
	if len(s) > 0 && s[0] == '/' {
		start = 1
	}
	for i := start; i < len(s); i++ {
		if s[i] == '/' {
			count++
		}
	}
	return count
}

func computeTrigrams(text string) map[int]bool {
	tris := make(map[int]bool)
	text = strings.ToLower(text)
	runes := []rune(text)
	if len(runes) < 3 {
		return tris
	}
	for i := 0; i < len(runes)-2; i++ {
		b1 := int(runes[i])
		b2 := int(runes[i+1])
		b3 := int(runes[i+2])
		key := b1*65536 + b2*256 + b3
		tris[key] = true
	}
	return tris
}

func getVirtualName(path string) string {
	filename := filepath.Base(path)
	isSpecial := false
	switch filename {
	case "index.js", "index.jsx", "index.ts", "index.tsx", "init.lua", "init.vim", "__init__.py":
		isSpecial = true
	}

	if isSpecial {
		parent := filepath.Base(filepath.Dir(path))
		if parent != "" && parent != "." && parent != "/" {
			return parent + "/" + filename
		}
	}
	return filename
}

func calculateProximity(targetPath string) float32 {
	if CurrentFileDir == "" || targetPath == "" {
		return 0
	}
	targetDir := filepath.Dir(targetPath)
	if targetDir == "." {
		return 0
	}
	targetDir = targetDir + "/"

	currentDirLen := len(CurrentFileDir)
	targetDirLen := len(targetDir)
	scanLen := currentDirLen
	if targetDirLen < scanLen {
		scanLen = targetDirLen
	}

	commonDepth := 0
	start := 0
	if currentDirLen > 0 && CurrentFileDir[0] == '/' {
		start = 1
	}

	for i := start; i < scanLen; i++ {
		if CurrentFileDir[i] != targetDir[i] {
			break
		}
		if CurrentFileDir[i] == '/' {
			commonDepth++
		}
	}

	if commonDepth == 0 {
		return 0
	}

	targetDepth := countSlashes(targetDir)
	maxDepth := CurrentFileDepth
	if targetDepth > maxDepth {
		maxDepth = targetDepth
	}

	if maxDepth == 0 {
		return 0
	}
	return float32(commonDepth) / float32(maxDepth)
}

func calculateTrigramSimilarity(virtualName string) float32 {
	if CurrentFileTrigramsLen == 0 || len(virtualName) < 3 {
		return 0
	}

	trigrams2 := computeTrigrams(virtualName)
	if len(trigrams2) == 0 {
		return 0
	}

	intersection := 0
	for k := range trigrams2 {
		if CurrentFileTrigrams[k] {
			intersection++
		}
	}

	return float32(2*intersection) / float32(CurrentFileTrigramsLen+len(trigrams2))
}

func (nn *Network) Forward(input []float32) float32 {
	if len(nn.Weights) == 0 {
		return 0
	}

	current := input
	for i := 0; i < len(nn.Weights); i++ {
		w := nn.Weights[i]
		b := nn.Biases[i]
		
		inputDim := len(w)
		outputDim := len(w[0])
		
		next := make([]float32, outputDim)
		for col := 0; col < outputDim; col++ {
			var val float32 = 0
			for row := 0; row < inputDim; row++ {
				if row < len(current) {
					val += current[row] * w[row][col]
				}
			}
			if col < len(b) {
				val += b[col]
			}
			next[col] = val
		}

		// Batch Normalization (hidden layers only)
		if i < len(nn.Weights)-1 && i < len(nn.Gammas) && i < len(nn.Betas) && i < len(nn.RunningMeans) && i < len(nn.RunningVars) {
			gamma := nn.Gammas[i]
			beta := nn.Betas[i]
			mean := nn.RunningMeans[i]
			variance := nn.RunningVars[i]
			
			for col := 0; col < outputDim; col++ {
				if col < len(gamma) && col < len(beta) && col < len(mean) && col < len(variance) {
					v := next[col]
					normalized := (v - mean[col]) / float32(math.Sqrt(float64(variance[col]+1e-8)))
					next[col] = normalized*gamma[col] + beta[col]
				}
			}
		}

		// Activation
		if i < len(nn.Weights)-1 {
			for col := 0; col < outputDim; col++ {
				v := next[col]
				if v < 0 {
					next[col] = v * 0.01
				}
			}
		} else {
			v := next[0]
			next[0] = 1.0 / (1.0 + float32(math.Exp(float64(-v))))
		}

		current = next
	}

	return current[0]
}

func CalculateNeuralScore(path string, matchScore float32, virtualMatchScore float32) float32 {
	InitNeural()
	if NeuralNet == nil {
		return 0
	}

	normPath := normalizePath(path)
	input := make([]float32, 11)

	// 1. match
	input[0] = matchScore / 5000.0
	if input[0] > 1.0 {
		input[0] = 1.0
	}

	// 2. virtual_name match
	input[1] = virtualMatchScore / 5000.0
	if input[1] > 1.0 {
		input[1] = 1.0
	}

	// 3. frecency
	input[2] = FrecencyScores[normPath]

	// 4. open buffer
	if OpenBuffers[normPath] {
		input[3] = 1.0
	} else {
		input[3] = 0.0
	}

	// 5. alt buffer
	if AltBuffer != "" && normPath == AltBuffer {
		input[4] = 1.0
	} else {
		input[4] = 0.0
	}

	// 6. proximity
	input[5] = calculateProximity(normPath)

	// 7. project
	cwd := os.Getenv("FZF_PROJECT_CWD")
	if cwd != "" && strings.HasPrefix(normPath, normalizePath(cwd)) {
		input[6] = 1.0
	} else {
		input[6] = 0.0
	}

	// 8. recency
	if rank, ok := MruMap[normPath]; ok {
		factor := float32(MruListSize-rank+1) / float32(MruListSize)
		if factor < 0 {
			factor = 0
		}
		if factor > 1.0 {
			factor = 1.0
		}
		input[7] = factor
	} else {
		input[7] = 0.0
	}

	// 9. trigram similarity
	vname := getVirtualName(normPath)
	input[8] = calculateTrigramSimilarity(vname)

	// 10. transition
	input[9] = TransitionScores[normPath]

	// 11. not_current
	if CurrentFile != "" && normPath == CurrentFile {
		input[10] = 0.0
	} else {
		input[10] = 1.0
	}

	return NeuralNet.Forward(input)
}
