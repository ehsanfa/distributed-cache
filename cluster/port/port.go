package port

import "strconv"

func ConvertPort(p string) (uint16, error) {
	port, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return uint16(port), nil
}
