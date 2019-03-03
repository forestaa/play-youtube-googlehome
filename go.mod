module github.com/forestaa/play-youtube-googlehome

require (
	github.com/barnybug/go-cast v0.0.0-20181026143840-5be61902f818
	github.com/gogo/protobuf v1.2.0 // indirect
	github.com/hashicorp/mdns v1.0.0
	github.com/ikasamah/homecast v0.0.0-20181120095505-ae646e87e54e
	github.com/micro/mdns v0.1.0
	golang.org/x/net v0.0.0-20190125091013-d26f9f9a57f3
)

replace github.com/barnybug/go-cast v0.0.0-20181026143840-5be61902f818 => github.com/forestaa/go-cast v0.0.0-20190303135850-a385b201c5f5

replace github.com/ikasamah/homecast v0.0.0-20181120095505-ae646e87e54e => github.com/forestaa/homecast v0.0.0-20190303140315-c7cf39f4cf00
