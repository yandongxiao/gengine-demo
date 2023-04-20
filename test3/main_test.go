package test3

import (
    "fmt"
    "github.com/bilibili/gengine/builder"
    "github.com/bilibili/gengine/context"
    "github.com/bilibili/gengine/engine"
    "github.com/sirupsen/logrus"
    "testing"
)

// 构建多条规则

type User struct {
}

func (u *User) Print(s string) {
    fmt.Println(s)
}

const rule1 = `
rule "rule 1" salience 0
begin
  User.Print(@name)
end

// 顺序执行模式下，因为 salience 优先级更高，而被优先执行。
rule "rule 2" salience 1
begin
  User.Print(@name)
end
`

func TestMultipleRule(t *testing.T) {
    user := &User{}

    // 1. 构建规则

    // 使用dataContext注入的(变量)数据,对加载到gengine中的所有规则均可见
    dataContext := context.NewDataContext()
    dataContext.Add("User", user)

    // init rule engine
    ruleBuilder := builder.NewRuleBuilder(dataContext)

    err := ruleBuilder.BuildRuleFromString(rule1) // string(bs)
    if err != nil {
        logrus.Errorf("err:%s ", err)
        return
    }

    // 2. 构建引擎
    eng := engine.NewGengine()

    // 3. 执行规则。
    err = eng.Execute(ruleBuilder, true)
    if err != nil {
        logrus.Errorf("execute rule error: %v", err)
    }
}
