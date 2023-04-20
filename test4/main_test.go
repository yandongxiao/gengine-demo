package test4

import (
    "fmt"
    "github.com/bilibili/gengine/builder"
    "github.com/bilibili/gengine/context"
    "github.com/bilibili/gengine/engine"
    "github.com/sirupsen/logrus"
    "testing"
    "time"
)

const rule_conc_statement = `
rule "conc_test" "test" 
begin
	conc  {
        // 下面每条语句的实际执行顺序是并发的。与下面的书写顺序不同
		println("AAA")
		println("BBB")
		println("CCC")
		println("DDD")
		println("EEE")
	}
end
`

func Sout(str string) {
    println("----", str)
}

func Test_conc_statement(t *testing.T) {

    dataContext := context.NewDataContext()
    dataContext.Add("println", logrus.Println)
    dataContext.Add("sout", Sout)

    // init rule engine
    ruleBuilder := builder.NewRuleBuilder(dataContext)

    // resolve rules from string
    start1 := time.Now().UnixNano()
    err := ruleBuilder.BuildRuleFromString(rule_conc_statement)
    end1 := time.Now().UnixNano()

    println(fmt.Sprintf("rules num:%d, load rules cost time:%d ns", len(ruleBuilder.Kc.RuleEntities), end1-start1))

    if err != nil {
        panic(err)
    }
    eng := engine.NewGengine()
    start := time.Now().UnixNano()
    // true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
    err = eng.Execute(ruleBuilder, true)
    end := time.Now().UnixNano()
    if err != nil {
        panic(err)
    }
    println(fmt.Sprintf("execute rule cost %d ns", end-start))

}
