package logparser

import (
	pb "artemisLogParser/protobuf"
	"fmt"
)

type Game struct {
	Data   *pb.Game
	Events []*Event

	items map[int32]string
}

type Event struct {
	Time  int64
	Event *pb.AnalyticsEvent_Event

	game *Game
}

func newGame() *Game {
	return &Game{
		Events: make([]*Event, 0),
		items:  make(map[int32]string),
	}
}

func (g *Game) appendEvent(events *pb.AnalyticsEvent) {
	eventList := make([]*Event, 0)
	for _, e := range events.GetEvents() {
		eventList = append(eventList, &Event{Time: events.GetEventTine(), Event: e, game: g})

		switch e.GetEvent().(type) {
		case *pb.AnalyticsEvent_Event_Object:
			// todo updated lookup

		case *pb.AnalyticsEvent_Event_Item:
			itemEvent := e.GetItem()
			g.items[itemEvent.Id] = itemEvent.Name

		case *pb.AnalyticsEvent_Event_Custom:
		case *pb.AnalyticsEvent_Event_Device:
		case *pb.AnalyticsEvent_Event_Map:
		}
	}

	g.Events = append(g.Events, eventList...)
}

func (g *Game) LookupID(id int32) string {
	name, _ := g.items[id]
	return name
}

func (g *Game) String() string {
	str := fmt.Sprintf("Game Version:\t\t%s\nAnalytics Version:\t%s\nTime:\t\t\t%d\nMetadata:\t\t\"%s\"\n",
		g.Data.GameVersion, g.Data.AnalyticsVersion, g.Data.GameTime, g.Data.Metadata)

	str += "Events:\n"
	for _, event := range g.Events {
		str = fmt.Sprintf("%s\t%s\n", str, event)
	}

	return str
}

func (e *Event) String() string {
	var event string

	switch e.Event.GetEvent().(type) {
	case *pb.AnalyticsEvent_Event_Object:
		object := e.Event.GetObject()

		var scripts string
		for _, script := range object.GetScripts() {
			scripts = fmt.Sprintf("Script \"%s\" Data: \"%s\",", e.game.LookupID(script.GetId()), script.GetData())
		}
		scripts = scripts[:len(scripts)-1]

		event = fmt.Sprintf("Object \"%s\" pos: { %s } rot: { %s } scripts: [%s]",
			e.game.LookupID(object.GetId()), object.GetPosition(), object.GetRotation(), scripts)

	case *pb.AnalyticsEvent_Event_Item:
		item := e.Event.GetItem()
		event = fmt.Sprintf("New Item \"%s\" ID: %d ", item.GetName(), item.GetId())

	case *pb.AnalyticsEvent_Event_Custom:
		custom := e.Event.GetCustom()
		event = fmt.Sprintf("Custom Event - %s - \"%X\"", custom.GetType(), custom.GetValue())

	case *pb.AnalyticsEvent_Event_Device:
		device := e.Event.GetDevice()
		var deviceStatus string
		switch device.GetEvent() {
		case pb.DeviceEvent_device_connected:
			deviceStatus = "connected"
		case pb.DeviceEvent_device_disconnected:
			deviceStatus = "disconnected"
		}

		event = fmt.Sprintf("Device %s - %s", deviceStatus, device.GetName())

	case *pb.AnalyticsEvent_Event_Map:
		mapE := e.Event.GetMap()
		event = fmt.Sprintf("Map set to %s", mapE.GetMapName())
	}

	return fmt.Sprintf("%d - %s", e.Time, event)
}
