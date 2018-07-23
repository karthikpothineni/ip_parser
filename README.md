# ip_parser
This Repo will be used for extracting location specific information from IP Address

Note:
Update "LOCAL_CITY_FILE_PATH" in "ip_parser.go"."LOCAL_CITY_FILE_PATH" refers to the path where mmdb file is downloaded.

How It Works:
1) Run ip_parser.go
2)Pass "X-Forwarded-For":"IPAddress" in Headers
3)Hit "http://IPAddress:31001/test" API With GET Call

Sample Output:
Country: Spain
Latency for above request: 127.667µs
Country: Spain
Latency for above request: 57.402µs
Country: Spain
Latency for above request: 59.724µs
Country: Spain
Latency for above request: 57.298µs
