package fakes

type FakeDisk struct {
	cid string

	NeedsMigrationInputs []NeedsMigrationInput
	needsMigrationOutput needsMigrationOutput
}

type NeedsMigrationInput struct {
	Size            int
	CloudProperties map[string]interface{}
}

type needsMigrationOutput struct {
	needsMigration bool
}

func NewFakeDisk(cid string) *FakeDisk {
	return &FakeDisk{
		cid:                  cid,
		NeedsMigrationInputs: []NeedsMigrationInput{},
	}
}

func (d *FakeDisk) CID() string {
	return d.cid
}

func (d *FakeDisk) NeedsMigration(size int, cloudProperties map[string]interface{}) bool {
	d.NeedsMigrationInputs = append(d.NeedsMigrationInputs, NeedsMigrationInput{
		Size:            size,
		CloudProperties: cloudProperties,
	})

	return d.needsMigrationOutput.needsMigration
}

func (d *FakeDisk) SetNeedsMigrationBehavior(needsMigration bool) {
	d.needsMigrationOutput = needsMigrationOutput{
		needsMigration: needsMigration,
	}
}
