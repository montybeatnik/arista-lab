package arista

type PayloadData struct {
	Method  string
	Version int
	Format  string
	Cmds    []string
	ID      int
}
