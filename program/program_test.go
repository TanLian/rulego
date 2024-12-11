package program

import (
	"fmt"
	"testing"
)

// bubbleSort 冒泡排序
var bubbleSort = `
fn BubbleSort(arr) {
	n = arr.Len();
	for i = 0; i < n-1; i++ {
		// 标志位，用于优化，记录这一轮是否有交换
		swapped = false;
		for j = 0; j < n-i-1; j += 1 {
			if arr[j] > arr[j+1] {
				tmp = arr[j];
				arr[j] = arr[j+1];
				arr[j+1] = tmp;
				swapped = true;
			}
		}

		// 如果没有交换，说明数组已经有序，可以提前结束
		if !swapped {
			break;
		}
	}
}

arr = [3,2,1,4];
BubbleSort(arr);
println(arr);
`

var input2 = `
fn construct2DArray(original, m, n) {
	if m*n != len(original) {
		return [];
	}

	res = [];
	for i = 0; i < len(original); i += n {
		res.Push(original[i:i+n]);
	}
	return res;
}

println(construct2DArray([1,2,3,4], 2, 2));
`

var input3 = `
fn subarraySum(nums, k) {
    sum = [];
    for i = 0; i < nums.Len(); i++ {
        tmp = 0;
        if i > 0 {
            tmp = nums[i] + sum[i-1];
        } else {
            tmp = nums[i];
        }
        sum.Push(tmp);
    }

    m = {};
    res = 0;
    for i = 0; i < nums.Len(); i += 1 {
        if sum[i] == k {
            res += 1;
        }
        res += m[sum[i]-k];
		m[sum[i]] += 1;
    }
    return res;
}

ccc = subarraySum([1,2,3],3);
println("res: ", ccc);
`

var testMap = `
// 定义并初始化
m = {"a":3};

// insert
m["b"] = 2;

// update
m["a"] += 2;
m["b"] += 2;
println(m);

d = m["c"]; // 查询一个不存在的key的value
e = d+1;
println("d:", d);
println("e:", e);

// 判断key是否存在
println("contains: ", m.ContainsKey("e"));
println("contains: ", m.ContainsKey("a"));

// 删除key
m.Remove("b");
println(m);
`

var testSlice = `
// 定义及初始化
a = [0;10]; // 长度为10，且每个元素都是0
a[1] = 1; // update
println(a);

b = [1,2,3];
println(b);
`

var testContinue = `
for i = 0; i < 5; i++ {
	if i == 3 {
		continue;
	}
	println(i);
}
`

var testMinusMinus = `
a = 5;
a--;
println(a);
`

var testPlusAssign = `
a = 10;
a += 6;
println(a);
`

var testMinusAssign = `
a = 10;
a -= 6;
println(a);
`

var input4 = `
fn aaa(nums) {
	a = 1;
	nums.Push(1);
	//nums.Push(b);
}

fn bbb() {
	nums = [];
	b = a;
	aaa(nums);
	println(nums);
}

a = 2;
c = 3;
bbb();
`

var testPrecedence = `
println(5 + 3 & 6); // 0
println(5 & 3 + 6); // 1
println(10 - 4 | 1); // 7
println((5 + 3) & 6); // 0
println(3 * 2 & 7); // 6
println(12 / 3 | 4); // 4
println(5 + 3 * 2 & 7 | 1); // 3
println((5 + 3) * (2 & 7) | 1); // 17
println((5 << 2) + (3 & 1)); // 21
println((20 >> 2) - (8 | 1)); // -4
println(0 & 1 + 2); // 0
println(-1 | 5 * 2); // -1
println(-5 ^ 4 - 3); // -6
println((1 << 31) + 5 & 7); // 5
println(((((5 + 3) * 2 & 7) | 1) ^ 4) << 2); // 20
println(2 > 1^3); // false
println(1<<3&77); // 8
println(1+2<<5); // 96
`

