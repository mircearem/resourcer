package rh

import (
	"context"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type ResourceHandler struct {
	// Mutex for thread safety
	mu sync.RWMutex
	// Context for deterministic runtime
	ctx context.Context
	// Channel for error checking
	errCh chan error
	// Channel to listen for request and send back response
	Platform platform     `json:"platform"`
	Cpu      cpuInfo      `json:"cpu"`
	Loads    systemLoad   `json:"systemLoad"`
	Uptime   systemUptime `json:"uptime"`
}

// Information about the platform
type platform struct {
	Arch     string `json:"arch"`
	Os       string `json:"os"`
	Platform string `json:"platform"`
	Family   string `json:"family"`
	Kernel   string `json:"kernel"`
}

// Information on the Cpu
type cpuInfo struct {
	VendorId  string  `json:"vendorId"`
	Mhz       float64 `json:"mhz"`
	CacheSize int32   `json:"cacheSize"`
	Cores     int     `json:"cores"`
	Threads   int     `json:"threads"`
}

// Information on current system resource usage
type systemLoad struct {
	Cpu    cpuLoad      `json:"cpu"`
	Memory systemMemory `json:"memory"`
}

// The load on the CPU
type cpuLoad struct {
	Usage float64   `json:"total"`
	Loads []cpuCore `json:"cpuCore"`
}

// The load information on a single Cpu core
type cpuCore struct {
	CoreNo int     `json:"coreNumber"`
	Load   float32 `json:"load"`
}

// The information on the memory usage
type systemMemory struct {
	Unit      string  `json:"unit"`
	Total     float32 `json:"total"`
	Available float32 `json:"available"`
	Used      float32 `json:"used"`
}

// The total time the system has been powered on
type systemUptime struct {
	Years   int `json:"years"`
	Months  int `json:"months"`
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

// Create a new Resource Handler
func NewHandler(ctx context.Context) *ResourceHandler {
	return &ResourceHandler{
		ctx:   ctx,
		errCh: make(chan error),
	}
}

// Initialize the Resource Handler
func (r *ResourceHandler) init() error {
	var wg sync.WaitGroup

	numCnt := 4
	wg.Add(numCnt)

	errors := make(chan error, numCnt)

	for i := 0; i < numCnt; i += 1 {
		go func(id int) {
			defer wg.Done()
			// Determine which function to run
			switch id {
			// Get information on the platform
			case 0:
				err := r.getPlatformInformation()
				if err != nil {
					errors <- err
				}
			// Get the core count
			case 1:
				err := r.getCpuCoreCount()
				if err != nil {
					errors <- err
				}
			// Get the thread count
			case 2:
				err := r.getCpuThreadCount()
				if err != nil {
					errors <- err
				}
			// Get the cpu information
			case 3:
				err := r.getCpuInformation()
				if err != nil {
					errors <- err
				}
			}
		}(i)
	}
	// Wait for all the go routines to finish their jobs
	go func() {
		wg.Wait()
		close(errors)
	}()
	// Return the errors if any
	for err := range errors {
		if err != nil {
			return err
		}
	}
	// No errors
	return nil
}

// Function to read the information about the platform
func (r *ResourceHandler) getPlatformInformation() error {
	errch := make(chan error)

	go func() {
		// Lock the mutex
		r.mu.Lock()
		defer r.mu.Unlock()
		// Read information about the platform
		platformInfo, err := host.InfoWithContext(r.ctx)
		// Error reading information about the platform
		if err != nil {
			errch <- err
			return
		}
		// Write the platform information
		r.Platform.Arch = platformInfo.KernelArch
		r.Platform.Platform = platformInfo.Platform
		r.Platform.Os = platformInfo.PlatformVersion
		r.Platform.Kernel = platformInfo.KernelVersion
		r.Platform.Family = platformInfo.PlatformFamily
		// No error
		errch <- nil
	}()

	err := <-errch
	close(errch)

	return err
}

// Function to read the number of physical cores on the CPU
func (r *ResourceHandler) getCpuCoreCount() error {
	errch := make(chan error)

	// Get the CPU physical cores
	go func() {
		// Lock the mutex
		r.mu.Lock()
		defer r.mu.Unlock()
		// Read the number of cpu cores
		cores, err := cpu.CountsWithContext(r.ctx, false)
		// Error reading the number of cores
		if err != nil {
			errch <- err
			return
		}
		// Write the number of CPU cores
		r.Cpu.Cores = cores
		// No error
		errch <- nil
	}()

	err := <-errch
	close(errch)

	return err
}

// Function to read the number of threads on the CPU
func (r *ResourceHandler) getCpuThreadCount() error {
	errch := make(chan error)

	// Get the CPU physical cores
	go func() {
		// Lock the mutex
		r.mu.Lock()
		defer r.mu.Unlock()
		// Read the number of cpu cores
		threads, err := cpu.CountsWithContext(r.ctx, true)
		// Error reading the number of cores
		if err != nil {
			errch <- err
			return
		}
		// Write the number of CPU cores
		r.Cpu.Cores = threads
		// No error
		errch <- nil
	}()

	err := <-errch
	close(errch)

	return err
}

// Function to get the information on the CPU
func (r *ResourceHandler) getCpuInformation() error {
	errch := make(chan error)

	// Get the CPU information
	go func() {
		// Lock the mutex
		r.mu.Lock()
		defer r.mu.Unlock()
		// Read other information on the CPU
		cpuInfo, err := cpu.InfoWithContext(r.ctx)
		// Error reading the CPU information
		if err != nil {
			errch <- err
			return
		}
		// Write the CPU information
		r.Cpu.VendorId = cpuInfo[0].VendorID
		r.Cpu.CacheSize = cpuInfo[0].CacheSize
		r.Cpu.Mhz = cpuInfo[0].Mhz
		// No error
		errch <- nil
	}()

	err := <-errch
	close(errch)

	return err
}

func (r *ResourceHandler) Run() {
	// Initialize the server first
	if err := r.init(); err != nil {
		log.Fatalln(err)
	}
	// Set the timers for the update functions
	t1 := time.NewTicker(3 * time.Second)
	// Ticker to read the memory usage
	t2 := time.NewTicker(5 * time.Second)
	// Ticker to read the system uptime
	t3 := time.NewTicker(1 * time.Second)
	// Update the values
	go func() {
		for {
			select {
			case <-t1.C:
				// Update the CPU loads
				if err := r.getCpuLoad(); err != nil {
					r.errCh <- err
				}

			case <-t2.C:
				// Update the memory usage
				log.Println("update memory")
			case <-t3.C:
				// Upate the system uptime
				r.getSystemUptime()
			}
		}
	}()
	// Check for errors
	for err := range r.errCh {
		log.Println(err)
	}
}

// To determine the largest unit of the memory (GB, MB, KB, etc) calculate the power of 10
func powerOfTen(input uint64) uint64 {
	// Convert the input into a string
	inputStr := strconv.FormatUint(input, 10)

	// Count the number of digits
	numDigits := uint64(len(inputStr))

	return numDigits - 1
}

// Read the current load on each CPU core
func (r *ResourceHandler) getCpuLoad() error {
	errch := make(chan error, 2)

	// Read the total CPU load
	go func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		// Get the usage total load of the CPU
		usage, err := cpu.PercentWithContext(r.ctx, time.Second, false)
		// Error when trying to read CPU load
		if err != nil {
			errch <- err
			return
		}
		r.Loads.Cpu.Usage = usage[0]
		errch <- nil
	}()
	// Read the usage for each core
	go func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		// Get the usage for each CPU core
		loads, err := cpu.PercentWithContext(r.ctx, time.Second, true)
		// Error when trying to read loads for each core
		if err != nil {
			errch <- err
			return
		}
		// Format the information
		cores := make([]cpuCore, len(loads))

		for cpu, load := range loads {
			core := cpuCore{
				CoreNo: cpu,
				Load:   float32(load),
			}

			cores[cpu] = core
		}
		// Write the information into the struct
		r.Loads.Cpu.Loads = cores
		errch <- nil
	}()
	// Check for errors
	for err := range errch {
		if err != nil {
			return err
		}
	}
	return nil
}

