//Texas Instruments INA219 high side current sensor

package ina219

import (
	"fmt"

	"github.com/NeuralSpaz/i2c"
)

const (
	INA219_ADDRESS01 = 0x40
	INA219_ADDRESS02 = 0x41

	INA219_CONFIG_RESET                     = 0x8000
	INA219_CONFIG_BVOLTAGERANGE_MASK        = 0x2000
	INA219_CONFIG_BVOLTAGERANGE_16V         = 0x0000
	INA219_CONFIG_BVOLTAGERANGE_32V         = 0x2000
	INA219_CONFIG_GAIN_MASK                 = 0x1800
	INA219_CONFIG_GAIN_1_40MV               = 0x0000
	INA219_CONFIG_GAIN_2_80MV               = 0x0800
	INA219_CONFIG_GAIN_4_160MV              = 0x1000
	INA219_CONFIG_GAIN_8_320MV              = 0x1800
	INA219_CONFIG_BADCRES_MASK              = 0x0780
	INA219_CONFIG_BADCRES_9BIT              = 0x0080
	INA219_CONFIG_BADCRES_10BIT             = 0x0100
	INA219_CONFIG_BADCRES_11BIT             = 0x0200
	INA219_CONFIG_BADCRES_12BIT             = 0x0400
	INA219_CONFIG_SADCRES_MASK              = 0x0078
	INA219_CONFIG_SADCRES_9BIT_1S_84US      = 0x0000
	INA219_CONFIG_SADCRES_10BIT_1S_148US    = 0x0008
	INA219_CONFIG_SADCRES_11BIT_1S_276US    = 0x0010
	INA219_CONFIG_SADCRES_12BIT_1S_532US    = 0x0018
	INA219_CONFIG_SADCRES_12BIT_2S_1060US   = 0x0048
	INA219_CONFIG_SADCRES_12BIT_4S_2130US   = 0x0050
	INA219_CONFIG_SADCRES_12BIT_8S_4260US   = 0x0058
	INA219_CONFIG_SADCRES_12BIT_16S_8510US  = 0x0060
	INA219_CONFIG_SADCRES_12BIT_32S_17MS    = 0x0068
	INA219_CONFIG_SADCRES_12BIT_64S_34MS    = 0x0070
	INA219_CONFIG_SADCRES_12BIT_128S_69MS   = 0x0078
	INA219_CONFIG_MODE_MASK                 = 0x0007
	INA219_CONFIG_MODE_POWERDOWN            = 0x0000
	INA219_CONFIG_MODE_SVOLT_TRIGGERED      = 0x0001
	INA219_CONFIG_MODE_BVOLT_TRIGGERED      = 0x0002
	INA219_CONFIG_MODE_SANDBVOLT_TRIGGERED  = 0x0003
	INA219_CONFIG_MODE_ADCOFF               = 0x0004
	INA219_CONFIG_MODE_SVOLT_CONTINUOUS     = 0x0005
	INA219_CONFIG_MODE_BVOLT_CONTINUOUS     = 0x0006
	INA219_CONFIG_MODE_SANDBVOLT_CONTINUOUS = 0x0007

	INA219_REG_CONFIG       = 0x00
	INA219_REG_SHUNTVOLTAGE = 0x01
	INA219_REG_BUSVOLTAGE   = 0x02
	INA219_REG_POWER        = 0x03
	INA219_REG_CURRENT      = 0x04
	INA219_REG_CALIBRATION  = 0x05
)

type INA219 struct {
	Dev              i2c.I2CBus
	init             bool
	Address          uint8
	Configuation     uint16
	CalibrationValue uint16
	Current          float64
	Shunt            float64
	Bus              float64
	Power            float64
	Load             float64
}

func (d *INA219) String() string {
	return fmt.Sprintf("BusVolts %f ShuntVolts: %f Current: %f Power: %f", d.Bus, d.Shunt, d.Current, d.Power)
}

func New(deviceAdress uint8, i2cbus byte) *INA219 {
	deviceBus := i2c.NewI2CBus(i2cbus)
	d := &INA219{
		Dev: deviceBus,
		Configuation: INA219_CONFIG_BVOLTAGERANGE_16V |
			INA219_CONFIG_GAIN_8_320MV |
			INA219_CONFIG_BADCRES_12BIT |
			INA219_CONFIG_SADCRES_12BIT_128S_69MS |
			INA219_CONFIG_MODE_SANDBVOLT_CONTINUOUS,
		CalibrationValue: 6826,
		Address:          deviceAdress,
		// err:   make(chan error),
		// temp:  make(chan float64),
		// state: make(chan statefn),
	}
	return d
}

// Fetch all values from INA219
func Fetch(d *INA219) error {

	if !d.init {
		// Write Calibration Register
		if err := d.Dev.WriteWordToReg(d.Address, INA219_REG_CONFIG, d.Configuation); err != nil {
			return err
		}
		// Write Configuration Register
		if err := d.Dev.WriteWordToReg(d.Address, INA219_REG_CALIBRATION, d.CalibrationValue); err != nil {
			return err
		}
		d.init = true

	}
	// Fetch Values Volatge
	Bus, err := d.Dev.ReadWordFromReg(d.Address, INA219_REG_BUSVOLTAGE)
	if err != nil {
		return err
	}
	Shunt, err := d.Dev.ReadWordFromReg(d.Address, INA219_REG_SHUNTVOLTAGE)
	if err != nil {
		return err
	}
	Current, err := d.Dev.ReadWordFromReg(d.Address, INA219_REG_CURRENT)
	if err != nil {
		return err
	}
	// fmt.Println(Current)
	// fmt.Println(int16(Current))
	Power, err := d.Dev.ReadWordFromReg(d.Address, INA219_REG_POWER)
	if err != nil {
		return err
	}
	// if Bus & 1 {
	// 	err := errors.New("INA219 Error: MathOverFlow, See data sheet for correct values")
	// 	return error
	// }
	d.Bus = float64(int16(Bus>>3)*4) * 0.001
	d.Shunt = float64(int16(Shunt)) * 0.00001

	d.Current = float64(int16(Current)) * 0.0004
	d.Power = float64(int16(Power)) * 4 * 0.4 * 5 * 0.001
	return nil

}
