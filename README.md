# tplink-exporter
[![version](https://images.microbadger.com/badges/version/zibby/tplink-exporter.svg)](https://microbadger.com/images/zibby/tplink-exporter) [![image](https://images.microbadger.com/badges/image/zibby/tplink-exporter.svg)](https://microbadger.com/images/zibby/tplink-exporter) ![issues](https://img.shields.io/github/issues-raw/Zibby/tplink-exporter/master.svg) ![last-commit](https://img.shields.io/github/last-commit/Zibby/tplink-exporter.svg) ![pullcount](https://img.shields.io/docker/pulls/zibby/tplink-exporter.svg) ![Build Status](https://jenkins.zibbytechnology.ddns.net/job/tplink-exporter/job/master/badge/icon) ![goscore](https://goreportcard.com/badge/github.com/Zibby/tplink-exporter)

Prometheus endpoint written in go exposed on port 8089.

~~~bash
docker run -d \
  --name=tplink_exporter \
  -p 8089:8089
  -e TPLINK_ADDR="[{"name":"first_plug","address":"xxx.xxx.xxx.xxx","legacy": true},{"name":"second_plug", "address":"xxx.xxx.xxx.xxx", "legacy": false}]" \
  zibby/tplink-exporter
~~~

If you have an older plug, try with the variable `legacy="true"`, there was some changes to the json that the newer firmware sends out that is not backwards compatible.

browse localhost:8089/first_plug or localhost:8089/second_plug

## Example docker-compose
~~~docker
  tplink-exporter:
    image: zibby/tplink-exporter:latest
    container_name: tplink-exporter
    restart: always
    environment:
      PLUGS: '[
        {"name":"server","address":"xxx.xxx.xxx.xxx","legacy":true},
        {"name":"pc","address":"xxx.xxx.xxx.xxx","legacy":false},
        {"name":"tv","address":"xxx.xxx.xxx.xxx","legacy":false},
        {"name":"lights","address":"xxx.xxx.xxx.xxx","legacy":false},
        {"name":"christmas-tree","address":"xxx.xxx.xxx.xxx","legacy":false}
      ]'
~~~