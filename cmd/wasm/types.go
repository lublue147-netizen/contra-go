package main

type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StateStageClear
	StateGameOver
	StateVictory
)

type Direction int

const (
	DirRight Direction = iota
	DirLeft
	DirUp
	DirDown
	DirUpRight
	DirUpLeft
	DirDownRight
	DirDownLeft
)

type WeaponType int

const (
	WeaponNormal WeaponType = iota // Standard rifle
	WeaponSpread                   // Spread Gun (S) - 5-way spread
	WeaponLaser                    // Laser Gun (L) - Piercing beam
	WeaponMachine                  // Machine Gun (M) - Rapid fire
	WeaponBarrier                  // Barrier (B) - Invincibility
)

type PlayerState int

const (
	PlayerIdle PlayerState = iota
	PlayerRunning
	PlayerJumping
	PlayerCrouching
	PlayerSwimming
	PlayerDying
)

type Bullet struct {
	X, Y        float64
	VX, VY      float64
	Damage      int
	Type        WeaponType
	IsEnemy     bool
	Life        int
	Radius      float64
	Color       string
	PierceCount int
}

type Particle struct {
	X, Y       float64
	VX, VY     float64
	Life       int
	MaxLife    int
	Color      string
	Size       float64
	IsSpark    bool
}

type FloatingText struct {
	Text  string
	X, Y  float64
	Life  int
	Color string
}

type DropItem struct {
	X, Y       float64
	VY         float64
	Type       WeaponType
	Letter     string
	Color      string
	Active     bool
	Life       int
	OnGround   bool
}

type EnemyType int

const (
	EnemyTypeSoldier EnemyType = iota
	EnemyTypeSniper
	EnemyTypeTurret
	EnemyTypeFalconCapsule
	EnemyTypeSensor
)

type Enemy struct {
	ID         int
	Type       EnemyType
	X, Y       float64
	VX, VY     float64
	Width      float64
	Height     float64
	HP         int
	MaxHP      int
	ScoreValue int
	FacingLeft bool
	AimAngle   float64
	ShootTimer int
	StateTimer int
	OnGround   bool
	Active     bool
	DropsItem  WeaponType
	HasDrop    bool
	AnimFrame  int
	AnimTimer  int
}

type BossPart struct {
	Name       string
	OffsetX    float64
	OffsetY    float64
	Width      float64
	Height     float64
	HP         int
	MaxHP      int
	Destroyed  bool
	ShootTimer int
	Angle      float64
	Flashing   int
}

type Boss struct {
	X, Y         float64
	Width        float64
	Height       float64
	Active       bool
	Defeated     bool
	CoreHP       int
	MaxCoreHP    int
	TopTurret    BossPart
	LeftTurret   BossPart
	RightTurret  BossPart
	SoldierTimer int
	Flashing     int
	DeathTimer   int
	DefeatAnim   int
}

type Platform struct {
	X, Y, W, H float64
	IsWater    bool
	IsDropThru bool // Can drop down with Down+Jump
}
