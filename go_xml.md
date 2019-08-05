# golang xml

## xml序列化

```
func Marshal(v interface{}) ([]byte, error)
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
```

###可以使用Marshal函数和MarshalIndent函数返回XML编码。MarshalIndent函数功能类似于Marshal函数。区别在于有无缩进。
### 序列化需要遵守以下规则:
```
- XMLName字段，如上所述，会省略
- 具有标签"-"的字段会省略
- 具有标签"name,attr"的字段会成为该XML元素的名为name的属性
- 具有标签",attr"的字段会成为该XML元素的名为字段名的属性
- 具有标签",chardata"的字段会作为字符数据写入，而非XML元素
- 具有标签",innerxml"的字段会原样写入，而不会经过正常的序列化过程
- 具有标签",comment"的字段作为XML注释写入，而不经过正常的序列化过程，该字段内不能有"--"字符串
- 标签中包含"omitempty"选项的字段如果为空值会省略
 空值为false、0、nil指针、nil接口、长度为0的数组、切片、映射
- 匿名字段（其标签无效）会被处理为其字段是外层结构体的字段
- 如果一个字段的标签为"a>b>c"，则元素c将会嵌套进其上层元素a和b中。如果该字段相邻的字段标签指定了同样的上层元素，则会放在同一个XML元素里。
```

## xml反序列化
```
func Unmarshal(data []byte, v interface{}) error
```
### Unmarshal需要注意以下规则：
```
- 如果结构体字段的类型为字符串或者[]byte，且标签为",innerxml"，
 Unmarshal函数直接将对应原始XML文本写入该字段，其余规则仍适用。
- 如果结构体字段类型为xml.Name且名为XMLName，Unmarshal会将元素名写入该字段
- 如果字段XMLName的标签的格式为"name"或"namespace-URL name"，
 XML元素必须有给定的名字（以及可选的名字空间），否则Unmarshal会返回错误。
- 如果XML元素的属性的名字匹配某个标签",attr"为字段的字段名，或者匹配某个标签为"name,attr"的字段的标签名，Unmarshal会将该属性的值写入该字段。
- 如果XML元素包含字符数据，该数据会存入结构体中第一个具有标签",chardata"的字段中，
 该字段可以是字符串类型或者[]byte类型。如果没有这样的字段，字符数据会丢弃。
- 如果XML元素包含注释，该数据会存入结构体中第一个具有标签",comment"的字段中，
 该字段可以是字符串类型或者[]byte类型。如果没有这样的字段，字符数据会丢弃。
- 如果XML元素包含一个子元素，其名称匹配格式为"a"或"a>b>c"的标签的前缀，反序列化会深入XML结构中寻找具有指定名称的元素，并将最后端的元素映射到该标签所在的结构体字段。
 以">"开始的标签等价于以字段名开始并紧跟着">" 的标签。
- 如果XML元素包含一个子元素，其名称匹配某个结构体类型字段的XMLName字段的标签名，
 且该结构体字段本身没有显式指定标签名，Unmarshal会将该元素映射到该字段。
- 如果XML元素的包含一个子元素，其名称匹配够格结构体字段的字段名，且该字段没有任何模式选项（",attr"、",chardata"等），Unmarshal会将该元素映射到该字段。
- 如果XML元素包含的某个子元素不匹配以上任一条，而存在某个字段其标签为",any"，
 Unmarshal会将该元素映射到该字段。
- 匿名字段被处理为其字段好像位于外层结构体中一样。
- 标签为"-"的结构体字段永不会被反序列化填写。
```