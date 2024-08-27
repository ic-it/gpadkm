package gpadkm

func syscallEVIOCGNAME(len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x06, len)
}

func _IOC(dir, typ, nr, size int) uintptr {
	return uintptr((dir << _IOC_DIRSHIFT) | (typ << _IOC_TYPESHIFT) | (nr << _IOC_NRSHIFT) | (size << _IOC_SIZESHIFT))
}

const (
	_IOC_READ  = 2
	_IOC_WRITE = 1

	_IOC_DIRSHIFT  = 30
	_IOC_TYPESHIFT = 8
	_IOC_NRSHIFT   = 0
	_IOC_SIZESHIFT = 16
)
