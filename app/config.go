package app

type ComponentType int

const (
	LINK ComponentType = 1 + iota
	TRIGGER
	HANDLER
)

var ctStr = [...]string{
	"all",
	"link",
	"trigger",
	"handler",
}

type Component struct {
	Name string
	Type ComponentType
	Ref  string
}

func (m ComponentType) String() string { return ctStr[m] }
