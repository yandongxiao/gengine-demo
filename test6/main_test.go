package test6

import (
    "fmt"
    "github.com/bilibili/gengine/builder"
    "github.com/bilibili/gengine/context"
    "github.com/bilibili/gengine/engine"
    "github.com/sirupsen/logrus"
    "testing"
)

// 每个Engine构建一个规则，确认两个规则间是否可以互相调用
const rule1 = `
rule "rule 1"
begin
  RunRule2()
end
`

const rule2 = `
rule "rule 2"
begin
  println("i am rule 2")
end
`

var eng *engine.Gengine
var ruleBuilder *builder.RuleBuilder

func init() {
    // 1. 构建规则
    dataContext := context.NewDataContext()
    dataContext.Add("println", fmt.Println)
    ruleBuilder = builder.NewRuleBuilder(dataContext)
    err := ruleBuilder.BuildRuleFromString(rule2) // string(bs)
    if err != nil {
        logrus.Errorf("err:%s ", err)
        return
    }
    // 2. 构建引擎
    eng = engine.NewGengine()
}

func RunRule2() {
    err := eng.Execute(ruleBuilder, true)
    if err != nil {
        logrus.Errorf("execute rule error: %v", err)
    }
}

func TestMultipleRule(t *testing.T) {
    dataContext := context.NewDataContext()
    dataContext.Add("RunRule2", RunRule2)
    ruleBuilder := builder.NewRuleBuilder(dataContext)
    err := ruleBuilder.BuildRuleFromString(rule1) // string(bs)
    if err != nil {
        logrus.Errorf("err:%s ", err)
        return
    }

    eng := engine.NewGengine()
    err = eng.Execute(ruleBuilder, true)
    if err != nil {
        logrus.Errorf("execute rule error: %v", err)
    }

}
