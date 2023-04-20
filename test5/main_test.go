package test5

import (
    "fmt"
    "github.com/bilibili/gengine/engine"
    "testing"
)

// 该测试主要介绍什么是规则引擎池（gengine pool）

// 为什么需要规则引擎池：加载到一个gengine实例中的一系列规则,在执行的时候,是有状态的。
// 当一个gengine实例作为服务核心中调用的的模块的时候,当规则少的时候,一个请求数据执行所
// 有规则的时间非常短,基于当前请求执行的规则的状态维持的时间也非常短,因此不会引发任何问题;
// 但当一个gengine实例中加载了几十个甚至上百个规则的时候,一个请求执行完所有的规则的时间就会变长,尤其是处于
// 高QPS的情况下,当前请求还未执行完并返回时,下一个请求就已经到来,这个时候下一个请求就极有可能破坏当前请求执
// 行规则的状态,并最终导致不可预知的结果
// 为了解决这个问题,gengine框架仿照”数据库连接池”,实现了一个”gengine池”.
// pool的所有API都是线程安全的

const pool_mix_model_rule = `
rule "best" "best"  salience 100
begin
	println("best....")
	Ps.P = true
	println("best....", Ps.P)
end

rule "better" "better"   salience 99
begin

if Ps.P {
	println("better....")

	Ps.R = true
	println("better....",Ps.R)
}
end


rule "good" "good"   salience 98
begin
	println("good....")
	println("good....",Ps.R, Ps.P)
end
`

type Ps struct {
    P bool
    R bool
}

func Test_mix_model(t *testing.T) {

    apis := make(map[string]interface{})

    // poolMinLen 池中初始实例化的gengine实例个数
    // poolMaxLen 池子中最多可实例化的gengine实例个数, 且poolMaxLen > poolMinLen; 当poolMinLen个实例不够用
    //            的时候,最多还可实例化(poolMaxLen-poolMinLen)个gengine实例
    // em 规则执行模式,em只能等于1、2、3、 4; em=1时,表示顺序执行规则,em=2的时候,表示并发执行规则,em=3的时候,表示
    //    混合模式执行规则,em=4的时候,表示逆混合模式执行规则;当使用ExecuteSelectedRules和
    //    ExecuteSelectedRulesConcurrent等指明了具体的执行模式的方法执行规则时,此参数自动失效
    // rulesStr要初始化的所有的规则字符串
    // apiOuter需要注入到gengine中使用的所有api,最好仅注入一些公共的无状态函数或参数;对于那些具体与某次请求(执行)相
    // 关的参数,则在执行规则方法时使用data map[string]interface{} 注入;这样会有利于状态管理。
    pool, e1 := engine.NewGenginePool(1, 3, 3, pool_mix_model_rule, apis)
    if e1 != nil {
        panic(e1)
    }

    println("pool.GetExecModel()==", pool.GetExecModel())
    data := make(map[string]interface{})
    Ps := &Ps{}
    data["Ps"] = Ps
    data["println"] = fmt.Println

    e, _ := pool.Execute(data, true)
    if e != nil {
        panic(e)
    }
}
