package test2

import (
    "github.com/bilibili/gengine/builder"
    "github.com/bilibili/gengine/context"
    "github.com/bilibili/gengine/engine"
    "github.com/sirupsen/logrus"
    "testing"
)

// 该测试用例希望将 User 对象注入引擎，通过调用它的 API 获取用户所有的订单信息，如果
// 如果用户的订单总金额大于等于 100，则返回 true，否则返回 false。

type Order struct {
    Id    int
    Money int
}

func (order *Order) CopyFrom(newOrder *Order) {
    order.Id = newOrder.Id
    order.Money = newOrder.Money
}

type User struct {
    Name   string
    Orders []*Order
}

func (u *User) LenOfOrders() int {
    return len(u.Orders)
}

func (u *User) GetOrders() []*Order {
    return []*Order{
        {
            Id:    1,
            Money: 20,
        },
        {
            Id:    2,
            Money: 80,
        },
    }
}

const rule = `
rule "check order" salience 0
begin
    orders := User.GetOrders()
    User.Orders = orders

    total = 0
    for i := 0; i<User.LenOfOrders(); i=i+1 {
        Order.CopyFrom(User.Orders[i])
        total = total + Order.Money
    }
    if total >= 100 {
        return true
    }
    return false
end
`

func TestCheckOrder(t *testing.T) {
    user := &User{
        Name: "Calo",
    }

    // 1. 构建规则
    dataContext := context.NewDataContext()
    dataContext.Add("User", user)
    dataContext.Add("Order", &Order{})

    // init rule engine
    ruleBuilder := builder.NewRuleBuilder(dataContext)

    err := ruleBuilder.BuildRuleFromString(rule)
    if err != nil {
        logrus.Errorf("err:%s ", err)
        return
    }

    // 2. 构建引擎
    eng := engine.NewGengine()

    // 3. 执行规则
    err = eng.Execute(ruleBuilder, true)
    if err != nil {
        logrus.Errorf("execute rule error: %v", err)
    }

    // 4. 检查获取结果
    results, _ := eng.GetRulesResultMap()
    for i := 0; i < len(user.Orders); i++ {
        logrus.Infoln(user.Orders[i], results["check order"])
    }
}
