package report

import "bytes"

var (
	toLevel = map[Level]string{
		Low: "Low",
		Medium: "Medium",
		High: "High",
		Critical: "Critical",
	}

	toDetectType = map[DetectType]string{
		Image: "Image",
		Container: "Container",
	}

	toEventType = map[EventType]string{
		Risk: "Risk",
		Invasion: "Invasion",
	}

	toAlertType = map[AlertType]string{
		Vulnerability: "Vulnerability",
		MaliciousFile: "MaliciousFile",
		Backdoor: "Backdoor",
		Sensitive: "Sensitive",
		AbnormalHistory: "AbnormalHistory",
		Weakpass: "Weakpass",
	}

	toWeakpassService = map[WeakpassService]string{
		SSH: "SSH",
	}
)

func (l Level) MarshalJSON()([]byte, error){
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toLevel[l])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d DetectType) MarshalJSON()([]byte, error){
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toDetectType[d])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e EventType) MarshalJSON()([]byte, error){
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toEventType[e])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (a AlertType) MarshalJSON()([]byte, error){
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toAlertType[a])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (w WeakpassService) MarshalJSON()([]byte, error){
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toWeakpassService[w])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}