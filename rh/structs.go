package rh

// Information about the platform
type Platform struct {
	Arch     string `json:"arch"`
	Os       string `json:"os"`
	Platform string `json:"platform"`
	Family   string `json:"family"`
	Kernel   string `json:"kernel"`
}

// Information on the Cpu
type CpuInfo struct {
	VendorId  string  `json:"vendorId"`
	Mhz       float64 `json:"mhz"`
	CacheSize int32   `json:"cacheSize"`
	Cores     int     `json:"cores"`
	Threads   int     `json:"threads"`
}

// Information on current system resource usage
type SystemLoad struct {
	Cpu    CpuLoad      `json:"cpu"`
	Memory SystemMemory `json:"memory"`
}

// The load on the CPU
type CpuLoad struct {
	Total   float64   `json:"jotal"`
	PerCore []CpuCore `json:"perCore"`
}

// The load information on a single Cpu core
type CpuCore struct {
	CoreNo int     `json:"coreNumber"`
	Load   float32 `json:"load"`
}

// The information on the memory usage
type SystemMemory struct {
	Unit      string  `json:"unit"`
	Total     float32 `json:"total"`
	Available float32 `json:"available"`
	Used      float32 `json:"used"`
}

// The total time the system has been powered on
type SystemUptime struct {
	Years   int `json:"years"`
	Months  int `json:"months"`
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}