// Read information on the system memory
func (r *ResourceHandler) getSystemMemory() error {
	errch := make(chan error)

	go func() {
		mem, err := mem.VirtualMemoryWithContext(r.ctx)
		// Error while reading the information on the system memory
		if err != nil {
			errch <- err
			return
		}
		// Calculate the values that will be used in the struct
		var total float64
		var available float64
		var unit string
		// Create the output data structures
		sysMem := new(systemMemory)
		powTen := powerOfTen(mem.Total)
		// Format the data
		switch powTen {
		// Megabytes
		case 6:
			total = float64(mem.Total) / math.Pow(10, float64(powTen))
			available = float64(mem.Available) / math.Pow(10, float64(powTen))
			unit = "mb"
		// Gigabyte
		case 9:
			total = float64(mem.Total) / math.Pow(10, float64(powTen))
			available = float64(mem.Available) / math.Pow(10, float64(powTen))
			unit = "gb"
		// Kilobytes
		default:
			total = float64(mem.Total) / math.Pow(10, 3)
			available = float64(mem.Available) / math.Pow(10, 3)
			unit = "kb"
		}
		// Update the data in the struct
		sysMem.Total = float32(total)
		sysMem.Available = float32(available)
		sysMem.Used = float32(mem.UsedPercent)
		sysMem.Unit = unit
		r.Loads.Memory = *sysMem
		// No errors
		errch <- nil
	}()

	err := <-errch
	return err
}

// Get the system uptime
func (r *ResourceHandler) getSystemUptime() error {
	errch := make(chan error)

	go func() {
		// Read the information on the system uptime
		up, err := host.UptimeWithContext(r.ctx)
		// Check for errors
		if err != nil {
			errch <- err
			return
		}
		// Get the current time to make the calculations
		now := time.Now()
		then := now.Add(-time.Second * time.Duration(up))
		// Calculate each filed of the uptime
		years := now.Year() - then.Year()
		months := int(now.Month() - then.Month())
		days := now.Day() - then.Day()
		hours := now.Hour() - then.Hour()
		minutes := now.Minute() - then.Minute()
		// Adjust for negative values in minutes
		if minutes < 0 {
			minutes += 60
			hours -= 1
		}
		// Adjust for negative values in hours
		if hours < 0 {
			hours += 24
			days -= 1
		}
		// Adjust for negative values in days
		if days < 0 {
			lastMonth := now.AddDate(0, 1, -now.Day())
			days += int(then.Sub(lastMonth).Hours() / 24)
			months += 1
		}
		// Adjust for negative values in months
		if months < 0 {
			months += 12
			years += 1
		}
		// Create the return object
		uptime := systemUptime{
			Years:   years,
			Months:  months,
			Days:    days,
			Hours:   hours,
			Minutes: minutes,
		}
		r.Uptime = uptime
		// No errors
		errch <- nil
	}()
	err := <-errch
	return err
}
