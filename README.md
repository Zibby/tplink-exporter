# tplink-exporter
[![version](https://images.microbadger.com/badges/version/zibby/tplink-exporter.svg)](https://microbadger.com/images/zibby/tplink-exporter) [![image](https://images.microbadger.com/badges/image/zibby/tplink-exporter.svg)](https://microbadger.com/images/zibby/tplink-exporter) ![issues](https://img.shields.io/github/issues-raw/Zibby/tplink-exporter/master.svg) ![last-commit](https://img.shields.io/github/last-commit/Zibby/tplink-exporter.svg) ![pullcount](https://img.shields.io/docker/pulls/zibby/tplink-exporter.svg) ![Build Status](https://jenkins.zibbytechnology.ddns.net/job/tplink-exporter/job/master/badge/icon) ![goscore](https://goreportcard.com/badge/github.com/Zibby/tplink-exporter)

Prometheus exporter for TPlink smart plugs written in go exposed on port 8089. This release uses the gorilla/mux package to handle routes dynamically, so a json file defining your plugs is no longer required. It is strongly recommended to switch to this version as it is much more resource friendly

~~~bash
docker run -d \
  --name=tplink_exporter \
  -p 8089:8089
  zibby/tplink-exporter
~~~

If you have an older plug, try with the variable `legacy="true"`, there was some changes to the json that the newer firmware sends out that is not backwards compatible.


Example:
browse localhost:8089/plugs/anyname?address=XXX.XXX.XXX.XXX&legacy=true 
or if you have older firmware 
localhost:8089/plugs/othername?address=XXX.XXX.XXX.XXX&legacy=false

## Example docker-compose

~~~docker
  tplink-exporter:
    image: zibby/tplink-exporter:latest
    container_name: tplink-exporter
    restart: always
~~~

## Example Prometheus Config

~~~
scrape_configs:
  - job_name: tplink-server
    params:
      address: ['192.168.0.188']
    static_configs:
            - targets: ['tplink-exporter:8089']

  - job_name: tplink-pc
    params:
      address: ['192.168.0.131']
    static_configs:
            - targets: ['tplink-exporter:8089']
~~~
