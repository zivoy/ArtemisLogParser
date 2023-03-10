package logparser

import (
	"fmt"
	pb "github.com/zivoy/ArtemisLogParser/protobuf"
	"time"
)

type Game struct {
	Data   *Data
	Events []*Event

	items map[int32]string
}

type Event struct {
	Time  time.Time
	Event *pb.AnalyticsEvent_Event

	game *Game
}

type Data struct {
	*pb.Game
	Time time.Time
}

func newGame() *Game {
	return &Game{
		Events: make([]*Event, 0),
		items:  make(map[int32]string),
	}
}

func newData(game *pb.Game) *Data {
	return &Data{
		Game: game,
		Time: time.UnixMilli(game.GetGameTime()),
	}
}

func (g *Game) appendEvent(events *pb.AnalyticsEvent) {
	eventTime := time.UnixMilli(events.GetEventTine())
	for _, e := range events.GetEvents() {
		g.Events = append(g.Events, &Event{Time: eventTime, Event: e, game: g})

		switch e.GetEvent().(type) {
		case *pb.AnalyticsEvent_Event_Object:
			object := e.GetObject()
			var lastObject *pb.ObjectEvent
			for i := len(g.Events) - 2; i > 0; i-- {
				event := g.Events[i].Event
				switch event.GetEvent().(type) {
				case *pb.AnalyticsEvent_Event_Object:
					obj := event.GetObject()
					if obj.GetId() != object.GetId() {
						continue
					}
					lastObject = obj
					i = -1 // break from loop
				}
			}

			if object.GetPosition() == nil && lastObject != nil {
				object.Position = lastObject.Position
			}
			if object.GetRotation() == nil && lastObject != nil {
				object.Rotation = lastObject.Rotation
			}
			if object.GetScripts() == nil && lastObject != nil && lastObject.GetScripts() != nil {
				object.Scripts = lastObject.Scripts
			}
			// todo updated lookup

		case *pb.AnalyticsEvent_Event_Item:
			itemEvent := e.GetItem()
			g.items[itemEvent.Id] = itemEvent.Name
		default:
		}
	}
}

func (g *Game) LookupID(id int32) string {
	name, _ := g.items[id]
	return name
}

func (g *Game) String() string {
	str := fmt.Sprintf("Game Version:\t\t%s\nAnalytics Version:\t%s\nTime:\t\t\t%s\nMetadata:\t\t\"%s\"\n",
		g.Data.GameVersion, g.Data.AnalyticsVersion, g.Data.Time, g.Data.Metadata)

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
		if len(scripts) > 0 {
			scripts = scripts[:len(scripts)-1]
		}

		position := object.GetPosition()
		pos := fmt.Sprintf("{ X:%f Y:%f Z:%f }", position.X, position.Y, position.Z)
		rotation := object.GetRotation()
		rot := fmt.Sprintf("{ W:%f X:%f Y:%f Z:%f }", rotation.W, rotation.X, rotation.Y, rotation.Z)

		event = fmt.Sprintf("Object \"%s\" pos: %s rot: %s scripts: [%s]",
			e.game.LookupID(object.GetId()), pos, rot, scripts)

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

	case *pb.AnalyticsEvent_Event_Death:
		death := e.Event.GetDeath()
		event = fmt.Sprintf("%s died", e.game.LookupID(death.GetId()))
	case *pb.AnalyticsEvent_Event_Respawn:
		respawn := e.Event.GetRespawn()
		event = fmt.Sprintf("%s respawned", e.game.LookupID(respawn.GetId()))
	case *pb.AnalyticsEvent_Event_Health:
		health := e.Event.GetHealth()
		event = fmt.Sprintf("%s now has %f health", e.game.LookupID(health.GetId()), health.GetAmount())
	case *pb.AnalyticsEvent_Event_Damage:
		damage := e.Event.GetDamage()
		event = fmt.Sprintf("%s was injured by %s", e.game.LookupID(damage.GetId()), e.game.LookupID(damage.GetDamager()))
	case *pb.AnalyticsEvent_Event_Stock:
		stock := e.Event.GetStock()
		event = fmt.Sprintf("%s now has %d stock", e.game.LookupID(stock.GetId()), stock.GetStock())

	default:
		event = e.Event.String()
	}

	return fmt.Sprintf("t+%s - %s", e.Time.Sub(e.game.Data.Time), event)
}
