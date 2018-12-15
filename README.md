# tplink-exporter
[![](https://images.microbadger.com/badges/version/zibby/tplink-exporter.svg)](https://microbadger.com/images/zibby/tplink-exporter "Get your own version badge on microbadger.com")

[![](https://images.microbadger.com/badges/image/zibby/tplink-exporter.svg)](https://microbadger.com/images/zibby/tplink-exporter "Get your own image badge on microbadger.com")

Prometheus endpoint written in go exposed on port 8089.

~~~bash
docker run -d \
  --name=tplink_exporter \
  -p 8089:8089
  -e TPLINK_ADDR="${IP_ADDRESS_OF_PLUG}" \
  zibby/tplink-exporter
~~~

browse localhost:8089/metrics

