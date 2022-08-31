
// +build !linux


package common


import "syscall"




func Syscall_Dup2or3 (_oldfd int, _newfd int) error {
	return syscall.Dup2 (_oldfd, _newfd)
}

