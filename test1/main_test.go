package test2

import (
    "fmt"
    "github.com/bilibili/gengine/builder"
    "github.com/bilibili/gengine/context"
    "github.com/bilibili/gengine/engine"
    "github.com/sirupsen/logrus"
    "testing"
    "time"
)

// User 是要被注入的结构体。但是注重的是API(方法)而非数据。
// gengine 是否可以自动获取从Prometheus获取数据？感觉是可以的，传递一个含有获取Prometheus数据方法的对象。
// gengine 规则执行时间？
type User struct {
    Name string
    Age  int64
    Male bool
}

func (u *User) GetNum(i int64) int64 {
    return i
}

func (u *User) Print(s string) {
    fmt.Println(s)
}

func (u *User) Say() {
    fmt.Println("hello world")
}

// 定义规则
// 规则名称 name test
// 规则描述(非必须) i can
// 规则优先级 salience 数字越大优先级越高。
// 这个规则就是对对象内容进行修改，为什么不直接在代码中指定规则: 可扩展性。
// begin/end中间的语法？支持定义不同类型的变量，支持变量之间的比较操作。
// 什么情况下会并发执行下面的规则？User 对象需要满足并发安全性。
const rule1 = `
rule "name test" "i can"  salience 0
begin
        // 支持单行注释，除了可以调用对象的方法，也可以调用对外暴露的函数
		if 7 == User.GetNum(7){
			User.Age = User.GetNum(89767) + 10000000
			User.Print("6666")
		}else{
			User.Name = "yyyy"
		}
end
`

func TestGetStarted(t *testing.T) {
    user := &User{
        Name: "Calo",
        Age:  0,
        Male: true,
    }

    // 1. 构建规则

    // 使用dataContext注入的(变量)数据,对加载到gengine中的所有规则均可见
    // 注入初始化的结构体，如何替换对象？所以，我们更应该注入的是一种API方法，而不是数据。
    // engine pool 可以在执行规则时，动态注入数据。
    // dataContext.Add("println",fmt.Println) 很棒的一种注入函数的方式。
    dataContext := context.NewDataContext()
    dataContext.Add("User", user)

    // init rule engine
    ruleBuilder := builder.NewRuleBuilder(dataContext)

    // 不停服更新规则的两种方式：全量更新 && 增量更新。
    // https://rencalo770.github.io/gengine_doc/#/compile
    start1 := time.Now()
    err := ruleBuilder.BuildRuleFromString(rule1) // string(bs)
    if err != nil {
        logrus.Errorf("err:%s ", err)
        return
    }
    logrus.Infof("the number of rules: %d, load rules cost time:%v", len(ruleBuilder.Kc.RuleEntities),
        time.Now().Sub(start1),
    )
    for key, rule := range ruleBuilder.Kc.RuleEntities {
        fmt.Println("rule", key, rule.RuleName, rule.RuleDescription, rule.Salience)
    }

    // 2. 构建引擎
    eng := engine.NewGengine()

    // 3. 执行规则
    // Execute 是顺序执行模式。
    // ExecuteConcurrent 是并发执行模式。
    // gengine 支持了非常丰富的执行模式，详情参见 https://rencalo770.github.io/gengine_doc
    start := time.Now().UnixNano()
    err = eng.Execute(ruleBuilder, true)
    if err != nil {
        logrus.Errorf("execute rule error: %v", err)
    }
    fmt.Println(user.Age)
    end := time.Now().UnixNano()
    logrus.Infof("execute rule cost %d ns", end-start)

    // 4. 检查获取结果
    logrus.Infof("user.Age=%d,Name=%s,Male=%t", user.Age, user.Name, user.Male)
}
