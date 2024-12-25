package warsgame

const (
	FieldWidth  = 2000
	FieldHeight = 1500

	Radius       = 30
	PortalRadius = 75

	Top    = 0 + Radius
	Bottom = FieldHeight - Radius
	Left   = 0 + Radius
	Right  = FieldWidth - Radius

	Acceleration       = 0.9
	Braking            = 0.85
	Friction           = 0.2
	WallElasticity     = 1.2
	PlayerElasticity   = 1.1
	BrickElasticity    = 0.8
	MaxVelocity        = 10
	MaxCollideVelocity = 12.5

	TPS = 60

	LineSpacing     = 1.1
	TextFieldHeight = 30
	TextFieldWidth  = 282
	MaxTextLength   = 20
)
