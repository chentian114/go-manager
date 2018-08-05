package builder

import (
	"errors"
	"strings"
	"sort"
)

/**
 用于构建增、删、改、查sql语句
 */

var (
	errUnsupportedOperator = errors.New(`[builder] unsupported operator`)
	errWhereInType = errors.New(`[builder] the value of "xxx in" must be of []interface{} type`)
	errGroupByValueType = errors.New(`[builder] the value of "_groupby" must be of string type`)

	defaultSortAlgorithm = sort.Strings		//默认排序算法，字符串排序
)

// 空函数
func emptyFunc(){}

//定义比较接口
type Comparable interface{
	Builder([] string,[]interface{})
}

//定义类型 =
type Eq map[string]interface{}
//定义类型 !=
type Ne map[string]interface{}
//定义类型 >
type Gt map[string]interface{}
//定义类型 <
type Lt map[string]interface{}
//定义类型 >=
type GtEq map[string]interface{}
//定义类型 <=
type LtEq map[string]interface{}

func (e Eq) Builder()([]string,[]interface{}){
	return build(e,"=")
}

func build(m map[string]interface{},op string)([]string,[]interface{}){
	if nil == m||  0 == len(m){
		return nil,nil
	}

	len := len(m)
	cond := make([]string,len)  //初始化查询条件集合
	values := make([]interface{},len) //初始化查询条件对应值集合

	var index int = 0
	for key,_ := range m {
		cond[index] = key
		index++
	}
	defaultSortAlgorithm(cond)	//查询条件排序

	for i:=0; i<len ; i++{
		values[index] = m[cond[i]]
		cond[i] = assembleExpression(cond[i],op)	//装配条件表达式
	}
	return cond,values
}

//装配条件表达式
func assembleExpression(field,op string)string{
	return field + op + "?"
}

//构建select语句，所需参数:表名、where语句、查询结果字段列表
//支持的操作符：=,in,>,>=,<,<=,<>,!=.不支持的操作符视为：=
//特殊条件语句
func BuilderSelect(table string,where map[string]interface{},selectFields []string) (cond string,vals []interface{},err error){

	copyWhere := copyWhere(where)	//拷贝where以便对where条件进行操作



	return "",nil,nil
}

//拷贝where
func copyWhere(source map[string]interface{})(copy map[string]interface{})  {
	copy = make(map[string]interface{})
	for key,val := range source {
		copy[key] = val
	}
	return
}




func BuilderInsert(){

}

func BuilderUpdate(){

}

func BuilderDelete(){

}

