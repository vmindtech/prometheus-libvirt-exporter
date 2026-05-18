# Prometheus Libvirt Exporter for OpenStack (PortvMind)

Prometheus exporter for libvirt host and VM metrics, maintained by PortvMind. Written in Go; listens on TCP port **9177** at `/metrics` by default.

Forked from [inovex/prometheus-libvirt-exporter](https://github.com/inovex/prometheus-libvirt-exporter) and [zhangjianweibj/prometheus-libvirt-exporter](https://github.com/zhangjianweibj/prometheus-libvirt-exporter).

## OpenStack labels

### `instance_uuid`

Domain metrics use the OpenStack instance UUID as `instance_uuid` instead of the libvirt domain name (e.g. `instance-00000001`). The value is read from SMBIOS sysinfo first, then falls back to the domain root `<uuid>`:

```xml
<sysinfo type='smbios'>
  <system>
    <entry name='uuid'>00000000-0000-4000-8000-000000000001</entry>
  </system>
</sysinfo>
```

### `project_id`

The Nova project UUID is exposed as `project_id`:

```xml
<nova:project uuid="00000000-0000-4000-8000-000000000002">example-project</nova:project>
```

Domains without Nova metadata export an empty `project_id` (`""`).

## Usage

### Running the Exporter

Run as a standalone binary or as a systemd service.

#### Standalone

Build the exporter ([Building and running](#building-and-running)), then:

```sh
./prometheus-libvirt-exporter
```

By default the exporter listens on port **9177** and serves metrics at `/metrics`.

#### With systemd

Sample unit file: [contrib/prometheus-libvirt-exporter.service](contrib/prometheus-libvirt-exporter.service)

```sh
sudo cp contrib/prometheus-libvirt-exporter.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now prometheus-libvirt-exporter
```

### Configuration

| Flag | Description | Default |
|------|-------------|---------|
| `--libvirt.uri` | Libvirt socket URI | `/var/run/libvirt/libvirt-sock-ro` |
| `--libvirt.driver` | Libvirt driver | `qemu:///system` |
| `--exporter.timeout` | Max libvirt API call duration (e.g. `3s`) | `3s` |
| `--exporter.max-concurrent-collects` | Concurrent collects (min: 1). Tune `libvirtd` `max_client_requests` accordingly | `4` |
| `--web.telemetry-path` | HTTP path for metrics | `/metrics` |

Full option list:

```sh
./prometheus-libvirt-exporter --help
```

### Prometheus Configuration

```yaml
scrape_configs:
  - job_name: 'libvirt'
    static_configs:
      - targets: ['localhost:9177']
```

Example query by project:

```promql
libvirt_domain_info_memory_usage_bytes{project_id="00000000-0000-4000-8000-000000000002"}
```

## Building and running

### Requirements

1. Goreleaser: `go install github.com/goreleaser/goreleaser@latest`
2. Taskfile: `go install github.com/go-task/task/v3/cmd/task@latest`

### Local Building

```sh
task build
```

Artifacts (binaries, packages, archives) are written to `dist/`.

## Metrics

Shared labels on domain metrics: `instance_uuid`, `project_id`. Additional labels depend on the metric type.

| Name | Labels | Description |
|------|--------|-------------|
| `libvirt_up` | — | Libvirt scrape status |
| `libvirt_domains` | — | Number of domains |
| `libvirt_domain_timed_out` | `instance_uuid`, `project_id` | Domain metric scrape timeout |
| `libvirt_domain_openstack_info` | `instance_uuid`, `instance_name`, `instance_id`, `flavor_name`, `user_name`, `user_id`, `project_name`, `project_id` | OpenStack metadata |
| `libvirt_domain_info` | `instance_uuid`, `project_id`, `os_type`, `os_type_arch`, `os_type_machine` | Domain OS metadata |
| `libvirt_domain_info_state` | `instance_uuid`, `project_id`, `state_desc` | Domain state code and description |
| `libvirt_domain_info_maximum_memory_bytes` | `instance_uuid`, `project_id` | Maximum allowed memory (bytes) |
| `libvirt_domain_info_memory_usage_bytes` | `instance_uuid`, `project_id` | Memory usage (bytes) |
| `libvirt_domain_info_virtual_cpus` | `instance_uuid`, `project_id` | vCPU count |
| `libvirt_domain_info_cpu_time_seconds_total` | `instance_uuid`, `project_id` | Total CPU time (seconds) |
| `libvirt_domain_memory_stats_*` | `instance_uuid`, `project_id` | Memory stats (swap, balloon, RSS, faults, etc.) |
| `libvirt_domain_block_stats_info` | `instance_uuid`, `project_id`, `disk_type`, `target_bus`, `driver_*`, `source_*`, `target_device`, `serial` | Block device metadata |
| `libvirt_domain_block_stats_*` | `instance_uuid`, `project_id`, `target_device` | Block I/O counters and timings |
| `libvirt_domain_interface_stats_info` | `instance_uuid`, `project_id`, `interface_type`, `source_bridge`, `target_device`, `mac_address`, `model_type`, `mtu_size` | Network interface metadata |
| `libvirt_domain_interface_stats_*` | `instance_uuid`, `project_id`, `target_device` | Network I/O counters |
| `libvirt_domain_job_*` | `instance_uuid`, `project_id` | Domain job metrics (migration, backup, etc.) |
| `libvirt_domain_vcpu_*` | `instance_uuid`, `project_id` [, `vcpu`] | vCPU metrics |
| `libvirt_storage_pool_*` | `storage_pool` | Storage pool metrics (no `instance_uuid` / `project_id`) |

## Example

```text
libvirt_domain_info_memory_usage_bytes{instance_uuid="00000000-0000-4000-8000-000000000001",project_id="00000000-0000-4000-8000-000000000002"} 8.388608e+09
libvirt_domain_block_stats_read_bytes_total{instance_uuid="00000000-0000-4000-8000-000000000001",project_id="00000000-0000-4000-8000-000000000002",target_device="vda"} 1.497283072e+09
libvirt_domain_openstack_info{instance_uuid="00000000-0000-4000-8000-000000000001",flavor_name="m1.small",instance_id="00000000-0000-4000-8000-000000000001",instance_name="example-vm",project_id="00000000-0000-4000-8000-000000000002",project_name="example-project",user_id="00000000-0000-4000-8000-000000000003",user_name="example-user"} 1
libvirt_domain_timed_out{instance_uuid="00000000-0000-4000-8000-000000000001",project_id="00000000-0000-4000-8000-000000000002"} 0
libvirt_storage_pool_capacity_bytes{storage_pool="example-pool"} 12573614080
```