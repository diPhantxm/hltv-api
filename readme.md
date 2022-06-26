# Getting started
This API provides several implementations for parsing data from [HLTV](https://hltv.org). To choose one replace controllers in strategy when creating instance of API in [main.go](cmd/app/main.go) or implement your own.
*Note: Using ActualImplementation may result in bottleneck because of delay between requests to website or you will get temporarily banned. I recommend using StoredImplementation.*

1. Define address and port for API server
2. Choose implementation of parsing data
3. If you selected StoredImplementation, add database configuration (Connection string and driver name)

# Stored Implementation
All controllers start goroutine to periodically (every 5 minutes) poll pages on HLTV and parse data (even if it didn't change).
It uses store (currently database, but you can implement your own class for repositories. For example: storing data in files or another database).
When new request come, controller form response from the data in store. So the time between actual and stored info is around 5 minutes.
To collect data from finished events or matches, run [this](cmd/parse/parseFinished.go).

# Actual Implementation
When new request come, it forms the url and parse the data from the page. Then response is returned.

# Benchmark
|function|time|memory|allocs|
|:--------:|:----:|:------:|:------:|
|ParseEvent|9990612 ns/op|1788898 B/op|15780 allocs/op|
|ParseEvent|37607000 ns/op|5186073 B/op|76247 allocs/op|
|ParseTeam|49537895 ns/op|10724525 B/op|67998 allocs/op|
|ParsePlayer|279530825 ns/op|18018990 B/op|64330 allocs/op|