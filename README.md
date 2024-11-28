# leveling
server
```shell
go run cmd/server/main.go
```
client
```shell
go run cmd/client/main.go
```

## design pattern
### Service Locator
see [ref](https://gameprogrammingpatterns.com/service-locator.html)

> service 類別定義了抽象的操作集合，由 service provider 具體實作。 <br>
> locator 透過尋找合適的 provider 存取特定的服務，同時隱藏具體的實作類型以及尋找他的過程

有些物件或系統，在各個地方遊走，幾乎遍佈在整個專案的角落，例如記憶體管理、日誌、亂數等，整個專案都需要他們。

選擇用靜態類或單例模式？

想像一下公共電話簿，要找到某人電話，只需要提供他的姓名。
如果這個人換了電話，只需要更新公共電話簿，而不需要告知每個人。


