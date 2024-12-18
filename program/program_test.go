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
assert_eq(arr, [1,2,3,4]);
`

var input2 = `
fn construct2DArray(original, m, n) {
	if m*n != original.Len() {
		return [];
	}

	res = [];
	for i = 0; i < original.Len(); i += n {
		res.Push(original[i:i+n]);
	}
	return res;
}

res = construct2DArray([1,2,3,4], 2, 2);
println(res);
assert_eq(res, [[1,2],[3,4]]);
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
assert_eq(ccc, 2);
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
assert_eq(m, {"a":5,"b":4});
d = m["c"]; // 查询一个不存在的key的value
e = d+1;
println("d:", d);
println("e:", e);
assert_eq(e, 1);

// 判断key是否存在
println("contains: ", m.ContainsKey("e"));
println("contains: ", m.ContainsKey("a"));
assert_eq(m.ContainsKey("e"), false);
assert_eq(m.ContainsKey("a"), true);

// 删除key
m.Remove("b");
println(m);
assert_eq(m, {"a":5});
`

var testSlice = `
// 定义及初始化
a = [0;10]; // 长度为10，且每个元素都是0
a[1] = 1; // update
println(a);
assert_eq(a, [0,1,0,0,0,0,0,0,0,0]);

b = [1,2,3];
println(b);
assert_eq(b, [1,2,3]);
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
assert_eq(a, 4);
`

var testPlusAssign = `
a = 10;
a += 6;
println(a);
assert_eq(a, 16);
`

var testMinusAssign = `
a = 10;
a -= 6;
println(a);
assert_eq(a, 4);
`

var input4 = `
fn aaa(nums) {
	nums.Push(a);
}

fn bbb() {
	nums = [];
	aaa(nums);
	println(nums);
	assert_eq(nums, [2]);
}

a = 2;
bbb();
`

var testPrecedence = `
assert_eq(5 + 3 & 6, 0);
assert_eq(5 & 3 + 6, 1);
assert_eq(10 - 4 | 1, 7);
assert_eq((5 + 3) & 6, 0);
assert_eq(3 * 2 & 7, 6);
assert_eq(12 / 3 | 4, 4);
assert_eq(5 + 3 * 2 & 7 | 1, 3);
assert_eq((5 << 2) + (3 & 1), 21);
assert_eq((20 >> 2) - (8 | 1), -4);
assert_eq(0 & 1 + 2, 0);
assert_eq(-1 | 5 * 2, -1);
assert_eq(-5 ^ 4 - 3, -6);
assert_eq((1 << 31) + 5 & 7, 5);
assert_eq(((((5 + 3) * 2 & 7) | 1) ^ 4) << 2, 20);
assert_eq(2 > 1^3, false);
assert_eq(1<<3&77, 8);
assert_eq(1+2<<5, 96);
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
assert_eq(p1.get_name(), "leo");

p1.set_name("leo2");
println(p1.get_name());
assert_eq(p1.get_name(), "leo2");

p2 = person{name:"george"};
println(p2.get_name());
assert_eq(p2.get_name(), "george");
`

var input5 = `
struct person {
	age,
	name,
}

a = [person{1,"leo"}, person{2,"george"}];
b = a[0].name;
println(b);
assert_eq(b, "leo");
fn aaa() {
	return person{age:100};
}
c = aaa();
println(c.age);
assert_eq(c.age, 100);
`

var testAssertEq = `
assert_eq(1, 1);
assert_eq([1,2,3], [1,2,3]);
`

var testClosure = `
fn outer(y) {
	x = 10;
	println("y: ", y);
	fn inner(z) {
		println("z: ", z);
		return x+y+z;
	}
	return inner;
}

b = outer(5);
c = b(3);
println(c);
assert_eq(c, 18);
`

var testLambda = `
x = lambda a, b, c : a + b + c;
println(x(5, 6, 2));
assert_eq(x(5, 6, 2), 13);

f = lambda: "Hello, world!";
println(f());
assert_eq(f(), "Hello, world!");
`

var testCombinationSum = `
fn backtracking(nums, sum, start, result, target, candidates) {
    if sum == target {
        result.Push(nums.Clone());
        return 1;
    }

    if sum > target {
        return 1;
    }

    for i = start; i < candidates.Len(); i++ {
        nums.Push(candidates[i]);
        sum += candidates[i];
        backtracking(nums, sum, i, result, target, candidates);
        sum -= candidates[i];
        nums.Pop();
    }
}

fn combinationSum(candidates, target) {
    result = [];
    backtracking([], 0, 0, result, target, candidates);
    return result;
}

res = combinationSum([2,3,7], 7);
assert_eq(res, [[2,2,3], [7]]);
`

var testIter = `
a = [1,2,3];
b = a.Iter().Map(lambda x: x*2).Collect();
assert_eq(b, [2,4,6]);
`

var input6 = `
enum Option {
	None
	Some(data),
}

a = Option.Some(10);
println(a);
assert_eq(a, Option.Some(10));

b = Option.None;
println(b);
assert_eq(b, Option.None);
`

// TODO:
/*
1. 去掉panic，改成return error  done
3. 支持打印AST done
4. 去掉所有的panic，优化错误信息
5. 支持闭包 done
6. 支持返回函数 done
7. 支持返回空值：return; done
8. 支持lambda表达式 done
9. 支持enum
10. 支持迭代器
*/

func TestProgram_Run(t *testing.T) {
	p := New()
	p.Run(testClosure)
	//fmt.Println(p.AST(testLambda))
}

func TestProgram_RunAll(t *testing.T) {
	tests := []string{
		bubbleSort, testReverseString, testPlusPlus,
		testFor, testVar, testIfElse, testSwitch,
		testInlineComments, testFib, testMap, testSlice,
		testContinue, testMinusMinus, testPlusAssign,
		testMinusAssign, input4, testPrecedence, testStruct,
		input5, testAssertEq, testLambda, testClosure,
		input2, input3, testBacktracking, testCombinationSum}
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
	p := New()
	fmt.Println(p.AST(input))
}

var testReverseString = `
s = "    hello world    ";
println(s.Reverse());
assert_eq(s.Reverse(), "    dlrow olleh    ");

s = s.Trim(" ");
println(s);
assert_eq(s, "hello world");

arr = s.Split(" ");
println(arr);
assert_eq(arr, ["hello", "world"]);

println(s.Contains("hello"));
assert_eq(s.Contains("hello"), true);
`

var testPlusPlus = `
a = 10;
a++;
println(a);
assert_eq(a, 11);
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
assert_eq(a, 20);
`

var testInlineComments = `
a = []; 
a.Push(1);
println(a);
assert_eq(a, [1]);
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
assert_eq(c, 89);
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
assert_eq(a, [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]);
`
