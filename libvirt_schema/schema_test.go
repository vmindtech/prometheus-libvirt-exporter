package libvirt_schema

import (
	"encoding/xml"
	"testing"
)

func TestInstanceUUIDFromSysinfo(t *testing.T) {
	const domainXML = `
<domain type='kvm'>
  <name>instance-00000d86</name>
  <uuid>afcff067-e71a-40e4-ba0c-b0d74f0cd411</uuid>
  <sysinfo type='smbios'>
    <system>
      <entry name='uuid'>afcff067-e71a-40e4-ba0c-b0d74f0cd411</entry>
    </system>
  </sysinfo>
</domain>`

	var domain Domain
	if err := xml.Unmarshal([]byte(domainXML), &domain); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := domain.InstanceUUID(); got != "afcff067-e71a-40e4-ba0c-b0d74f0cd411" {
		t.Fatalf("InstanceUUID() = %q, want afcff067-e71a-40e4-ba0c-b0d74f0cd411", got)
	}
}

func TestInstanceUUIDFallsBackToDomainUUID(t *testing.T) {
	const domainXML = `
<domain type='kvm'>
  <name>instance-00000d86</name>
  <uuid>afcff067-e71a-40e4-ba0c-b0d74f0cd411</uuid>
</domain>`

	var domain Domain
	if err := xml.Unmarshal([]byte(domainXML), &domain); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := domain.InstanceUUID(); got != "afcff067-e71a-40e4-ba0c-b0d74f0cd411" {
		t.Fatalf("InstanceUUID() = %q, want domain uuid fallback", got)
	}
}
