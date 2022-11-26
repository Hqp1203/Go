# Golang知识点  

## 一. 入门篇  

### 1. 配置  

#### 1)配置GOROOT环境变量
#### 2）Go设置
```powershell
$env:GO111MODULE = "on"
$env:GOPROXY = "http://goproxy.cn"
```


### 2. 命令行运行Go程序

#### a. 编译加运行（产生可执行文件）

```sh
# 编译go源码
go build -o [myApp.exe] main.go
# 运行可执行文件
./myApp.exe
```

#### b.一次性运行（直接执行代码）

```sh
go run main.go
```

## 二. 基础语法  

### 基本数据类型和字符串的转化

#### a、基本数据类型转字符串

##### 1）fmt.Sprintf()
函数原型

```go
func Sprintf(format string,a ...interface{})string
```
代码示例
```go
func main() {
	var num1 int = 12
	var num2 float32 = 3.24
	var num3 bool = true
	var num4 byte = 'a'
	var str string = "i am sdad"
	str = fmt.Sprintf("%X\t%f\t%t\t%c", num1, num2, num3, num4)
    # %q 带``的字符串,不会转移特殊符号
    fmt.Printf("%q",str)
}
```

##### 2)strconv.Xxx()

函数原型
```go
func Itoa(i int)string
func FormatBool(b bool)string
func FormatInt(i int64, base int)string
func FormatUint(i uint64, base int)string
func FormatFloat(f float64, fmt byte, prec, bitSize int)string
```

代码示例
```go
func main() {
	var num1 int = 12
	var num2 float32 = 3.24
	var num3 bool = true
	var str string = "i am sdad"
	str = strconv.FormatInt(int64(num1),10) + strconv.Itoa(num1) +
	strconv.FormatBool(num3) + strconv.FormatFloat(float64(num2),'f',10,32)
	fmt.Printf("%q",str)
}
```

#### b、字符串转基本数据类型
strconv包,若不能转换，<font color='red'>会产生0值</font>  
函数原型
```go
func ParseBool(str string)(value bool,err error)
func ParseFloat(str string, bitSize int)(f float64,err error)
func ParseInt(str string, base, bitSize int)(i int64,err error)
func ParseUint(str string, b, bitSize int)(n uint64,err error)
```

```go
func main() {
	var num1 string = "12"
	var num2 string = "3.24"
	var num3 string = "true"
	a1,_ := strconv.ParseInt(num1,10,32)
	a2,_ := strconv.ParseFloat(num2,32)
	a3,_ := strconv.ParseBool(num3)
	fmt.Println(a1,a2,a3)
}
```

### 指针

 ```go
 func main() {
	var a int = 5
	var p *int = &a
	*p = 27
	fmt.Printf("%v %v",a,&p)
}
 ```

### 值类型和引用类型
a. 值类型<font color='red'>通常</font>在栈区  
特点：变量存储值，内存在栈中分配
```
基本数据类型系列：
int系列，float系列，bool系列，string，数组和结构体struct
```

b. 引用类型<font color='red'>通常</font>在堆区
特点：变量存储的是一个地址，这个地址对应的空间才真正存储数据，内存通常在堆中分配，当没有任何变量引用这个地址时，该地址对应的数据空间就成为一个垃圾，由GC来回收
```
引用数据类型：
指针，切片slice，map，管道chan，interface等都是引用类型
```
<font color='red'>“通常”</font>的出现是由于Golang的逃逸分析机制

### 标识符
#### 命名规则
```
1. 由26个字母，数字，_组成
2. 不能以数字开头
3. 严格区分大小写
4. 不包含空格
5. 不能使用系统保留的关键字
```
系统保留关键字  
||||||  
|:-:|:-:|:-:|:-:|:-:|  
|break|default|fnuc|interface|select|  
|case|defer|go|map|struct|  
|chan|else|goto|package|switch|  
|const|fallthrough|if|range|type|  
|continue|for|inport|return|var|  

#### 命名规范
```
1. package的名字和目录保持一致
2. 变量，函数名，常量采用驼峰命名法
3. 变量，函数名，常量首字母大写表示是共有的，首字母小写表示是私有的
```

