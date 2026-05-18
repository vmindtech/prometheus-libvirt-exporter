package libvirt_schema

import (
	"encoding/xml"
	"testing"
)

const exampleInstanceUUID = "00000000-0000-4000-8000-000000000001"

func TestInstanceUUIDFromSysinfo(t *testing.T) {
	const domainXML = `
<domain type='kvm'>
  <name>instance-00000001</name>
  <uuid>` + exampleInstanceUUID + `</uuid>
  <sysinfo type='smbios'>
    <system>
      <entry name='uuid'>` + exampleInstanceUUID + `</entry>
    </system>
  </sysinfo>
</domain>`

	var domain Domain
	if err := xml.Unmarshal([]byte(domainXML), &domain); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := domain.InstanceUUID(); got != exampleInstanceUUID {
		t.Fatalf("InstanceUUID() = %q, want %s", got, exampleInstanceUUID)
	}
}

func TestInstanceUUIDFallsBackToDomainUUID(t *testing.T) {
	const domainXML = `
<domain type='kvm'>
  <name>instance-00000001</name>
  <uuid>` + exampleInstanceUUID + `</uuid>
</domain>`

	var domain Domain
	if err := xml.Unmarshal([]byte(domainXML), &domain); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := domain.InstanceUUID(); got != exampleInstanceUUID {
		t.Fatalf("InstanceUUID() = %q, want domain uuid fallback", got)
	}
}
