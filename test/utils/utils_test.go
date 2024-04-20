package utils

import (
	"fmt"
	"github.com/fatih/color"
	"promptrun-api/utils"
	"regexp"
	"testing"
)

func TestGenUUID(t *testing.T) {
	// 测试生成的随机串长度
	uuid := utils.GenUUID()
	fmt.Println(uuid)
	if len(uuid) != 32 {
		t.Errorf("Generated UUID length should be 32, got %d", len(uuid))
	}

	// 测试生成的随机串是否只包含字母和数字
	match, err := regexp.MatchString("^[a-zA-Z0-9]+$", uuid)
	if err != nil {
		t.Errorf("Error matching UUID pattern: %v", err)
	}
	if !match {
		t.Errorf("Generated UUID should only contain letters and numbers")
	}
}

func TestColor(t *testing.T) {
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("Prints cyan text with an underline.")
}

func TestGenSnowflakeId(t *testing.T) {
	// 测试生成的唯一 ID 是否为 64bit
	id := utils.GenSnowflakeId()
	fmt.Println(id)
}
