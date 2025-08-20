# TradeMatching

<br />

## 限價單（Limit Order）的核心邏輯
<br />
買入限價單：
<br />
指定最高願意買入的價格，只會成交價格 ≤ 此價格的賣單。
<br />
從當前最低賣價開始成交，直到數量用完、賣單沒了，或遇到價格高於限價的賣單為止。
<br />
賣出限價單：
<br />
指定最低願意賣出的價格，只會成交價格 ≥ 此價格的買單。
<br />
從當前最高買價開始成交，直到數量用完、買單沒了，或遇到價格低於限價的買單為止。
<br />
特點：
<br />
可能立即成交（全額或部分）或完全不成交。
<br />
未成交的部分會掛在訂單簿中，等待將來有符合價格的對手單成交。
<br />
成交速度取決於市場對手單的價格與數量。
<br />
<br />

![image](https://github.com/lzz0826/MatchmakingTransactions/blob/main/imges/002.png)

## 市價單（Market Order）的核心邏輯:
買入市價單：
<br />
不指定價格，以當前賣單的最低價格開始成交，直到數量用完或賣單沒了。
<br />
賣出市價單：
<br />
不指定價格，以當前買單的最高價格開始成交，直到數量用完或買單沒了。
<br />
特點：
<br />
不會掛在訂單簿（因為立刻成交）。
<br />
可能會部分成交（如果流動性不足）。
<br />

<br />
<br />

![image](https://github.com/lzz0826/MatchmakingTransactions/blob/main/imges/001.png)
<br />
<br />
測試用API: test/Trad.postman_collection.json
<br />

### 瀏覽 Swagger UI
http://localhost:8081/swagger/index.html
<br />
初始化 Swagger 文件(有改動需要執行): swag init
<br />
![image](https://github.com/lzz0826/MatchmakingTransactions/blob/main/imges/003.png)


### TODO
訂單.交易記錄服務拆分(使用中間鍵MQ)
<br />
DB事務回滾後 買 賣 單簿(GO內部記憶體MAP會不一致)

