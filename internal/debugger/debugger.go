package debugger

type Debugger interface {
	ReadMemoryHook(memory uint32, addr uint64, width int)
	WriteMemoryHook(memory uint32, addr uint64, width int, data uint64)
	PrintLog(message string)
	InstructionHook(id int, pc uint64) // breakpointのチェックはここで行う
}
