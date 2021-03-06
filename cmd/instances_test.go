package cmd_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/cmd"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	fakedir "github.com/cloudfoundry/bosh-cli/director/directorfakes"
	fakeui "github.com/cloudfoundry/bosh-cli/ui/fakes"
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
)

var _ = Describe("InstancesCmd", func() {
	var (
		ui         *fakeui.FakeUI
		deployment *fakedir.FakeDeployment
		command    InstancesCmd
	)

	BeforeEach(func() {
		ui = &fakeui.FakeUI{}
		deployment = &fakedir.FakeDeployment{}
		command = NewInstancesCmd(ui, deployment)
	})

	Describe("Run", func() {
		var (
			opts InstancesOpts
		)

		BeforeEach(func() {
			opts = InstancesOpts{}
		})

		act := func() error { return command.Run(opts) }

		Context("when instances are successfully retrieved", func() {
			var (
				infos          []boshdir.VMInfo
				procCPUTotal   float64
				procMemPercent float64
				procMemKB      uint64
				procUptime     uint64
			)

			BeforeEach(func() {
				index1 := 1
				index2 := 2

				procCPUTotal = 50.40
				procMemPercent = 11.10
				procMemKB = 8000
				procUptime = 349350

				infos = []boshdir.VMInfo{
					{
						JobName:      "job-name",
						Index:        &index1,
						ProcessState: "in1-process-state",
						ResourcePool: "in1-rp",

						IPs: []string{"in1-ip1", "in1-ip2"},
						DNS: []string{"in1-dns1", "in1-dns2"},

						State:              "in1-state",
						VMID:               "in1-cid",
						AgentID:            "in1-agent-id",
						ResurrectionPaused: false,
						Ignore:             true,
						DiskIDs:            []string{"diskcid1", "diskcid2"},

						Vitals: boshdir.VMInfoVitals{
							Load: []string{"0.02", "0.06", "0.11"},

							CPU:  boshdir.VMInfoVitalsCPU{Sys: "0.3", User: "1.2", Wait: "2.1"},
							Mem:  boshdir.VMInfoVitalsMemSize{Percent: "20", KB: "2000"},
							Swap: boshdir.VMInfoVitalsMemSize{Percent: "21", KB: "2100"},

							Disk: map[string]boshdir.VMInfoVitalsDiskSize{
								"system":     boshdir.VMInfoVitalsDiskSize{Percent: "35"},
								"ephemeral":  boshdir.VMInfoVitalsDiskSize{Percent: "45"},
								"persistent": boshdir.VMInfoVitalsDiskSize{Percent: "55"},
							},
						},

						Processes: []boshdir.VMInfoProcess{
							{
								Name:  "in1-proc1-name",
								State: "in1-proc1-state",

								CPU: boshdir.VMInfoVitalsCPU{
									Total: &procCPUTotal,
								},
								Mem: boshdir.VMInfoVitalsMemIntSize{
									Percent: &procMemPercent,
									KB:      &procMemKB,
								},
								Uptime: boshdir.VMInfoVitalsUptime{
									Seconds: &procUptime,
								},
							},
							{
								Name:  "in1-proc2-name",
								State: "in1-proc2-state",
							},
						},
					},
					{
						JobName:      "job-name",
						Index:        &index2,
						ProcessState: "in2-process-state",
						AZ:           "in2-az",
						ResourcePool: "in2-rp",

						IPs: []string{"in2-ip1"},
						DNS: []string{"in2-dns1"},

						State:              "in2-state",
						VMID:               "in2-cid",
						AgentID:            "in2-agent-id",
						ResurrectionPaused: true,
						Ignore:             false,
						DiskIDs:            []string{"diskcid1", "diskcid2"},

						Vitals: boshdir.VMInfoVitals{
							Load: []string{"0.52", "0.56", "0.51"},

							CPU:  boshdir.VMInfoVitalsCPU{Sys: "50.3", User: "51.2", Wait: "52.1"},
							Mem:  boshdir.VMInfoVitalsMemSize{Percent: "60", KB: "6000"},
							Swap: boshdir.VMInfoVitalsMemSize{Percent: "61", KB: "6100"},

							Disk: map[string]boshdir.VMInfoVitalsDiskSize{
								"system":     boshdir.VMInfoVitalsDiskSize{Percent: "75"},
								"ephemeral":  boshdir.VMInfoVitalsDiskSize{Percent: "85"},
								"persistent": boshdir.VMInfoVitalsDiskSize{Percent: "95"},
							},
						},

						Processes: []boshdir.VMInfoProcess{
							{
								Name:  "in2-proc1-name",
								State: "in2-proc1-state",
							},
						},
					},
					{
						JobName:      "",
						Index:        nil,
						ProcessState: "unresponsive agent",
						ResourcePool: "",
					},
				}

				deployment.InstanceInfosReturns(infos, nil)
			})

			It("lists instances for the deployment", func() {
				Expect(act()).ToNot(HaveOccurred())

				Expect(ui.Table).To(Equal(boshtbl.Table{
					Content: "instances",

					HeaderVals: []boshtbl.Value{
						boshtbl.NewValueString("Instance"),
						boshtbl.NewValueString("Process State"),
						boshtbl.NewValueString("AZ"),
						boshtbl.NewValueString("IPs"),
					},

					SortBy: []boshtbl.ColumnSort{
						{Column: 0, Asc: true},
						{Column: 1, Asc: true},
					},

					Sections: []boshtbl.Section{
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-process-state"), true),
									boshtbl.ValueString{},
									boshtbl.NewValueStrings([]string{"in1-ip1", "in1-ip2"}),
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in2-process-state"), true),
									boshtbl.NewValueString("in2-az"),
									boshtbl.NewValueStrings([]string{"in2-ip1"}),
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("?"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("?"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("unresponsive agent"), true),
									boshtbl.ValueString{},
									boshtbl.ValueStrings{},
								},
							},
						},
					},

					Notes: []string{""},
				}))
			})

			It("lists instances with processes", func() {
				opts.Processes = true

				Expect(act()).ToNot(HaveOccurred())

				Expect(ui.Table).To(Equal(boshtbl.Table{
					Content: "instances",

					HeaderVals: []boshtbl.Value{
						boshtbl.NewValueString("Instance"),
						boshtbl.NewValueString("Process"),
						boshtbl.NewValueString("Process State"),
						boshtbl.NewValueString("AZ"),
						boshtbl.NewValueString("IPs"),
					},

					SortBy: []boshtbl.ColumnSort{
						{Column: 0, Asc: true},
						{Column: 1, Asc: true},
					},

					Sections: []boshtbl.Section{
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.ValueString{},
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-process-state"), true),
									boshtbl.ValueString{},
									boshtbl.NewValueStrings([]string{"in1-ip1", "in1-ip2"}),
								},
								{
									boshtbl.ValueString{},
									boshtbl.NewValueString("in1-proc1-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-proc1-state"), true),
									nil,
									nil,
								},
								{
									boshtbl.ValueString{},
									boshtbl.NewValueString("in1-proc2-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-proc2-state"), true),
									nil,
									nil,
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.ValueString{},
									boshtbl.NewValueFmt(boshtbl.NewValueString("in2-process-state"), true),
									boshtbl.NewValueString("in2-az"),
									boshtbl.NewValueStrings([]string{"in2-ip1"}),
								},
								{
									boshtbl.ValueString{},
									boshtbl.NewValueString("in2-proc1-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in2-proc1-state"), true),
									nil,
									nil,
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("?"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("?"),
									boshtbl.ValueString{},
									boshtbl.NewValueFmt(boshtbl.NewValueString("unresponsive agent"), true),
									boshtbl.ValueString{},
									boshtbl.ValueStrings{},
								},
							},
						},
					},

					Notes: []string{""},
				}))
			})

			It("lists instances for the deployment including details", func() {
				opts.Details = true

				Expect(act()).ToNot(HaveOccurred())

				Expect(ui.Table).To(Equal(boshtbl.Table{
					Content: "instances",

					HeaderVals: []boshtbl.Value{
						boshtbl.NewValueString("Instance"),
						boshtbl.NewValueString("Process State"),
						boshtbl.NewValueString("AZ"),
						boshtbl.NewValueString("IPs"),
						boshtbl.NewValueString("State"),
						boshtbl.NewValueString("VM CID"),
						boshtbl.NewValueString("VM Type"),
						boshtbl.NewValueString("Disk CIDs"),
						boshtbl.NewValueString("Agent ID"),
						boshtbl.NewValueString("Index"),
						boshtbl.NewValueString("Resurrection\nPaused"),
						boshtbl.NewValueString("Bootstrap"),
						boshtbl.NewValueString("Ignore"),
					},

					SortBy: []boshtbl.ColumnSort{
						{Column: 0, Asc: true},
						{Column: 1, Asc: true},
					},

					Sections: []boshtbl.Section{
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-process-state"), true),
									boshtbl.ValueString{},
									boshtbl.NewValueStrings([]string{"in1-ip1", "in1-ip2"}),
									boshtbl.NewValueString("in1-state"),
									boshtbl.NewValueString("in1-cid"),
									boshtbl.NewValueString("in1-rp"),
									boshtbl.NewValueStrings([]string{"diskcid1", "diskcid2"}),
									boshtbl.NewValueString("in1-agent-id"),
									boshtbl.NewValueInt(1),
									boshtbl.NewValueBool(false),
									boshtbl.NewValueBool(false),
									boshtbl.NewValueBool(true),
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in2-process-state"), true),
									boshtbl.NewValueString("in2-az"),
									boshtbl.NewValueStrings([]string{"in2-ip1"}),
									boshtbl.NewValueString("in2-state"),
									boshtbl.NewValueString("in2-cid"),
									boshtbl.NewValueString("in2-rp"),
									boshtbl.NewValueStrings([]string{"diskcid1", "diskcid2"}),
									boshtbl.NewValueString("in2-agent-id"),
									boshtbl.NewValueInt(2),
									boshtbl.NewValueBool(true),
									boshtbl.NewValueBool(false),
									boshtbl.NewValueBool(false),
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("?"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("?"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("unresponsive agent"), true),
									boshtbl.ValueString{},
									boshtbl.ValueStrings{},
									boshtbl.NewValueString(""),
									boshtbl.ValueString{},
									boshtbl.ValueString{},
									boshtbl.ValueStrings{},
									boshtbl.ValueString{},
									boshtbl.NewValueInt(0),
									boshtbl.NewValueBool(false),
									boshtbl.NewValueBool(false),
									boshtbl.NewValueBool(false),
								},
							},
						},
					},

					Notes: []string{""},
				}))
			})

			It("lists instances for the deployment including dns", func() {
				opts.DNS = true

				Expect(act()).ToNot(HaveOccurred())

				Expect(ui.Table).To(Equal(boshtbl.Table{
					Content: "instances",

					HeaderVals: []boshtbl.Value{
						boshtbl.NewValueString("Instance"),
						boshtbl.NewValueString("Process State"),
						boshtbl.NewValueString("AZ"),
						boshtbl.NewValueString("IPs"),
						boshtbl.NewValueString("DNS A Records"),
					},

					SortBy: []boshtbl.ColumnSort{
						{Column: 0, Asc: true},
						{Column: 1, Asc: true},
					},

					Sections: []boshtbl.Section{
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-process-state"), true),
									boshtbl.ValueString{},
									boshtbl.NewValueStrings([]string{"in1-ip1", "in1-ip2"}),
									boshtbl.NewValueStrings([]string{"in1-dns1", "in1-dns2"}),
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in2-process-state"), true),
									boshtbl.NewValueString("in2-az"),
									boshtbl.NewValueStrings([]string{"in2-ip1"}),
									boshtbl.NewValueStrings([]string{"in2-dns1"}),
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("?"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("?"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("unresponsive agent"), true),
									boshtbl.ValueString{},
									boshtbl.ValueStrings{},
									boshtbl.ValueStrings{},
								},
							},
						},
					},

					Notes: []string{""},
				}))
			})

			It("lists instances for the deployment including vitals and processes", func() {
				opts.Vitals = true
				opts.Processes = true

				Expect(act()).ToNot(HaveOccurred())

				Expect(ui.Table).To(Equal(boshtbl.Table{
					Content: "instances",

					HeaderVals: []boshtbl.Value{
						boshtbl.NewValueString("Instance"),
						boshtbl.NewValueString("Process"),
						boshtbl.NewValueString("Process State"),
						boshtbl.NewValueString("AZ"),
						boshtbl.NewValueString("IPs"),
						boshtbl.NewValueString("Uptime"),
						boshtbl.NewValueString("Load\n(1m, 5m, 15m)"),
						boshtbl.NewValueString("CPU\nTotal"),
						boshtbl.NewValueString("CPU\nUser"),
						boshtbl.NewValueString("CPU\nSys"),
						boshtbl.NewValueString("CPU\nWait"),
						boshtbl.NewValueString("Memory\nUsage"),
						boshtbl.NewValueString("Swap\nUsage"),
						boshtbl.NewValueString("System\nDisk Usage"),
						boshtbl.NewValueString("Ephemeral\nDisk Usage"),
						boshtbl.NewValueString("Persistent\nDisk Usage"),
					},

					SortBy: []boshtbl.ColumnSort{
						{Column: 0, Asc: true},
						{Column: 1, Asc: true},
					},

					Sections: []boshtbl.Section{
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.ValueString{},
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-process-state"), true),
									boshtbl.ValueString{},
									boshtbl.NewValueStrings([]string{"in1-ip1", "in1-ip2"}),
									ValueUptime{},
									boshtbl.NewValueString("0.02, 0.06, 0.11"),
									ValueCPUTotal{},
									NewValueStringPercent("1.2"),
									NewValueStringPercent("0.3"),
									NewValueStringPercent("2.1"),
									ValueMemSize{boshdir.VMInfoVitalsMemSize{Percent: "20", KB: "2000"}},
									ValueMemSize{boshdir.VMInfoVitalsMemSize{Percent: "21", KB: "2100"}},
									ValueDiskSize{boshdir.VMInfoVitalsDiskSize{Percent: "35"}},
									ValueDiskSize{boshdir.VMInfoVitalsDiskSize{Percent: "45"}},
									ValueDiskSize{boshdir.VMInfoVitalsDiskSize{Percent: "55"}},
								},
								{
									boshtbl.ValueString{},
									boshtbl.NewValueString("in1-proc1-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-proc1-state"), true),
									nil,
									nil,
									ValueUptime{&procUptime},
									nil,
									ValueCPUTotal{&procCPUTotal},
									nil,
									nil,
									nil,
									ValueMemIntSize{boshdir.VMInfoVitalsMemIntSize{Percent: &procMemPercent, KB: &procMemKB}},
									nil,
									nil,
									nil,
									nil,
								},
								{
									boshtbl.ValueString{},
									boshtbl.NewValueString("in1-proc2-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-proc2-state"), true),
									nil,
									nil,
									ValueUptime{},
									nil,
									ValueCPUTotal{},
									nil,
									nil,
									nil,
									ValueMemIntSize{},
									nil,
									nil,
									nil,
									nil,
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.ValueString{},
									boshtbl.NewValueFmt(boshtbl.NewValueString("in2-process-state"), true),
									boshtbl.NewValueString("in2-az"),
									boshtbl.NewValueStrings([]string{"in2-ip1"}),
									ValueUptime{},
									boshtbl.NewValueString("0.52, 0.56, 0.51"),
									ValueCPUTotal{},
									NewValueStringPercent("51.2"),
									NewValueStringPercent("50.3"),
									NewValueStringPercent("52.1"),
									ValueMemSize{boshdir.VMInfoVitalsMemSize{Percent: "60", KB: "6000"}},
									ValueMemSize{boshdir.VMInfoVitalsMemSize{Percent: "61", KB: "6100"}},
									ValueDiskSize{boshdir.VMInfoVitalsDiskSize{Percent: "75"}},
									ValueDiskSize{boshdir.VMInfoVitalsDiskSize{Percent: "85"}},
									ValueDiskSize{boshdir.VMInfoVitalsDiskSize{Percent: "95"}},
								},
								{
									boshtbl.ValueString{},
									boshtbl.NewValueString("in2-proc1-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in2-proc1-state"), true),
									nil,
									nil,
									ValueUptime{},
									nil,
									ValueCPUTotal{},
									nil,
									nil,
									nil,
									ValueMemIntSize{},
									nil,
									nil,
									nil,
									nil,
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("?"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("?"),
									boshtbl.ValueString{},
									boshtbl.NewValueFmt(boshtbl.NewValueString("unresponsive agent"), true),
									boshtbl.ValueString{},
									boshtbl.ValueStrings{},
									ValueUptime{},
									boshtbl.ValueString{},
									ValueCPUTotal{},
									NewValueStringPercent(""),
									NewValueStringPercent(""),
									NewValueStringPercent(""),
									ValueMemSize{},
									ValueMemSize{},
									ValueDiskSize{},
									ValueDiskSize{},
									ValueDiskSize{},
								},
							},
						},
					},

					Notes: []string{""},
				}))
			})

			It("lists failing (non-running) instances", func() {
				opts.Failing = true

				// Hides second VM
				infos[1].ProcessState = "running"
				infos[1].Processes[0].State = "running"

				Expect(act()).ToNot(HaveOccurred())

				Expect(ui.Table).To(Equal(boshtbl.Table{
					Content: "instances",

					HeaderVals: []boshtbl.Value{
						boshtbl.NewValueString("Instance"),
						boshtbl.NewValueString("Process State"),
						boshtbl.NewValueString("AZ"),
						boshtbl.NewValueString("IPs"),
					},

					SortBy: []boshtbl.ColumnSort{
						{Column: 0, Asc: true},
						{Column: 1, Asc: true},
					},

					Sections: []boshtbl.Section{
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-process-state"), true),
									boshtbl.ValueString{},
									boshtbl.NewValueStrings([]string{"in1-ip1", "in1-ip2"}),
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("?"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("?"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("unresponsive agent"), true),
									boshtbl.ValueString{},
									boshtbl.ValueStrings{},
								},
							},
						},
					},

					Notes: []string{""},
				}))
			})

			It("includes failing processes when listing failing (non-running) instances and processes", func() {
				opts.Failing = true
				opts.Processes = true

				// Hides first process in the first VM
				infos[0].Processes[0].State = "running"

				// Hides second VM completely
				infos[1].ProcessState = "running"
				infos[1].Processes[0].State = "running"

				Expect(act()).ToNot(HaveOccurred())

				Expect(ui.Table).To(Equal(boshtbl.Table{
					Content: "instances",

					HeaderVals: []boshtbl.Value{
						boshtbl.NewValueString("Instance"),
						boshtbl.NewValueString("Process"),
						boshtbl.NewValueString("Process State"),
						boshtbl.NewValueString("AZ"),
						boshtbl.NewValueString("IPs"),
					},

					SortBy: []boshtbl.ColumnSort{
						{Column: 0, Asc: true},
						{Column: 1, Asc: true},
					},

					Sections: []boshtbl.Section{
						{
							FirstColumn: boshtbl.NewValueString("job-name"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("job-name"),
									boshtbl.ValueString{},
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-process-state"), true),
									boshtbl.ValueString{},
									boshtbl.NewValueStrings([]string{"in1-ip1", "in1-ip2"}),
								},
								{
									boshtbl.ValueString{},
									boshtbl.NewValueString("in1-proc2-name"),
									boshtbl.NewValueFmt(boshtbl.NewValueString("in1-proc2-state"), true),
									nil,
									nil,
								},
							},
						},
						{
							FirstColumn: boshtbl.NewValueString("?"),
							Rows: [][]boshtbl.Value{
								{
									boshtbl.NewValueString("?"),
									boshtbl.ValueString{},
									boshtbl.NewValueFmt(boshtbl.NewValueString("unresponsive agent"), true),
									boshtbl.ValueString{},
									boshtbl.ValueStrings{},
								},
							},
						},
					},

					Notes: []string{""},
				}))
			})
		})

		It("returns error if instances cannot be retrieved", func() {
			deployment.InstanceInfosReturns(nil, errors.New("fake-err"))

			err := act()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})
	})
})
