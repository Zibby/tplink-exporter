# tplink-exporter
[![](https://images.microbadger.com/badges/version/zibby/tplink-exporter.svg)](https://microbadger.com/images/zibby/tplink-exporter "Get your own version badge on microbadger.com") [![](https://images.microbadger.com/badges/image/zibby/tplink-exporter.svg)](https://microbadger.com/images/zibby/tplink-exporter "Get your own image badge on microbadger.com") 
![](https://img.shields.io/github/issues-raw/Zibby/tplink-exporter/master.svg) ![](https://img.shields.io/github/last-commit/Zibby/tplink-exporter.svg) ![](https://img.shields.io/docker/automated/zibby/tplink-exporter.svg) ![](https://img.shields.io/docker/pulls/zibby/tplink-exporter.svg)
![Build Status](https://jenkins.zibbytechnology.ddns.net/job/tplink-exporter/job/master/badge/icon)

Prometheus endpoint written in go exposed on port 8089.

~~~bash
docker run -d \
  --name=tplink_exporter \
  -p 8089:8089
  -e TPLINK_ADDR="${IP_ADDRESS_OF_PLUG}" \
  -e LATER_FW="true"
  zibby/tplink-exporter
~~~

If you have an older plug, try with the env variable LASTER_FW="false", there was some changes to the json that the newer firmware sends out that is not backwards compatible. 

browse localhost:8089/metrics

