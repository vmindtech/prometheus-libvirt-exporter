# Prometheus Libvirt Exporter (PortvMind)

[![Build and Test](https://github.com/vmindtech/prometheus-libvirt-exporter/actions/workflows/build_and_test.yml/badge.svg)](https://github.com/vmindtech/prometheus-libvirt-exporter/actions/workflows/build_and_test.yml)
[![Lint Go Code](https://github.com/vmindtech/prometheus-libvirt-exporter/actions/workflows/lint.yml/badge.svg)](https://github.com/vmindtech/prometheus-libvirt-exporter/actions/workflows/lint.yml)

PortvMind tarafından bakımı yapılan bu exporter, libvirt üzerinden host ve sanal makine (VM) metriklerini Prometheus formatında sunar. Go ile yazılmıştır ve varsayılan olarak TCP portu **9177** üzerinde `/metrics` yolunda dinler.

[go-libvirt](https://github.com/digitalocean/go-libvirt) (DigitalOcean) paketi üzerine kuruludur; libvirt ile saf Go RPC arayüzü kullanır. Libvirt API referansı: [libvirt.org](https://libvirt.org/html/index.html).

Bu proje [inovex/prometheus-libvirt-exporter](https://github.com/inovex/prometheus-libvirt-exporter) ve [zhangjianweibj/prometheus-libvirt-exporter](https://github.com/zhangjianweibj/prometheus-libvirt-exporter) projelerinden türetilmiştir.

## OpenStack `project_id`

OpenStack Nova ortamında çalışan domain metriklerine, libvirt domain XML metadata’sındaki Nova `project` UUID’si `project_id` label’ı olarak eklenir. Böylece proje bazlı filtreleme ve sorgulama Prometheus/Grafana tarafında doğrudan yapılabilir.

`project_id` değeri şu kaynaktan okunur:

```xml
<nova:project uuid="93ac887ce5794c778320a88c3024b1ad">my-project</nova:project>
```

Nova metadata’sı olmayan domain’lerde `project_id` boş string (`""`) olarak yayınlanır.

## Usage

### Running the Exporter

Exporter’ı tek başına binary veya systemd servisi olarak çalıştırabilirsiniz.

#### Standalone

Exporter’ı derleyin ([Building and running](#building-and-running)), ardından:

```sh
./prometheus-libvirt-exporter
```

Varsayılan olarak exporter **9177** portunda dinler ve metrikleri `/metrics` altında sunar.

#### With systemd

Örnek systemd unit dosyası: [contrib/prometheus-libvirt-exporter.service](contrib/prometheus-libvirt-exporter.service)

```sh
sudo cp contrib/prometheus-libvirt-exporter.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now prometheus-libvirt-exporter
```

### Configuration

Komut satırı bayrakları:

| Flag | Açıklama | Varsayılan |
|------|----------|------------|
| `--libvirt.uri` | Libvirt soket URI’si | `/var/run/libvirt/libvirt-sock-ro` |
| `--libvirt.driver` | Libvirt sürücüsü | `qemu:///system` |
| `--exporter.timeout` | Libvirt API çağrı üst sınırı (örn. `3s`) | `3s` |
| `--exporter.max-concurrent-collects` | Eşzamanlı collect sayısı (min: 1). `libvirtd` içindeki `max_client_requests` değerini buna göre ayarlayın | `4` |
| `--web.telemetry-path` | Metriklerin yayınlandığı HTTP yolu | `/metrics` |

Tüm seçenekler için:

```sh
./prometheus-libvirt-exporter --help
```

### Prometheus Configuration

`prometheus.yml` örneği:

```yaml
scrape_configs:
  - job_name: 'libvirt'
    static_configs:
      - targets: ['localhost:9177']
```

Proje bazlı örnek sorgu:

```promql
libvirt_domain_info_memory_usage_bytes{project_id="93ac887ce5794c778320a88c3024b1ad"}
```

## Building and running

### Requirements

1. Goreleaser: `go install github.com/goreleaser/goreleaser@latest`
2. Taskfile: `go install github.com/go-task/task/v3/cmd/task@latest`

### Local Building

```sh
task build
```

Derleme çıktıları `dist/` klasöründe oluşur (binary, paketler, arşivler).

## Metrics

Domain ile ilişkili metriklerde ortak label’lar: `domain`, `project_id`. Ek label’lar metrik türüne göre değişir.

| Name | Labels | Description |
|------|--------|-------------|
| `libvirt_up` | — | Libvirt scrape durumu |
| `libvirt_domains` | — | Domain sayısı |
| `libvirt_domain_timed_out` | `domain`, `project_id` | Domain metrik scrape zaman aşımı |
| `libvirt_domain_openstack_info` | `domain`, `instance_name`, `instance_id`, `flavor_name`, `user_name`, `user_id`, `project_name`, `project_id` | OpenStack metadata (toplu label seti) |
| `libvirt_domain_info` | `domain`, `project_id`, `os_type`, `os_type_arch`, `os_type_machine` | Domain OS metadata |
| `libvirt_domain_info_state` | `domain`, `project_id`, `state_desc` | Domain durum kodu ve açıklaması |
| `libvirt_domain_info_maximum_memory_bytes` | `domain`, `project_id` | İzin verilen maksimum bellek (byte) |
| `libvirt_domain_info_memory_usage_bytes` | `domain`, `project_id` | Bellek kullanımı (byte) |
| `libvirt_domain_info_virtual_cpus` | `domain`, `project_id` | vCPU sayısı |
| `libvirt_domain_info_cpu_time_seconds_total` | `domain`, `project_id` | Toplam CPU süresi (saniye) |
| `libvirt_domain_memory_stats_*` | `domain`, `project_id` | Bellek istatistikleri (swap, balloon, RSS, fault, vb.) |
| `libvirt_domain_block_stats_info` | `domain`, `project_id`, `disk_type`, `target_bus`, `driver_*`, `source_*`, `target_device`, `serial` | Blok cihaz metadata |
| `libvirt_domain_block_stats_*` | `domain`, `project_id`, `target_device` | Blok I/O sayaçları ve süreleri |
| `libvirt_domain_interface_stats_info` | `domain`, `project_id`, `interface_type`, `source_bridge`, `target_device`, `mac_address`, `model_type`, `mtu_size` | Ağ arayüzü metadata |
| `libvirt_domain_interface_stats_*` | `domain`, `project_id`, `target_device` | Ağ I/O sayaçları |
| `libvirt_domain_job_*` | `domain`, `project_id` | Domain job (migration, backup, vb.) metrikleri |
| `libvirt_domain_vcpu_*` | `domain`, `project_id` [, `vcpu`] | vCPU metrikleri |
| `libvirt_storage_pool_*` | `storage_pool` | Storage pool metrikleri (domain/project_id yok) |

## Example

```text
libvirt_domain_info_memory_usage_bytes{domain="instance-0001e06e",project_id="93ac887ce5794c778320a88c3024b1ad"} 1.7179869184e+10
libvirt_domain_block_stats_read_bytes_total{domain="instance-0001e06e",project_id="93ac887ce5794c778320a88c3024b1ad",target_device="vda"} 1.497283072e+09
libvirt_domain_interface_stats_receive_bytes_total{domain="instance-0001e06e",project_id="93ac887ce5794c778320a88c3024b1ad",target_device="tapab672ce4-11"} 1.589638794e+09
libvirt_domain_openstack_info{domain="instance-0001e06e",flavor_name="z1.4xlarge",instance_id="a12423b02-4a36-4530-bf25-acb8ba80b1b1",instance_name="openstackInstanceName",project_id="93ac887ce5794c778320a88c3024b1ad",project_name="openstackProjectName",user_id="",user_name="openstackUserName"} 1
libvirt_domain_timed_out{domain="instance-0001e06e",project_id="93ac887ce5794c778320a88c3024b1ad"} 0
libvirt_storage_pool_capacity_bytes{storage_pool="testpool"} 12573614080
```

## License

MIT — bkz. [LICENSE](LICENSE).
