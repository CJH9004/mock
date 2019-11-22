# mock

Mock data according to tags

## 支持的数据类型

int, uint, float, string, slice, array, map, struct

## 示例

```go
type Struct struct {
  Field string `mock: "type(sentence) range(10)"`
}

val tag = "type(sentence) range(10)"
```

## 标签函数

### type

- 支持的参数：eamil, date, phone, url, ipv4, domain, word, sentence
- date支持string和int64，其它类型仅支持string

### range

- 默认[1,10)
- range(n): [1, n) or [n, 1)
- range(min, max): [min, max)

### value

- value(v1, v2, v3, ...): [v1, v2, v3]中随机取值

### mock

- 自定义mock函数名

### key

- 为map类型的key指定tag

### elem

- 为map, slice, array类型的元素指定tag

### format

- 为date类型指定格式

### tag

- 为当前field指定tag

## 详细使用请查看mock_test.go

## Todos

- Valid