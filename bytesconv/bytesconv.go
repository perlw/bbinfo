package bytesconv

import "fmt"

const (
	SizeKB = 1024
	SizeMB = SizeKB * 1024
	SizeGB = SizeMB * 1024
	SizeTB = SizeGB * 1024
	SizePB = SizeTB * 1024
)

func Qualifier(bytes int) int {
	switch {
	case bytes >= SizePB:
		return SizePB
	case bytes >= SizeTB:
		return SizeTB
	case bytes >= SizeGB:
		return SizeGB
	case bytes >= SizeMB:
		return SizeMB
	case bytes >= SizeKB:
		return SizeKB
	default:
		return 1
	}
}

func QualifierToString(qualifier int) string {
	switch qualifier {
	case SizePB:
		return "PB"
	case SizeTB:
		return "TB"
	case SizeGB:
		return "GB"
	case SizeMB:
		return "MB"
	case SizeKB:
		return "KB"
	default:
		return "B"
	}
}

func ToHumanReadable(bytes int) string {
	qual := Qualifier(bytes)
	return fmt.Sprintf("%.2f%s", float32(bytes)/float32(qual), QualifierToString(qual))
}

func QualifyTransfer(upBytes, downBytes int) (up float32, down float32, qualifier string) {
	upQualifier := Qualifier(upBytes)
	downQualifier := Qualifier(downBytes)

	if upQualifier > downQualifier {
		up = float32(upBytes) / float32(upQualifier)
		down = float32(downBytes) / float32(upQualifier)
		qualifier = QualifierToString(upQualifier)
	} else {
		up = float32(upBytes) / float32(downQualifier)
		down = float32(downBytes) / float32(downQualifier)
		qualifier = QualifierToString(downQualifier)
	}

	return up, down, qualifier
}
