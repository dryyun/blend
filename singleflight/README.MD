# dataflight 

用于解决缓存失效时，同时多个请求引起的缓存雪崩问题  

完全复制于 [singleflight](https://github.com/golang/groupcache/tree/master/singleflight)     

groupcache 是 `Apache License 2.0` 协议的，所以我只要指明来源就行了把？？不是很了解。。  


## 使用方法  

### 读取同一个 key  
```go

var g Group 

// 判断 key 缓存是否存在，不存在就 g.Do

g.Do(key, func() (i interface{}, e error) {
	
	// 判断 key 缓存时候存在，存在直接返回 
	// 这步基于考虑，同时进入 Do 的几个请求，第一个已经得出结果，第二个如果不加这步，会在执行 fn 的
	
	// fn 主体 ，获取 key 对应的值
	
	// 设置 key 的缓存
	
})


```

### 同时读取不同的 key  

```go
// 一定要使用 goroutine ，否则会被上个key 的读取卡主

go func {
	g.Do(key,fn)
}()

go func {
	g.Do(key,fn)
}()

// 传递值可以通过 channel 进行
```




