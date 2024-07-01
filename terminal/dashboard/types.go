package dashboard

var (
	statusLabelMap map[int]string = map[int]string{
		0: "Online",
		1: "Offline",
		2: "Disabled",
		3: "Hashrate Error",
		4: "Temperature Error",
		5: "Fan Speed Error",
		6: "Missing Hashboard Error",
		7: "Pool Share Error",
	}

	modeLabelMap map[int]string = map[int]string{
		0: "Normal",
		1: "Sleep",
		2: "Low Power",
	}
)
