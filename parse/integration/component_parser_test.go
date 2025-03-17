package integration

import (
	"testing"

	"github.com/zhangshanwen/html2go/parse"
)

func TestParseComponentData(t *testing.T) {
	// 解析组件数据
	components, err := parse.ParseComponentData()
	if err != nil {
		t.Fatalf("解析组件数据失败: %v", err)
	}

	// 验证数据不为空
	if len(components) == 0 {
		t.Fatal("解析的组件数据为空")
	}

	// 测试几个特定组件是否存在并验证其结构
	// 测试vuetifyx中的组件
	vxBtn, exists := components["vx-btn"]
	if !exists {
		t.Fatal("找不到组件vx-btn")
	}

	// 验证vx-btn组件的属性
	expected := parse.ComponentDefinition{
		Go:     "VXBtn",
		Accept: "none",
	}

	// 只比较Go和Accept字段
	if vxBtn.Go != expected.Go || vxBtn.Accept != expected.Accept {
		t.Errorf("vx-btn组件错误: 期望 %+v, 得到 %+v", expected, vxBtn)
	}

	// 验证vx-btn组件的某些属性是否存在
	attrs := []string{"color", "disabled", "text"}
	for _, attr := range attrs {
		if _, exists := vxBtn.Attrs[attr]; !exists {
			t.Errorf("vx-btn组件缺少属性: %s", attr)
		}
	}

	// 测试vuetify中的组件
	vAlert, exists := components["v-alert"]
	if !exists {
		t.Fatal("找不到组件v-alert")
	}

	// 验证v-alert组件的属性
	expectedAlert := parse.ComponentDefinition{
		Go:     "VAlert",
		Accept: "...h.HTMLComponent",
	}

	// 只比较Go和Accept字段
	if vAlert.Go != expectedAlert.Go || vAlert.Accept != expectedAlert.Accept {
		t.Errorf("v-alert组件错误: 期望 %+v, 得到 %+v", expectedAlert, vAlert)
	}

	// 验证组件数量至少有两个文件中的组件数量之和
	minCount := 10 // 保守估计两个文件中有多少组件
	if len(components) < minCount {
		t.Errorf("组件数量过少: 期望至少 %d, 得到 %d", minCount, len(components))
	}

	// 遍历几个组件，并打印部分信息作为参考
	t.Log("解析的部分组件信息:")
	count := 0
	for name, comp := range components {
		if count >= 5 {
			break
		}
		t.Logf("组件: %s, Go类型: %s, Accept: %s, 属性数量: %d",
			name, comp.Go, comp.Accept, len(comp.Attrs))
		count++
	}
}
