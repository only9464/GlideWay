package util

import (
	"fmt"
	"reflect"
)

// 全局函数映射
var functionRegistry = map[string]reflect.Value{}

// RegisterFunction 注册函数到全局映射
func RegisterFunction(funcName string, fn interface{}) {
	functionRegistry[funcName] = reflect.ValueOf(fn)
}

// CallFunction 动态调用已注册的函数，并返回结果和错误
func CallFunction(funcName string, params ...interface{}) (interface{}, error) {
	// log.Printf("函数名为：%s", funcName) // 使用 log.Printf 进行格式化输出
	// for i, param := range params {
	// 	log.Printf("参数 %d: %v", i, param)
	// }
	if fn, ok := functionRegistry[funcName]; ok {
		// log.Printf("函数 %s 已注册", funcName)
		// log.Printf("函数类型: %v", fn.Type())
		// log.Printf("函数参数数量: %d", fn.Type().NumIn())
		// log.Printf("函数返回值数量: %d", fn.Type().NumOut())
		// 确保参数数量匹配
		if len(params) != fn.Type().NumIn() {
			return nil, fmt.Errorf("参数数量不匹配")
		}

		// 准备参数
		in := make([]reflect.Value, len(params))
		for i, param := range params {
			in[i] = reflect.ValueOf(param)
		}

		// 调用函数
		results := fn.Call(in)

		// 处理返回值
		if len(results) > 0 {
			// 检查最后一个返回值是否为错误类型
			lastIndex := len(results) - 1
			if err, ok := results[lastIndex].Interface().(error); ok && err != nil {
				return nil, err
			}

			// 如果只有一个返回值，直接返回
			if len(results) == 1 {
				return results[0].Interface(), nil
			}

			// 如果有多个返回值，返回所有值
			returnValues := make([]interface{}, len(results))
			for i, result := range results {
				returnValues[i] = result.Interface()
			}
			return returnValues, nil
		}
		return nil, fmt.Errorf("没有返回值")
	}
	return nil, fmt.Errorf("函数 %s 未注册", funcName)
}
