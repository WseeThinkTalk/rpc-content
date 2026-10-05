package sensitive

import (
	"testing"
)

// TestSensitiveFilter 验证敏感词检测准确性与干扰符号穿透防护
func TestSensitiveFilter(t *testing.T) {
	words := []string{"暴力", "涉黄", "赌博", "badword"}
	f := NewFilter(words)

	cases := []struct {
		input     string
		sensitive bool
	}{
		{"这是一篇正常的技术笔记分享", false},
		{"本文包含涉黄不良内容", true},
		{"涉*黄信息拦截测试", true},
		{"BADWORD test in English", true},
		{"赌   博 网站", true},
	}

	for _, v := range cases {
		got := f.IsSensitive(v.input)
		if got != v.sensitive {
			t.Errorf("测试用例 [%s] 预期敏感判定为 %v, 实际得到 %v", v.input, v.sensitive, got)
		}
	}
}

// TestReplaceSensitive 验证敏感词脱敏替换功能
func TestReplaceSensitive(t *testing.T) {
	words := []string{"违禁词"}
	f := NewFilter(words)
	res := f.Replace("这是包含违禁词的内容", '*')
	expected := "这是包含***的内容"
	if res != expected {
		t.Fatalf("预期脱敏结果为 %s, 实际得到 %s", expected, res)
	}
}
