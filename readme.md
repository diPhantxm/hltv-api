# Getting started
This API provides several implementations for parsing data from [HLTV](https://hltv.org). To choose one replace controllers in strategy when creating instance of API in [main.go](cmd/app/main.go) or implement your own. <br/>
*Note: Using ActualImplementation may result in bottleneck because of delay between requests to website or you will get temporarily banned. I recommend using StoredImplementation.* <br/>

1. Define address and port for API server
2. Choose implementation of parsing data
3. If you selected StoredImplementation, add database configuration (Connection string and driver name)

# Stored Implementation
All controllers start goroutine to periodically (every 5 minutes) poll pages on HLTV and parse data (even if it didn't change). <br/>
It uses store (currently database, but you can implement your own class for repositories. For example: storing data in files or another database). <br/>
When new request come, controller form response from the data in store. So the time between actual and stored info is around 5 minutes. <br/>
To collect data from finished events or matches, run [this](cmd/parse/parseFinished.go).

# Actual Implementation
When new request come, it forms the url and parse the data from the page. Then response is returned. <br/>

# Benchmark
|function|time|memory|allocs|
|:--------:|:----:|:------:|:------:|
|ParseEvent|1,829,307,000 ns/op|4,129,744 B/op|18,994 allocs/op|
|ParseMatch|40,668,946 ns/op|5,280,532 B/op|76,710 allocs/op|
|ParseTeam|55,601,717 ns/op|11,166,684 B/op|71,360 allocs/op|
|ParsePlayer|19,932,157 ns/op|4,352,330 B/op|23,424 allocs/op|

# Test Coverage
|Parser|Rate|
|:--------:|:----:|
|Event|83.3%|
|Match|83.5%|
|Team|86.7%|
|Player|82.0%|

Other ~17% of coverage rate is checking error for nil:
```
if err != nil {
    return nil, err
}
```