var testStruct = `
struct person {
	age,
	name,
}

impl person {
	fn get_name(self) {
		self.name
	}

	fn set_name(self, name) {
		self.name = name;
	}
}

p1 = person{1,"leo"};
println(p1.get_name());

p1.set_name("leo2");
println(p1.get_name());

p2 = person{name:"george"};
println(p2.get_name());
`

var input5 = `
struct person {
	age,
	name,
}

a = [person{1,"leo"}, person{2,"george"}];
b = a[0].name;
println(b);

fn aaa() {
	return person{age:100};
}
c = aaa();
println(c.age);
`

var input6 = `
fn fib(n) {
	if n < 2 {
		return 1;
	}
	return fib(n-1) + fib(n-2);
}

a = fib(10);
println(a);
`

// TODO:
/*
1. 去掉panic，改成return error  done
2. 报错行号支持 done
3. 支持打印AST done
*/

func TestProgram_Run(t *testing.T) {
	p := New()
	p.Run(input6)
}

func TestProgram_RunAll(t *testing.T) {
	tests := []string{
		bubbleSort, testReverseString, testPlusPlus,
		testFor, testVar, testIfElse, testSwitch,
		testInlineComments, testFib, testMap, testSlice,
		testContinue, testMinusMinus, testPlusAssign,
		testMinusAssign, input4, testPrecedence, testStruct,
		input5, input6}
	for i, v := range tests {
		p := New()
		if err := p.Run(v); err != nil {
			fmt.Println("err: ", err, " i: ", i)
			break
		}
	}
}

func TestProgramAST_Run(t *testing.T) {
	input := `
a = 10;
a += 6;
println(a);
`
	p := New()
	fmt.Println(p.AST(input))
}

var testReverseString = `
s = "    hello world    ";
println(s.Reverse());
s = s.Trim(" ");
println(s);

arr = s.Split(" ");
println(arr);

println(s.Contains("hello"));
`

var testPlusPlus = `
a = 10;
a++;
println(a);
`

var testFor = `
for a = 0; a <= 10; a++ {
	println(a);
}
`

// 变量
// 语法：变量名 = 值;
var testVar = `
	intVal = 1; // 整形变量
	floatVal = 3.14; // 浮点型变量
	boolVal = true; // bool类型
	strVal = "hello world"; // 字符串
	arrVal = ["hello", "world", 100, 101, 102.1]; // slice，各个元素可以是不同类型
	mapVal = {"name":"leo", "age":20}; // map
`

// if else
/*
语法：
if expr {
	state1;
	state2;
} else if expr {
	state1;
	state2;
} else {
	state1;
	state2;
}
*/
var testIfElse = `
	a = 21;
	if a > 20 {
		println("greater than 20");
	} else if a >= 10 {
		println("greater or equal to 10");
	} else {
		println("less than 10");
	}
`

var testSwitch = `
a = 10;
switch a {
case 1:
	println("111");
	a++;
case 10:
	println("hhhh");
	a += 10;
default:
	println("dddd");
	a += 5;
}
println(a);
`

var testInlineComments = `
a = []; 
// a.Push(1);
println(a);
`

var testFib = `
fn fib(n) {
	if n < 2 {
		return 1;
	}
	return fib(n-2)+fib(n-1);
}

c = fib(10);
println(c);
`

// testBacktracking 测试全排列
var testBacktracking = `
fn backtracking(nums, ns, used, result) {
    if ns.Len() == nums.Len() {
        result.Push(ns.Clone());
    }

    for i = 0; i < nums.Len(); i++ {
        if used[i] {
            continue;
        }

        used[i] = true;
        ns.Push(nums[i]);
        backtracking(nums, ns, used, result);
        used[i] = false;
        ns.Pop();
    }
}

fn permute(nums) {
    used = [false; nums.Len()];
    ns = [];
    result = [];
    backtracking(nums, ns, used, result);
    return result;
}

a = permute([1,2,3]);
println(a);
`
