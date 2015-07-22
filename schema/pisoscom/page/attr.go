package page
import "github.com/seccijr/quintocrawl/model"

var FuncMapper = map[string]func (model.Flat) (model.Flat, error) {
	"SUPERFICIE:.*(\\d+) M² CONSTRUIDOS" =
}

const (
	SURFACE_BUILT_ATTR = iota
	SURFACE_BUILT_REGEXP =
	SURFACE_USE_ATTR = iota
	SURFACE_USE_REGEXP = "SUPERFICIE:.*(\\d+) M² ÚTILES"
	ROOMS_ATTR = iota
	ROOMS_REGEXP = "HABITACIONES: (\\d+) HABITACIONES"
	BATHROOMS_ATTR = iota
	BATHROOMS_REGEXP = "BAÑOS: (\\d+) BAÑOS"
	FLOOR_ATTR = iota
	FLOOR_REGEXP = "PLANTA: (\\d+) BAÑOS"
	AGE_ATTR = iota
	AGE_REGEXP = "ANTIGÜEDAD: (.*)"
	MAINTAIN_ATTR = iota
	MAINTAIN_REGEXP = "CONSERVACIÓN: (.*)"
	CERT_STATE_ATTR = iota
	CERT_STATE_REGEXP = "CLASIFICACIÓN: (.*)"
	CERT_LVL_ATTR = iota
	CERT_LVL_REGEXP = ""
	CERT_INT_ATTR = iota
	CERT_INT_REGEXP = "(\\d+) Kwh/m2año"
	CERT_BROAD_ATTR = iota
	CERT_BROAD_REGEXP = "(\\d+) Kg CO2/m2 año"
	GARDEN_ATTR = iota
	GARAGE_REGEXP = "JARDÍN: (.*)"
	CUPBOARD_ATTR = iota
	HEATING_ATTR = iota
	GROUND_ATTR = iota
	CARPENTRY_ATTR = iota
	AIR_COND_ATTR = iota
	ORIENTATION_ATTR = iota
	POOL_ATTR = iota
	COM_FEES_ATTR = iota
	LIGHT_ATTR = iota
	ELEVATOR_ATTR = iota
	BOXROOM_ATTR = iota
	KIT_EQUIP_ATTR = iota
	LAUNDRY_ATTR = iota
	DINING_ATTR = iota
	TERRACE_ATTR = iota
	FURNITURE_ATTR = iota
	REINF_DOOR_ATTR = iota
	INTERCOM_ATTR = iota
	DOUBLE_GLASS_ATTR = iota
	WATER_ATTR = iota
	CHIMNEY_ATTR = iota
	SUNNY_ATTR = iota
	SEC_SYST_ATTR = iota
	TLF_ATTR = iota
	BALCONY_ATTR = iota
	HAND_FLAT_ATTR = iota
)
