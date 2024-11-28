package program

import (
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
arr = [1,2];
if arr[1] > arr[0] {
	println("ok");
}
`

func TestProgram_Run(t *testing.T) {
	p := New()
	p.Run(testSwitch)
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
	a = 9;
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
