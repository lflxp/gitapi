package gitapi

type Git struct {
	Url	string
	Version string
}

func (this *Git) init() {
	this.Version = ExecCommand("git version")
}

func (this *Git) status(path string) string {
	println("ok")
	return "ok"
}
