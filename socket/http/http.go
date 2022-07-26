package main

/*
	使用 net/http 实现安全通信

	HTTPS协议是用来解决传统HTTP协议明文传输不安全的问题的, 与普通HTTP协议不同,
	HTTPS协议在传输层(TCP协议)和应用层(HTTP协议)之间增加了一个安全传输层(SSL);

	采用HTTPS协议后, 新网络协议栈在应用层和传输层之间新增了一个安全传输层;
	安全传输层通常采用SSL(Secure Socket Layer)或TLS(Transport Layer Security)
	协议实现(Go标准库支持TLS 1.3版本协议); 这一层负责HTTP协议传输的内容加密、
	通信双方身份验证等;

*/
