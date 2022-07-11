package main

/*

	包内测试: 将测试代码放在与被测包同名的包中的测试方法
	包外测试: 将测试代码放在名为被测包包名+"_test"的包中的测试方法

	包内测试这种方法本质上是一种白盒测试方法, 由于测试代码与被测包源码在
	同一包名下, 测试代码可以访问该包下的所有符号, 无论是导出符号还是未导
	出符号; 并且由于包的内部实现逻辑对测试代码是透明的, 包内测试可以更为
	直接地构造测试数据和实施测试逻辑, 可以很容易地达到较高的测试覆盖率;
	因此对于追求高测试覆盖率的项目而言, 包内测试是不二之选

	但是包内测试也会有以下问题:
	- 测试代码自身需要经常性的维护
		包内测试的白盒测试本质意味着它是一种面向实现的测试, 测试代码的测试
		数据构造和测试逻辑通常与被测包的特定数据结构设计和函数/方法的具体
		实现逻辑是紧耦合的,	而包的内部实现逻辑又是易变的, 其优化调整是一
		种经常性行为, 这就意味着采用包内测试的测试代码也需要经常性的维护;
	- 硬伤: 包循环引用
		Go标准库对strings包的测试采用包外测试, 因为 testing 包中导入了
		strings 包, 如果在 strings 包内进行测试导入 testing 包时会引起
		循环引用;

	与包内测试本质是面向实现的白盒测试不同, 包外测试的本质是一种面向接口
	的黑盒测试; 接口指被测试包对外导出的API, 这些API是被测包与外部交互的契约;
	契约一旦确定就会长期保持稳定, 无论被测包的内部实现逻辑和数据结构设计
	如何调整与优化, 一般都不会影响这些契约; 这一本质让包外测试代码与被测试
	包充分解耦, 使得针对这些导出API进行测试的包外测试代码表现出十分健壮的
	特性, 即很少随着被测代码内部实现逻辑的调整而进行调整和维护;

	包外测试这种纯黑盒的测试还有一个功能域之外的好处, 那就是可以更加聚焦
	地从用户视角验证被测试包导出API的设计的合理性和易用性;

	同时包外测试存在测试盲区:
		由于测试代码与被测试目标并不在同一包名下, 测试代码仅有权访问被测包
		的导出符号, 并且仅能通过导出API这一有限的"窗口"并结合构造特定数据来
		验证被测包行为; 在这样的约束下, 很容易出现对被测试包的测试覆盖不足
		的情况。
	为了解决盲区问题可以为被测包"安插后门", 即在被测包内定义 export_test.go
	文件, 该文件中的代码位于被测包名下, 但既不会被包含在正式产品代码中, 也不
	包含任何测试代码, 而仅用于将被测包的内部符号在测试阶段暴露给包外测试代码,
	或者定义一些辅助包外测试的代码, 比如扩展被测包的方法集合;


	如何选择:
	包外测试由于将测试代码放入独立的包中, 更适合编写偏向集成测试的用例,
	它可以任意导入外部包, 并测试与外部多个组件的交互;
	包内测试更聚焦于内部逻辑的测试, 通过给函数/方法传入一些特意构造的
	数据的方式来验证内部逻辑的正确性, 比如net/http包的response_test.go;

	当运用包外测试与包内测试共存的方式时, 可考虑让包外测试和包内测试
	聚焦于不同的测试类别;
*/
