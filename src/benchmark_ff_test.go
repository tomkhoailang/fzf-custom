package fzf

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/junegunn/fzf/src/algo"
	"github.com/junegunn/fzf/src/util"
)

func BenchmarkFF(b *testing.B) {
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = "/home/hj/coding/ruby/gitlab"
	out, err := cmd.Output()
	if err != nil {
		b.Fatalf("Failed to run git ls-files: %v", err)
	}

	lines := strings.Split(string(out), "\n")
	var formattedItems []*Item

	mruList := []string{
		"app/controllers/projects/issues_controller.rb",
		"jest.config.snapshots.js",
	}

	algo.Init("filename-first")
	algo.MruMap = make(map[string]int)
	for i, m := range mruList {
		algo.MruMap[m] = i + 1
	}

	for _, file := range lines {
		if file == "" {
			continue
		}
		dir := filepath.Dir(file)
		base := filepath.Base(file)

		var formatted string
		if dir == "." {
			formatted = fmt.Sprintf("  %s", base)
		} else {
			formatted = fmt.Sprintf("  %s  \x1b[38;5;244;3m%s\x1b[0m", base, dir)
		}

		chars := util.ToChars([]byte(formatted))
		item := &Item{text: chars}
		formattedItems = append(formattedItems, item)
	}

	// Pre-populate cache (simulates second+ keystroke, cache already warm)
	for _, item := range formattedItems {
		algo.GetOrInitFfCache(&item.text)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pattern := &Pattern{
			fuzzy:         true,
			extended:      true,
			caseSensitive: true,
			normalize:     true,
			forward:       true,
			withPos:       true,
			text:          []rune("iss con"),
			termSets:      parseTerms(true, CaseSmart, true, "iss con"),
		}
		pattern.procFun[termFuzzy] = algo.FuzzyMatchV2FilenameFirstNoNeural

		slab := util.MakeSlab(16*1024, 32*1024)

		for _, item := range formattedItems {
			pattern.MatchItem(item, true, slab)
		}
	}
}
