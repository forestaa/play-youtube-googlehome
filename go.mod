module github.com/forestaa/play-youtube-googlehome

require (
	github.com/barnybug/go-cast v0.0.0-20181026143840-5be61902f818
	github.com/gogo/protobuf v1.2.0 // indirect
	github.com/hashicorp/mdns v1.0.0
	github.com/ikasamah/homecast v0.0.0-20181120095505-ae646e87e54e
	github.com/micro/mdns v0.0.0-20181201230301-9c3770d4057a
	github.com/miekg/dns v1.1.2 // indirect
	golang.org/x/crypto v0.0.0-20190103213133-ff983b9c42bc // indirect
	golang.org/x/net v0.0.0-20190107210223-45ffb0cd1ba0
	golang.org/x/sys v0.0.0-20190107173414-20be8e55dc7b // indirect
)

replace github.com/barnybug/go-cast v0.0.0-20181026143840-5be61902f818 => github.com/forestaa/go-cast v0.0.0-20190130234615-c427c194fc28

replace github.com/ikasamah/homecast v0.0.0-20181120095505-ae646e87e54e => github.com/forestaa/homecast v0.0.0-20190131005902-c3495bc84acc
