package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// テストケースをテーブル駆動テストで定義
	tests := []struct {
		name     string
		args     []string
		wantOut  string
		contains []string // 出力に含まれるべき文字列
	}{
		{
			name:     "引数なしの場合",
			args:     []string{"cmd"},
			contains: []string{"👋 Welcome to your CLI, yagi2!", "No args passed."},
		},
		{
			name:     "引数1つの場合",
			args:     []string{"cmd", "hello"},
			contains: []string{"👋 Welcome to your CLI, yagi2!", "You passed: [hello]"},
		},
		{
			name:     "引数複数の場合",
			args:     []string{"cmd", "hello", "world", "test"},
			contains: []string{"👋 Welcome to your CLI, yagi2!", "You passed: [hello world test]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 標準出力をキャプチャ
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// os.Argsを保存して後で復元
			oldArgs := os.Args
			os.Args = tt.args

			// main関数を実行
			main()

			// 出力を復元
			os.Args = oldArgs
			_ = w.Close()
			os.Stdout = oldStdout

			// キャプチャした出力を読み取り
			out, _ := io.ReadAll(r)
			gotOut := string(out)

			// 期待される文字列が含まれているか確認
			for _, want := range tt.contains {
				if !bytes.Contains(out, []byte(want)) {
					t.Errorf("出力に %q が含まれていません\n実際の出力: %q", want, gotOut)
				}
			}
		})
	}
}

// ベンチマークテストの例
func BenchmarkMain(b *testing.B) {
	// 標準出力を無効化してベンチマーク実行
	oldStdout := os.Stdout
	os.Stdout = nil
	defer func() { os.Stdout = oldStdout }()

	oldArgs := os.Args
	os.Args = []string{"cmd", "bench", "test"}
	defer func() { os.Args = oldArgs }()

	for i := 0; i < b.N; i++ {
		main()
	}
}

// Example関数の例（GoDocに表示される）
func ExampleMain() {
	// os.Argsを設定
	oldArgs := os.Args
	os.Args = []string{"cmd", "example"}
	defer func() { os.Args = oldArgs }()

	main()
	// Output:
	// 👋 Welcome to your CLI, yagi2!
	// You passed: [example]
}
